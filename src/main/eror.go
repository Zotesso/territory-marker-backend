package main

import (
	"log"
	"encoding/json"
	"errors"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Territorey struct {
	name string
}

func maine() {
	lambda.Start(handler)
}

func handlere(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("Function Invoked")
  cfg, err := config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
		o.Region = "us-east-1"
		return nil
	})

	if err != nil {
		panic(err)
	}

	svc := dynamodb.NewFromConfig(cfg)

	result, err := svc.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("territories"),
			Key: map[string]types.AttributeValue{
				"name": &types.AttributeValueMemberS{Value: "teste"},
			},
			AttributesToGet: []string{"name"},
	})

	log.Println("resultado -")
	log.Println(result.Item)

	if err != nil {
		log.Fatalf("Got error calling GetItem: %s", err)
	}

	if result == nil {
		msg := "Could not find Item"
		return events.APIGatewayProxyResponse{}, errors.New(msg)
	}

	items := Territory{}
	err = attributevalue.UnmarshalMap(result.Item, &items)
	log.Println("items -")
	log.Println(items)
	if err != nil {
		panic(err)
	}

	itemJson, err := json.Marshal(items)
	log.Println("itemJson -")
	log.Println(itemJson)
	if err != nil {
		panic(err)
	}

  return events.APIGatewayProxyResponse{Body: string(itemJson), StatusCode: 200}, nil
}