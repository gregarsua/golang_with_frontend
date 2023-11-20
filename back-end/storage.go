package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// connect to the database
func NewMongoDbStore() (*MongoStore, error) {
	opts := options.Client().ApplyURI("mongodb+srv://gregarsua:1o5Jze45Y5EAHlLc@clustergolang.px9nxvw.mongodb.net/?retryWrites=true&w=majority")
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	db := client.Database("Users")

	return &MongoStore{
		client: client,
		db:     db,
	}, nil
}

// get the users from the users collection
func (s *MongoStore) GetUsers() ([]*User, error) {
	collection := s.db.Collection("users")

	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var users []*User
	for cursor.Next(context.TODO()) {
		var user User
		if err := cursor.Decode(&user); err != nil {
			log.Printf("Error decoding user document: %v", err)
			continue // Skip to the next document in case of decoding error
		}
		users = append(users, &user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
