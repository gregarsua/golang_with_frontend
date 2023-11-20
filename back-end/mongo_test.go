package main

import (
	"fmt"
	"testing"
)

func TestMongoDBIntegration(t *testing.T) {
	// Create a new MongoDB store
	mongoStore, err := NewMongoDbStore()
	if err != nil {
		t.Fatalf("Error creating MongoDB store: %v", err)
	}

	// Retrieve users from MongoDB
	users, err := mongoStore.GetUsers()
	if err != nil {
		t.Fatalf("Error retrieving users from MongoDB: %v", err)
	}

	// Print or log the retrieved users
	for _, user := range users {
		fmt.Printf("User ID: %v, FirstName: %v, LastName: %v\n", user.ID, user.FirstName, user.LastName)
	}

	// Optionally, add assertions to validate the results
	if len(users) == 0 {
		t.Error("Expected at least one user, but got none.")
	}
}
