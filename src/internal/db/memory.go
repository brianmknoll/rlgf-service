package db

import (
	"context"
	"fmt"
)

type DbMemory struct {
	Memory string `json:"memory"`
}

const memoryDocId = "omni-memory"

func (f *FirestoreDatabase) ReadMemories(guildId string) (string, error) {
	memoryRef := f.client.
		Collection("guilds").
		Doc(guildId).
		Collection("memories").
		Doc(memoryDocId)

	doc, err := memoryRef.Get(context.Background())
	if err != nil {
		return "", err
	}
	if !doc.Exists() {
		return "", nil
	}

	var mem DbMemory
	err = doc.DataTo(&mem)
	if err != nil {
		return "", err
	}
	return mem.Memory, nil
}

func (f *FirestoreDatabase) CreateMemory(guildId, memory string) error {
	memoryRef := f.client.
		Collection("guilds").
		Doc(guildId).
		Collection("memories").
		Doc(memoryDocId)

	newMemory := memory

	doc, err := memoryRef.Get(context.Background())
	if err != nil {
		return err
	}
	if doc.Exists() {
		var mem DbMemory
		err := doc.DataTo(&mem)
		if err != nil {
			return err
		}
		newMemory = mem.Memory + "\n" + memory
	}

	if doc.Exists() {
		_, err := memoryRef.Set(context.Background(), DbMemory{
			Memory: newMemory,
		})
		if err != nil {
			return err
		}
		fmt.Println("Updated memory document")
	} else {
		_, err := memoryRef.Create(context.Background(), DbMemory{
			Memory: newMemory,
		})
		if err != nil {
			return err
		}
		fmt.Println("Created new memory document")
	}
	return nil
}
