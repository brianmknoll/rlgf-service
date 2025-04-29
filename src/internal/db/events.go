package db

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type DbEvent struct {
	EventId string `firestore:"eventId"`
	Name    string `firestore:"name"`
}

func (f *FirestoreDatabase) CreateEvent(guildId, name string) error {
	// TODO: Choose something more uniuque than just the name.
	eventRef := f.client.Collection("guilds").Doc(guildId).Collection("events").Doc(name)
	wr, err := eventRef.Create(context.Background(), DbEvent{
		EventId: uuid.New().String(),
		Name:    name,
	})
	if err != nil {
		return err
	}
	fmt.Println(wr)
	return nil
}
