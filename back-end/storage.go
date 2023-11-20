package main

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// connect to the database
func NewMongoDbStore() (*MongoStore, error) {
	mongoURI := os.Getenv("MONGO_URI")
	opts := options.Client().ApplyURI(mongoURI)
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

func (s *MongoStore) GetUserByID(userID primitive.ObjectID) (*User, error) {
	collection := s.db.Collection("users")

	filter := bson.D{{Key: "_id", Value: userID}}

	var user User
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *MongoStore) CreateUser(user *User) error {
	collection := s.db.Collection("users")

	_, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Printf("error inserting data %v", err)
		return err
	}

	return nil
}

func (s *MongoStore) UpdateUser(user *User) error {
	return nil
}

func (s *MongoStore) DeleteUser(user *User) error {
	return nil
}
