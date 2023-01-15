package main

import (
	"os"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/aws/session"
)

type Territory struct {
	Name string `json:"name"`
}

const tableName = "territories"

var(
	dynaClient dynamodbiface.DynamoDBAPI
)

func main() {
	region := os.Getenv("AWS_REGION")
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},)

		if err != nil {
			return
		}
	dynaClient = dynamodb.New(awsSession)
	lambda.Start(handler)
}

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	resp := events.APIGatewayProxyResponse{Headers: map[string]string{"Content-Type":"application/json"}}

	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"name": {S: aws.String("teste")},
		},
		TableName: aws.String(tableName),
	}

	result, err := dynaClient.GetItem(input)

	if err != nil {
		return nil, err
	}

	item := new(Territory)
	err = dynamodbattribute.UnmarshalMap(result.Item, item)

	if err != nil {
		return nil, err
	}

	stringBody, _ := json.Marshal(item)

	resp.StatusCode = 200
	resp.Body = string(stringBody)
	return &resp, nil
}