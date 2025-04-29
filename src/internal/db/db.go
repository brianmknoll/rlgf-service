package db

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
)

const projectId = "real-life-group-finder"

type Database interface {
	CreateEvent(guildId, name string) error
}

type FirestoreDatabase struct {
	client *firestore.Client
}

func NewFirestoreDatabase() Database {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectId)
	if err != nil {
		panic(fmt.Sprintf("Failed to create Firestore client: " + err.Error()))
	}
	return &FirestoreDatabase{client: client}
}
