# README #
To run the API, execute:

export AWS_REGION=us-east-1 && export AWS_REGION_BLACK_SEA_BUCKET=us-west-1 && go run main.go.

You can find samples of the REQUESTS in the Postman file: printer-timeline-backend.postman_collection.json at the root level of the project. 
Time type can be relative or absolute. You can find more information on the request parameters in app/internal/queryparams/query_params_test.go. 

package queryparams extracts the query parameters common to more than one endpoint.
In app/internal/api/api.go you will find the different endpoints and the corresponding handlers.

You can follow a top-down approach to better understand the code. Also take a look at the test files.

### What is this repository for? ###

* It is the back-end of the Printer-timeline. 
This project is a Go API that let's you retrieve in chronological order different data elements of the Cloud Connector (OpenXMls, CloudJsons, etc.) specifying a time range and a pn (Product Number) and sn(Serial number), if you don't specify them you will obtain all the pn's and sn's in that time range. It also has an endpoint to retrieve AWS S3 objects by bucket name and object key.

### How do I get set up? ###

* Summary of set up
Export the environment variables in dev.env.
Execute the app: cd app && go run main.go.

* Configuration
The API is using AWS Services under the hood. It is using the AWS SDK for the Go programming language [AWS SDK for Go](https://aws.amazon.com/sdk-for-go/).

The file dev.env contains the env variables.

AWS uses different regions and depending where the AWS Services are. You will need to change the environment variables to match the AWS Service's regions (the file dev.env already contains the correct ones).s 

The API uses AWS Cloudwatch Insights: AWS CloudWatch is in the region 'us-east-1'. 
The Cloud Connector bucket 'cloudconnector-core-production' is in us-east-1 and the bucket 'cloudconnector-to-blacksea-production' is in us-west-1. 

The env var AWS_PROFILE is not mandatory as is used when you have multiple profiles and CT-PRODUCTION is the name of your profile that has the credentials of the account 'latex-dev'. You just need to use it if you have multiple profiles.

This API is intended to work with the account 'latex-dev' so you will need an IAM user and the corresponding AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY for that account.
You don't need to use the AWS CLI but the setup instructions are here (and the same credentials are needed for using the AWS CLI or the AWS SDK): [AWS Credentials Configuration](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html).

For more info about AWS accounts contact the Cloud Connector team.

* How to run tests
cd app && go test ./... -cover

* Deployment instructions
TODO. It will be an AWS Lambda.

### Who do I talk to? ###

* Cloud Connector team
* Aldo Fuster Turpin
