test:
	go test ./...

build:
	echo "Compiling..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o lambda cmd/printer-timeline-backend/main.go

zip: 
	zip -m lambda.zip lambda

create-function:
	aws lambda create-function --region us-east-1 --function-name printer-timeline-backend --runtime go1.x \
	--zip-file fileb://lambda.zip --handler lambda --timeout 900 --memory-size 128 \
	--role arn:aws:iam::103691013827:role/lambda_basic_execution \
	--environment "Variables={$$(grep = ../dev.env | xargs | tr " " ",")}" \
	--description "lambda that retrieves data related to printers"

update-function:
	aws lambda update-function-code --region us-east-1 --function-name printer-timeline-backend --zip-file fileb://lambda.zip --publish

clean:
	rm lambda.zip

run-local:
	export $$(grep = 'dev.env' | xargs); go run cmd/printer-timeline-backend/main.go

first-all: build zip create-function clean

all: build zip update-function clean