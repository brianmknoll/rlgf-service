package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"fmt"
	"slices"
)

const (
	region          = "us-west-1"
	migrationsTable = "migrations"
	migrationId     = "migrationId"
)

type Migration struct {
	migrationId string
	job         func(*dynamodb.DynamoDB)
}

var migrations = []Migration{
	{
		migrationId: "create_events_table",
		job:         createEventTable,
	},
}

func (db *DynamoDatabase) RunMigrations() {
	describeInput := &dynamodb.DescribeTableInput{
		TableName: aws.String(migrationsTable),
	}

	_, err := db.dyn.DescribeTable(describeInput)
	if err != nil {
		createMigrationsTable(db.dyn)
	}

	migrationIds, err := getMigrationIdsAlreadyRan(db.dyn)
	if err != nil {
		panic(fmt.Sprintf("Failed to get migration IDs: %v", err))
	}

	for _, migration := range migrations {
		if !slices.Contains(migrationIds, migration.migrationId) {
			fmt.Printf("Running migration: %s\n", migration.migrationId)
			migration.job(db.dyn)

			err := markMigrationAsRan(db.dyn, migration.migrationId)
			if err != nil {
				panic(fmt.Sprintf("Failed to mark migration as ran: %v", err))
			}
		} else {
			fmt.Printf("Migration %s already ran, skipping...\n", migration.migrationId)
		}
	}
}

func markMigrationAsRan(dyn *dynamodb.DynamoDB, migrationId string) error {
	input := &dynamodb.PutItemInput{
		TableName: aws.String(migrationsTable),
		Item: map[string]*dynamodb.AttributeValue{
			migrationId: {
				S: aws.String(migrationId),
			},
		},
	}

	_, err := dyn.PutItem(input)
	if err != nil {
		return fmt.Errorf("failed to mark migration as ran: %w", err)
	}

	fmt.Printf("Migration %s marked as ran\n", migrationId)
	return nil
}

func getMigrationIdsAlreadyRan(dyn *dynamodb.DynamoDB) ([]string, error) {
	var allMatchingItems []map[string]*dynamodb.AttributeValue

	var lastEvaluatedKey map[string]*dynamodb.AttributeValue
	scanPageCount := 0

	fmt.Println("Scanning migration table...")

	for {
		scanPageCount++
		fmt.Printf("Scanning page %d...\n", scanPageCount)

		scanInput := &dynamodb.ScanInput{
			TableName:            aws.String(migrationsTable),
			ProjectionExpression: aws.String(migrationId),
			ExclusiveStartKey:    lastEvaluatedKey,
		}

		output, err := dyn.Scan(scanInput)
		if err != nil {
			panic(fmt.Sprintf("Scan failed: %v", err))
		}

		if len(output.Items) > 0 {
			allMatchingItems = append(allMatchingItems, output.Items...)
		}

		lastEvaluatedKey = output.LastEvaluatedKey
		if lastEvaluatedKey == nil {
			break
		}
	}

	fmt.Printf("Scan complete. Found %d matching items.\n", len(allMatchingItems))

	return unmarshalMigrationIds(allMatchingItems)
}

type MigrationRow struct {
	MigrationId string `json:"migrationId"`
}

func unmarshalMigrationIds(rows []map[string]*dynamodb.AttributeValue) ([]string, error) {
	var unmarshalledRows []MigrationRow
	var migrationIds []string

	err := dynamodbattribute.UnmarshalListOfMaps(rows, &unmarshalledRows)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal items: %v", err))
	}

	for _, row := range unmarshalledRows {
		migrationIds = append(migrationIds, row.MigrationId)
	}

	return migrationIds, nil
}

func createMigrationsTable(dyn *dynamodb.DynamoDB) {
	input := &dynamodb.CreateTableInput{
		TableName: aws.String(migrationsTable),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String(migrationId),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String(migrationId),
				KeyType:       aws.String("HASH"),
			},
		},
		BillingMode: aws.String(dynamodb.BillingModePayPerRequest),
	}

	_, err := dyn.CreateTable(input)
	if err != nil {
		panic(fmt.Sprintf("Failed to create table: %v", err))
	}

	fmt.Println("migrations table created successfully")
}

func createEventTable(dyn *dynamodb.DynamoDB) {
	input := &dynamodb.CreateTableInput{
		TableName: aws.String("events"),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("eventId"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("eventId"),
				KeyType:       aws.String("HASH"),
			},
		},
		BillingMode: aws.String(dynamodb.BillingModePayPerRequest),
	}

	_, err := dyn.CreateTable(input)
	if err != nil {
		panic(fmt.Sprintf("Failed to create table: %v", err))
	}

	fmt.Println("events table created successfully")
}
