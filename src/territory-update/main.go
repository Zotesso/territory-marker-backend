package main

import (
	"os"
	"log"
	"github.com/reecerussell/aws-lambda-multipart-parser/parser"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/aws/session"
)

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

	territoryRequestItems, err := parser.Parse(request)

	if err != nil {
		log.Fatalf("Got error calling UpdateItem: %s", err)
	}

	last_worked_date, ok := territoryRequestItems.Get("last_worked_date")

	if !ok  {
		resp.StatusCode = 400
		resp.Body = string("Missing field last_worked_date")
		return &resp, nil
	}

	last_given_date, ok := territoryRequestItems.Get("last_given_date")

	if !ok  {
		resp.StatusCode = 400
		resp.Body = string("Missing field last_given_date")
		return &resp, nil
	}

	owner, ok := territoryRequestItems.Get("owner")

	if !ok  {
		resp.StatusCode = 400
		resp.Body = string("Missing field owner")
		return &resp, nil
	}

	// Update item in table Movies
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":r": {
				S: aws.String(last_worked_date),
			},
			":c": {
				S: aws.String(last_given_date),
			},
			":d": {
				S: aws.String(owner),
			},
		},
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
      "id": {
          N: aws.String("1"),
      },
		},
		ReturnValues: aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set last_worked_date= :r, last_given_date= :c, owner= :d"),
	}

	_, err = dynaClient.UpdateItem(input)

	if err != nil {
			log.Fatalf("Got error calling UpdateItem: %s", err)
	}

	resp.StatusCode = 200
	resp.Body = string(request.Body)
	return &resp, nil
}