package main

import (
	"context"
	"database/sql"
	"emailserver-saimu/models"
	"emailserver-saimu/utils/dbpostgres"
	"emailserver-saimu/utils/email"
	"emailserver-saimu/utils/logs"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"os"
	"strconv"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		logs.Error("DotEnv|:" + err.Error())
	}

	sqsURL := os.Getenv("SQS_URL")
	logs.Info("SQS: " + sqsURL)
}

func connectDBPostgres() (*dbpostgres.DBPG, *sql.DB) {
	port, err := strconv.Atoi(os.Getenv("PG_PORT"))
	if err != nil {
		logs.Error("PostgresPort|:" + err.Error())
		os.Exit(1)
	}

	dbPG := dbpostgres.NewDBPostgres()
	dbPGCli, err := dbPG.Connect(
		os.Getenv("PG_USER"),
		os.Getenv("PG_PASSWORD"),
		os.Getenv("PG_HOST"),
		os.Getenv("PG_DBNAME"),
		port,
	)
	if err != nil {
		logs.Error("Postgres.Connect|:" + err.Error())
		os.Exit(1)
	}

	// set the maximum life time of a connection
	dbPGCli.SetConnMaxLifetime(5 * time.Minute)

	// set the maximum number of connections in the pool
	dbPGCli.SetMaxOpenConns(50)

	// set the maximum idle connections
	dbPGCli.SetMaxIdleConns(25)

	if err := dbPGCli.Ping(); err != nil {
		logs.Error("Postgres.Ping|:" + err.Error())
		os.Exit(1)
	}

	return dbPG, dbPGCli
}

func connectSQS() *sqs.Client {
	cfg, err := awsconfig.LoadDefaultConfig(context.Background())
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}

	// Create an SQS client
	svc := sqs.NewFromConfig(cfg)

	return svc
}

func main() {

	sqsURL := os.Getenv("SQS_URL")

	queue := connectSQS()

	smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		logs.Error("SMTPPortNumber|:" + err.Error())
		os.Exit(1)
	}

	emailServer, err := email.NewEmailServer(
		os.Getenv("SERVER_EMAIL"),
		os.Getenv("SERVER_PASSWORD"),
		os.Getenv("SMTP_SERVER"),
		smtpPort,
	)
	if err != nil {
		logs.Error("MailServer|:" + err.Error())
		os.Exit(1)
	}
	logs.Info("Starting Queue Service")
	var processedMessageIDs = make(map[string]bool)

	for {
		result, err := queue.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
			QueueUrl:            &sqsURL,
			MaxNumberOfMessages: *aws.Int32(10), // Maximum number of messages to retrieve
		})
		if err != nil {
			panic("unable to receive message from queue" + err.Error())
		} else {
			var wg sync.WaitGroup
			errChan := make(chan error, len(result.Messages))

			wg.Add(len(result.Messages))

			for _, msg := range result.Messages {
				go func(msg types.Message) {

					messageID := *msg.MessageId
					if processedMessageIDs[messageID] {
						// Skip processing if the message has already been processed
						errChan <- nil
						return
					}

					logs.Info(fmt.Sprintf("Received message: %s: %s", *msg.MessageId, *msg.Body))

					msgBody := strings.ReplaceAll(*msg.Body, `'`, `"`)

					var emailNoti models.EmailNoti
					json.Unmarshal([]byte(msgBody), &emailNoti)

					err = emailServer.SendEmailLuckyShirt(emailNoti)
					if err != nil {
						errChan <- fmt.Errorf(fmt.Sprintf("send email error: %s", err))
					} else {
						logs.Info(fmt.Sprintf("send email to %s successfully", emailNoti.Email))
					}

					processedMessageIDs[messageID] = true

					_, err := queue.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
						QueueUrl:      &sqsURL,
						ReceiptHandle: msg.ReceiptHandle,
					})
					if err != nil {
						errChan <- fmt.Errorf(fmt.Sprintf("error deleting message: %s", err))
					}

					errChan <- nil
				}(msg)

			}

			go func() {
				wg.Wait()
				close(errChan)
			}()

			for range result.Messages {
				err := <-errChan
				if err != nil {
					logs.Error(err)
				}

			}

		}

	}

}
