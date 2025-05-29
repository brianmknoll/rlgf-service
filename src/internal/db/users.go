package db

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DbUser struct {
	Id       int       `firestore:"id"`
	Username string    `firestore:"username"`
	Name     string    `firestore:"name"`
	Joined   time.Time `firestore:"joined"`
}

func (f *FirestoreDatabase) ReadUsers(guildId string) ([]DbUser, error) {
	ctx := context.Background()

	usersRef := f.client.
		Collection("guilds").
		Doc(guildId).
		Collection("users")

	iter := usersRef.Documents(ctx)
	defer iter.Stop()

	var users []DbUser
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate over users: %w", err)
		}

		var user DbUser
		if err := doc.DataTo(&user); err != nil {
			fmt.Printf("Failed to convert document data to DbUser for doc ID %s: %v\n", doc.Ref.ID, err)
			continue
		}
		users = append(users, user)
	}

	return users, nil
}

func (f *FirestoreDatabase) CreateUser(guildId string, user DbUser) error {
	userRef := f.client.
		Collection("guilds").
		Doc(guildId).
		Collection("users").
		Doc(fmt.Sprintf("%d", user.Id))

	wr, err := userRef.Create(context.Background(), user)

	if err != nil {
		if st, _ := status.FromError(err); st.Code() == codes.AlreadyExists {
			return ErrAlreadyExists
		}
		return err
	}

	fmt.Println(wr)
	return nil
}
