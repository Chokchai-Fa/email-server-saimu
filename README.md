# Email Server Saimu
Email Notification Server with SQS Consumer

with this project we need to set up aws sqs first -> see document at https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/welcome.html

## How to run with Go.
1. Install Library and Dependencies: "go mod tidy"
2. run "go run main.go"

## How to run with DockerFile.
1. "docker build -t emailserver-saimu ."
2. "docker run  -d -t be-emailserver-saimu"

## If you want to deploy with container based, you can use this docker image
https://hub.docker.com/r/chokchaifa/saimu-email-server

### example
1. docker pull chokchaifa/saimu-email-server
2. docker run -d -t chokchaifa/saimu-email-server