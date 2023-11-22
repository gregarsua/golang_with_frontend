package storage

import (
	"back-end/types"
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStore struct {
	client *mongo.Client
	db     *mongo.Database
}

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

// post user to database
func (s *MongoStore) CreateUser(user *types.User) error {
	collection := s.db.Collection("users")

	_, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Printf("error inserting data %v", err)
		return err
	}

	return nil
}

// get the users from the users collection
func (s *MongoStore) GetUsers() ([]*types.User, error) {
	collection := s.db.Collection("users")

	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var users []*types.User
	for cursor.Next(context.TODO()) {
		var user types.User
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

// get the user in the db using id
func (s *MongoStore) GetUserByID(userID primitive.ObjectID) (*types.User, error) {
	collection := s.db.Collection("users")

	filter := bson.D{{Key: "_id", Value: userID}}

	var user types.User
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// delete user from database
func (s *MongoStore) DeleteUser(userID primitive.ObjectID) error {
	collection := s.db.Collection("users")

	_, err := collection.DeleteOne(context.TODO(), bson.M{"_id": userID})
	if err != nil {
		log.Printf("id not found %v", err)
		return err
	}

	return nil
}

// update user from database
func (s *MongoStore) UpdateUser(userID primitive.ObjectID, updatedUser *types.User) error {
	collection := s.db.Collection("users")

	filter := bson.M{"_id": userID}
	update := bson.M{
		"$set": bson.M{
			"firstName":   updatedUser.FirstName,
			"lastName":    updatedUser.LastName,
			"company":     updatedUser.Company,
			"phoneNumber": updatedUser.PhoneNumber,
			"updatedAt":   time.Now().UTC(),
		},
	}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Printf("error updating user %v", err)
		return err
	}

	return nil
}
