package email

import (
	"emailserver-saimu/models"
	"emailserver-saimu/utils/logs"
	"fmt"
	"strings"

	"gopkg.in/gomail.v2"
)

type EmailServer struct {
	From       string
	SMTPServer string
	SMTPPort   int
	Email      string
	Password   string
}

func NewEmailServer(serverEmail string, serverPassword string, smtpServer string, smtpPort int) (EmailServer, error) {
	return EmailServer{serverEmail, smtpServer, smtpPort, serverEmail, serverPassword}, nil
}

func (h *EmailServer) SendSaimu(data models.EmailNoti) error {
	message := gomail.NewMessage()
	message.SetHeader("From", h.From)
	message.SetHeader("To", data.Email)
	message.SetHeader("Subject", "สีเสื้อมงคลประจำวันนี้")

	msgBody := strings.Replace(data.Conclusion, "\n", "<br>", -1)

	message.SetBody("text/html", fmt.Sprintf(`<!DOCTYPE html>
	<html>
	<head>
	  <link rel="preconnect" href="https://fonts.googleapis.com">
	  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
	  <link href="https://fonts.googleapis.com/css2?family=Sarabun:ital,wght@0,100;0,200;0,300;0,400;0,500;0,600;0,700;0,800;1,100;1,200;1,300;1,400;1,500;1,600;1,700;1,800&display=swap" rel="stylesheet">
	  <style>
		body {
		  font-family: Sarabun, sans-serif;
		  font-size: 17px;
		}
		.container {
		  max-width: 600px;
		  margin: 0 auto;
		  padding: 20px;
		  border: 1px solid #ddd;
		  border-radius: 5px;
		}
		.header {
		  margin-bottom: 20px;
		  color: #06142B;
		}
		.case-summary {
		  background-color:  #E6007E;
		  padding: 10px;
		  border-radius: 5px;
		  display: inline-block;
		  border: none;
		  color: #ffff
		}
		.case-info {
		  margin-bottom: 20px;
		  background-color: #F9FAFB;
		  padding: 15px;
		  border-top: 2px solid #091E42;
		}
		.case-info{
		  margin-bottom: 5px;
		}
		.label-text{
		  color: #626F86;
		  font-size: 15px;
		}
		.approve-text{
		  color: #008060;
		  font-size: 17px;
		}
		.reject-text{
		  color:#BC2200;
		  font-size: 17px;
		}
	  </style>
	</head>
	<body>
	  <div class="container">
		<div class="case-info">
		  <p style="font-size: 20px;color: #E6007E;">สีเสื้อมงคลประจำวันนี้</h1>
		  <p >%s<p>
		</div>
	  </div>
	</body>
	</html>`, msgBody))

	// Build the dialer and send the email.
	dialer := gomail.NewDialer(h.SMTPServer, h.SMTPPort, h.Email, h.Password)

	err := dialer.DialAndSend(message)
	if err != nil {
		errMsg := fmt.Sprintf("error on send email | %s", data.Email)
		logs.Error(errMsg)
	}

	return err
}
