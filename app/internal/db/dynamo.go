package db

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/pkg/errors"
)

type PrinterSubscriptionFetcher interface {
	GetPrinterSubscriptions(ctx context.Context, printerId string) ([]*CCPrinterSubscriptionModel, error)
}

type CCPrinterSubscriptionCollection struct {
	*dynamodb.DynamoDB
	tableName string
}

// NewCCPrinterSubscriptionCollection configures a Collection to connect to dynamo CCPrinterSubscription table
// given a session and based on the environment variable table name
func NewCCPrinterSubscriptionCollectionWithSession(s *session.Session) (*CCPrinterSubscriptionCollection, error) {
	envVar := "TABLE_CC_PRINTER_SUBSCRIPTION"
	tableName, exist := os.LookupEnv(envVar)
	if !exist {
		return nil, errors.Errorf("you have to define the environment variable %s to work with dynamo", envVar)
	}

	return &CCPrinterSubscriptionCollection{
		DynamoDB:  dynamodb.New(s),
		tableName: tableName,
	}, nil
}

// Put stores a subscription in dynamo
func (col *CCPrinterSubscriptionCollection) Put(ctx context.Context, subscription *CCPrinterSubscriptionModel) error {

	newSubscription, err := dynamodbattribute.MarshalMap(subscription)
	if err != nil {
		return fmt.Errorf("error marshalling printer. cause: %w", err)
	}
	condition := aws.String("attribute_not_exists(PrinterID) AND attribute_not_exists(AccountID)")

	input := &dynamodb.PutItemInput{
		Item:                newSubscription,
		ConditionExpression: condition,
		TableName:           aws.String(col.tableName),
	}

	_, err = col.PutItemWithContext(ctx, input)

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			switch awsErr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				return ConditionalPutErr
			default:
				return err
			}
		}
	}
	return nil
}

// GetPrinterSubscriptions queries all subscriptions in dynamo for the specified printerID
func (col *CCPrinterSubscriptionCollection) GetPrinterSubscriptions(ctx context.Context, printerId string) ([]*CCPrinterSubscriptionModel, error) {

	queryInput := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":printerId": {
				S: aws.String(printerId),
			},
		},
		KeyConditionExpression: aws.String("PrinterID = :printerId"),
		TableName:              aws.String(col.tableName),
	}

	result, err := col.QueryWithContext(ctx, queryInput)
	if err != nil {
		return nil, err
	}
	if len(result.Items) == 0 {
		return nil, NotFoundErr
	}

	var subscriptionsArr []*CCPrinterSubscriptionModel

	for _, item := range result.Items {
		subscription := CCPrinterSubscriptionModel{}
		err := dynamodbattribute.UnmarshalMap(item, &subscription)
		if err != nil {
			return nil, err
		}
		subscriptionsArr = append(subscriptionsArr, &subscription)
	}

	return subscriptionsArr, nil
}

// Get retrieves the subscription from dynamo for the specified printerID and accountID
func (col *CCPrinterSubscriptionCollection) Get(ctx context.Context, printerId string, accountId string) (CCPrinterSubscriptionModel, error) {
	subscription := CCPrinterSubscriptionModel{}

	result, err := col.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(col.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"PrinterID": {
				S: aws.String(printerId),
			},
			"AccountID": {
				S: aws.String(accountId),
			},
		},
	})

	if err != nil {
		return CCPrinterSubscriptionModel{}, err
	}

	if len(result.Item) == 0 {
		return CCPrinterSubscriptionModel{}, NotFoundErr
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &subscription)
	if err != nil {
		return CCPrinterSubscriptionModel{}, fmt.Errorf("error unmarshalling printer subscription. cause: %w", err)
	}

	return subscription, err
}

// Delete removes a row of CCPrinterSubscriptions table based on the printerID and accountID specified
func (col *CCPrinterSubscriptionCollection) Delete(ctx context.Context, printerId string, accountId string) error {

	condition := aws.String("attribute_exists(PrinterID) AND attribute_exists(AccountID)")
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(col.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"PrinterID": {
				S: aws.String(printerId),
			},
			"AccountID": {
				S: aws.String(accountId),
			},
		},
		ConditionExpression: condition,
	}

	_, err := col.DeleteItemWithContext(ctx, input)

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			switch awsErr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				return ConditionalDelErr
			default:
				return err
			}
		}
	}
	return nil
}
