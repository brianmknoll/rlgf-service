package db

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DynamoDatabase struct {
	dyn *dynamodb.DynamoDB
}

func NewDynamoDatabase(s *session.Session) *DynamoDatabase {
	dyn := dynamodb.New(s)
	return &DynamoDatabase{
		dyn: dyn,
	}
}
