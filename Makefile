test:
	go version


build:
	echo "Compiling..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o lambda cmd/printer-timeline-backend/main.go


zip: 
	zip -m lambda.zip lambda


redeploy:
	aws lambda update-function-code --region us-east-1 --function-name printer-timeline-backend --zip-file fileb://lambda.zip --publish


all: 
	compile zip redeploy
