package main

import (
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Storage interface {
	GetUsers() ([]*User, error)
	// GetUsersById(int) (*Users, error)
}

type APIServer struct {
	listenAddr string
	store      Storage
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

type MongoStore struct {
	client *mongo.Client
	db     *mongo.Database
}

type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	FirstName   string             `bson:"firstName"`
	LastName    string             `bson:"lastName"`
	Company     string             `bson:"company"`
	PhoneNumber int64              `bson:"phoneNumer"`
}
