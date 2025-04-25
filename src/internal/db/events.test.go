package db

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
)

func TestCreateEvent(t *testing.T) {
	sess := session.Must(session.NewSession())
	db := NewDynamoDatabase(sess)

	err := db.CreateEvent("Test Event")
	if err != nil {
		t.Fatalf("Failed to create event: %v", err)
	}

	// Verify the event was created (this part would require a read operation)
}
