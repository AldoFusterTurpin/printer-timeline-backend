package awslambda

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func newLambdaOkResponse(response []byte) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       fmt.Sprintf("%s", response),
	}, nil
}

func newLambdaError(httpStatus int, err error) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		StatusCode: httpStatus,
		Body:       fmt.Sprintf("%s", err.Error()),
	}, nil
}
