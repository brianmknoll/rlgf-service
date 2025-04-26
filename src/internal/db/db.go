package db

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws/session"
)

type Database interface {
	CreateEvent(name string) error
}

type DynamoDatabase struct {
	dyn *dynamodb.Client
}

func NewDynamoDatabase(s *session.Session) Database {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}
	dyn := dynamodb.NewFromConfig(cfg)
	return &DynamoDatabase{
		dyn: dyn,
	}
}
