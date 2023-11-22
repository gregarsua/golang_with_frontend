package main

import (
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Storage interface {
	GetUsers() ([]*User, error)
	GetUserByID(userID primitive.ObjectID) (*User, error)
	CreateUser(user *User) error
	UpdateUser(userID primitive.ObjectID, updatedUser *User) error
	DeleteUser(userID primitive.ObjectID) error
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
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FirstName   string             `json:"firstName" bson:"firstName"`
	LastName    string             `json:"lastName" bson:"lastName"`
	Company     string             `json:"company" bson:"company"`
	PhoneNumber int64              `json:"phoneNumber" bson:"phoneNumber"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type CreateUserRequest struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Company     string `json:"company"`
	PhoneNumber int64  `json:"phoneNumber"`
}
