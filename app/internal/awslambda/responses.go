package awslambda

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func newLambdaOkResponse(response []byte) (*events.APIGatewayProxyResponse, error) {
	headers := map[string]string{
		"Content-type": "text",
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       fmt.Sprintf("%s", response),
		Headers:    headers,
	}, nil
}

func newLambdaError(httpStatus int, err error) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		StatusCode: httpStatus,
		Body:       fmt.Sprintf("%s", err.Error()),
	}, nil
}
