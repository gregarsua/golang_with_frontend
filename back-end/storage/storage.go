package storage

import (
	"back-end/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Storage interface {
	GetUsers() ([]*types.User, error)
	GetUserByID(userID primitive.ObjectID) (*types.User, error)
	CreateUser(user *types.User) error
	UpdateUser(userID primitive.ObjectID, updatedUser *types.User) error
	DeleteUser(userID primitive.ObjectID) error
}
