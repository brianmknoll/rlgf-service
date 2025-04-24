package main

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Server interface {
	eventsPostHandler(w http.ResponseWriter, r *http.Request)
	runMigrations()
}

type RlgfServer struct {
	svc *dynamodb.DynamoDB
}

func NewRlgServer() *RlgfServer {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := dynamodb.New(sess)
	return &RlgfServer{
		svc: svc,
	}
}
