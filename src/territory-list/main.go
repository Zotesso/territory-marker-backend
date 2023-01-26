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
	Id int `json:"id"`
	Name string `json:"name"`
	Polylines string `json:"polylines"`
	Last_worked_date string `json:"last_worked_date"`
	Last_given_date string `json:"last_given_date"`
	Owner string `json:"owner"`
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
	resp := events.APIGatewayProxyResponse{Headers: map[string]string{
		"Content-Type":"application/json",
		"Access-Control-Allow-Origin":"*",
		"Access-Control-Allow-Credentials":"true",
	}}

	input := &dynamodb.ScanInput {
		TableName: aws.String(tableName),
	}

	result, err := dynaClient.Scan(input)

	if err != nil {
		return nil, err
	}

	item := []Territory{}
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &item)

	if err != nil {
		return nil, err
	}

	stringBody, _ := json.Marshal(item)

	resp.StatusCode = 200
	resp.Body = string(stringBody)
	return &resp, nil
}