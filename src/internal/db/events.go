package db

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

type DbEvent struct {
	EventId string `dynamodbav:"eventId"`
	Name    string `dynamodbav:"name"`
}

func (d *DynamoDatabase) CreateEvent(name string) error {
	item, err := dynamodbattribute.MarshalMap(DbEvent{
		EventId: uuid.New().String(),
		Name:    name,
	})
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String("events"),
	}

	_, err = d.dyn.PutItem(context.TODO(), input)
	if err != nil {
		return err
	}

	fmt.Printf("Event created: %s\n", name)
	return nil
}
