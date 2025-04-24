package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Event struct {
	Name string `json:"name"`
}

func (s *RlgfServer) eventsPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var e Event

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&e)
	if err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = createEvent(s.svc, &e)
	if err != nil {
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("Received: Name=%s\n", e.Name)
	w.WriteHeader(http.StatusCreated)
}

type DbEvent struct {
	EventId string `json:"eventId"`
	Name    string `json:"name"`
}

func createEvent(svc *dynamodb.DynamoDB, e *Event) error {
	edb, err := dynamodbattribute.MarshalMap(e)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      edb,
		TableName: aws.String("events"),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		return err
	}

	fmt.Printf("Event created: %s\n", e.Name)
	return nil
}
