# README #
You can find samples of the REQUESTS in the Postman file: printer-timeline-backend.postman_collection.json at the root level of the project. 

Time type can be relative or absolute. You can find more information on the request parameters in app/internal/queryparams/query_params_test.go 

package queryparams extracts the query parameters common to more than one endpoint.

In app/internal/api/api.go you will find the different endpoints and the corresponding handlers.

You can follow a top-down approach to better understand the code. Also take a look at the test files.

### What is this repository for? ###

* It is the back-end of the Printer-timeline. 
This project is a Go API that let's you retrieve in chronological order different data elements of the Cloud Connector (OpenXMls, CloudJsons, etc.) specifying a time range and a pn (Product Number) and sn(Serial number), if you don't specify them you will obtain all the pn's and sn's in that time range. It also has an endpoint to retrieve AWS S3 objects by bucket name and object key.

### How do I get set up? ###
The file dev.env contains the environment variables. 

Set env variable DEVELOPMENT=true in file dev.env when whant to use development mode (it will create a local server at http://0.0.0.0:8080).

Set env variable DEVELOPMENT=false in file dev.env when whant to deploy it to AWS, it will not run in a local server. Instead it will be expecting to execute in an AWS Lambda.

Set env variable MAX_TIME_DIFF_IN_MINUTES with the value of the max difference between start time and end time that you want to allow for the queries. 

Nevertheless, it also has a limit hardcoded in init.go (setMaxTimeDiffInMinutes()). It has that limit to prevent expensive query scans.
If MAX_TIME_DIFF_IN_MINUTES is not set or it is not a valid integer, the default value will be 60 minutes.

* Summary of set up (check that you are at the root of the project before executing the commands because some commands include a 'cd' as part of the command itself 🧐).

First export the environment variables of dev.env. 

Yo can execute:

$ for line in $(cat dev.env); do export $line; done; 

Execute the app to run in local server:
First, set env variable DEVELOPMENT=true in file dev.env.

Then execute the progarm:

$ cd app/cmd/printer-timeline-backend && go run main.go

* Configuration

The API is using AWS Services under the hood. It is using the AWS SDK for the Go programming language [AWS SDK for Go](https://aws.amazon.com/sdk-for-go/).

AWS uses different regions and depending where the AWS Services are. You will need to change the environment variables to match the AWS Service's regions (the file dev.env already contains the correct ones). Check the AWS console (web) to see the regions of your services.

The API uses AWS Cloudwatch Insights: AWS CloudWatch is in the region 'us-east-1'. 
The Cloud Connector bucket 'cloudconnector-core-production' is in us-east-1 and the bucket 'cloudconnector-to-blacksea-production' is in us-west-1. 

This API is intended to work with the account 'latex-dev' so you will need an IAM user and the corresponding AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY for that account.
The setup instructions for the credentials and the AWS CLI are here: [AWS Credentials Configuration](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html).

For more info about AWS accounts contact the Cloud Connector team.

* How to run tests

cd app/internal  && go test ./... -cover;

* Deployment instructions

In order to deploy the program you will need the correspnding AWS credentials and the AWS CLI (more info in previous section).

- Step 0: Set env variable DEVELOPMENT=false in dev.env

- Step 1: Build (compile the program)

$ GOOS=linux GOARCH=amd64 go build -o lambda cmd/printer-timeline-backend/main.go

- Step 2: Compress the binary

$ zip -m lambda.zip lambda

- Step 3a: Create the AWS Lambda (ONLY IF THE AWS LAMBDA does not exist):

$ aws lambda create-function --region us-east-1 --function-name printer-timeline-backend --runtime go1.x \
  --zip-file fileb://lambda.zip --handler lambda --timeout 300 --memory-size 128 \
  --role arn:aws:iam::103691013827:role/lambda_basic_execution \
  --description "Lambda that returns the last json of particular printer" \
  --environment "Variables={$(grep = ../dev.env | xargs | tr " " ",")}"

- Step 3b: Update the AWS Lambda with the zip containing the new binary created in previous steps (here you are uploading the new version of the code to AWS).

$ aws lambda update-function-code --region us-east-1 --function-name printer-timeline-backend --zip-file fileb://lambda.zip --publish

### Who do I talk to? ###

* Cloud Connector team
* Aldo Fuster Turpin
