package db

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type DbMessage struct {
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

func (f *FirestoreDatabase) CreateMessage(guildId, channelId string, message DbMessage) error {
	eventRef := f.client.
		Collection("guilds").
		Doc(guildId).
		Collection("channels").
		Doc(channelId).
		Collection("messages").
		Doc(message.Timestamp.Format("2006-01-02T15:04:05.000Z07:00"))
	wr, err := eventRef.Create(context.Background(), message)
	if err != nil {
		return err
	}
	fmt.Println(wr)
	return nil
}

func (f *FirestoreDatabase) ReadRecentMessages(guildId, channelId string) ([]DbMessage, error) {
	ctx := context.Background()

	recent := time.Now().Add(-12 * time.Hour)

	messagesCollection := f.client.
		Collection("guilds").
		Doc(guildId).
		Collection("channels").
		Doc(channelId).
		Collection("messages")

	query := messagesCollection.Where("timestamp", ">", recent).OrderBy("timestamp", firestore.Asc)

	iter := query.Documents(ctx)
	defer iter.Stop()

	var messages []DbMessage
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate over messages: %w", err)
		}

		var msg DbMessage
		if err := doc.DataTo(&msg); err != nil {
			fmt.Printf("Failed to convert document data to DbMessage for doc ID %s: %v\n", doc.Ref.ID, err)
			continue
		}
		messages = append(messages, msg)
	}

	return messages, nil
}
