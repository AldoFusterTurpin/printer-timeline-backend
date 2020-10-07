package awslambda

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func newLambdaOkResponse(headers map[string]string, response []byte) (*events.APIGatewayProxyResponse, error) {
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
