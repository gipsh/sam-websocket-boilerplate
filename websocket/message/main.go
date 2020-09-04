package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
)

// RequestPayload represents the request body sent by the Socket

type ConnectionItem struct {
	ConnectionID string   `json:"ConnectionID"`
	Created      string   `json:"Created"`
	SrcIP        string   `json:"SrcIP"`
	Messages     []string `json:"Messages"`
}

func main() {
	lambda.Start(HandleConnect)
}

func updateConnection(ctx context.Context, item ConnectionItem, message string) {

	config, _ := external.LoadDefaultAWSConfig()

	dynamodbSession := dynamodb.New(config)

	item.Messages = append(item.Messages, message)

	attributeValues, _ := dynamodbattribute.MarshalMap(item)

	putInput := &dynamodb.PutItemInput{
		Item:      attributeValues,
		TableName: aws.String(os.Getenv("tableName")),
	}

	putReq := dynamodbSession.PutItemRequest(putInput)

	_, putErr := putReq.Send(ctx)
	if putErr != nil {
		fmt.Println(putErr)
	}

}

func findConnection(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) ConnectionItem {

	config, _ := external.LoadDefaultAWSConfig()

	dynamodbSession := dynamodb.New(config)

	input := &dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("tableName")),
		Key: map[string]dynamodb.AttributeValue{
			"ConnectionID": {
				S: aws.String(request.RequestContext.ConnectionID),
			},
		},
	}

	req := dynamodbSession.GetItemRequest(input)
	resp, err := req.Send(context.TODO())
	if err == nil {
		fmt.Println(resp)
	}

	item := ConnectionItem{}
	err = dynamodbattribute.UnmarshalMap(resp.Item, &item)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	return item

}

func responseToConnection(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) {

	// TODO: use an external env variable to override the endpoint
	config, _ := external.LoadDefaultAWSConfig(request.RequestContext.DomainName)
	config.EndpointResolver = aws.ResolveWithEndpointURL(
		fmt.Sprintf("https://%s/%s",
			request.RequestContext.DomainName,
			request.RequestContext.Stage))

	apigw := apigatewaymanagementapi.New(config)

	input := &apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: aws.String(request.RequestContext.ConnectionID),
		Data:         []byte("hello from server"),
	}

	req := apigw.PostToConnectionRequest(input)

	_, err := req.Send(ctx)
	if err != nil {
		fmt.Println(err)
	}
}

func HandleConnect(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {

	// response to the message
	responseToConnection(ctx, request)

	// find the connection on dynamdb
	connectioItem := findConnection(ctx, request)

	// and update the object with the new message
	updateConnection(ctx, connectioItem, request.Body)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "OK",
	}, nil

}
