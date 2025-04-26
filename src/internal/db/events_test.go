package db

import (
	"context"
	"net/url"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	transport "github.com/aws/smithy-go/endpoints"
)

const (
	localDynamoDBEndpoint = "http://localhost:8000"
	awsRegion             = "us-west-1"
)

type localResolver struct{}

func (localResolver) ResolveEndpoint(ctx context.Context, params dynamodb.EndpointParameters) (transport.Endpoint, error) {
	u, _ := url.Parse("http://localhost:8000")
	return transport.Endpoint{URI: *u}, nil
}

func newTestDynamoDatabase() (Database, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if err != nil {
		return nil, err
	}
	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.EndpointResolverV2 = localResolver{}
	})
	return &DynamoDatabase{
		dyn: client,
	}, nil
}

func TestCreateEvent(t *testing.T) {
	db, err := newTestDynamoDatabase()
	if err != nil {
		t.Fatalf("Failed to create test DynamoDB client: %v", err)
	}

	err = db.CreateEvent("Test Event")
	if err != nil {
		t.Fatalf("Failed to create event: %v", err)
	}

	// Verify the event was created (this part would require a read operation)
}
