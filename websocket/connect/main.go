package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
)

type ConnectionItem struct {
	// ID  string `json:"connectionID"`
	ConnectionID string   `json:"ConnectionID"`
	Created      string   `json:"Created"`
	SrcIP        string   `json:"SrcIP"`
	Messages     []string `json:"Messages"`
}

func main() {
	lambda.Start(HandleConnect)
}

func HandleConnect(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Create a representation of the Row in the table
	fmt.Println("HOLA!!!")

	ms := make([]string, 0)

	connectionItem := ConnectionItem{
		//         ID: request.RequestContext.ConnectionID,
		ConnectionID: request.RequestContext.ConnectionID,
		Created:      time.Now().String(),
		SrcIP:        request.RequestContext.Identity.SourceIP,
		Messages:     ms,
	}
	attributeValues, _ := dynamodbattribute.MarshalMap(connectionItem)

	fmt.Println(attributeValues)

	input := &dynamodb.PutItemInput{
		Item:      attributeValues,
		TableName: aws.String(os.Getenv("tableName")),
	}

	// Insert the Row in the table
	config, _ := external.LoadDefaultAWSConfig()

	dynamodbSession := dynamodb.New(config)

	req := dynamodbSession.PutItemRequest(input)

	_, err := req.Send(ctx)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 200, Body: "ERROR"}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "connected",
	}, nil

}
