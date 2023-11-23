package api

import (
	"back-end/helpers"
	"back-end/storage"
	"back-end/types"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type APIServer struct {
	listenAddr string
	store      storage.Storage
}

func NewAPIServer(listenAddr string, store storage.Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/users", helpers.MakeHTTPHandleFunc(s.handleUsers))
	router.HandleFunc("/users/{id}", helpers.MakeHTTPHandleFunc(s.handleGetUserById))

	log.Println("JSON API server running on port:", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleUsers(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetUsers(w, r)
	case "POST":
		return s.handleCreateUser(w, r)
	default:
		return fmt.Errorf("method not allowed %s", r.Method)
	}
}

// GET all users on the database
func (s *APIServer) handleGetUsers(w http.ResponseWriter, r *http.Request) error {
	users, err := s.store.GetUsers()
	if err != nil {
		return err
	}

	return helpers.WriteJSON(w, http.StatusOK, users)
}

// GET user by using an id
func (s *APIServer) handleGetUserById(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		id, err := helpers.GetID(r)
		if err != nil {
			return err
		}

		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return fmt.Errorf("error creating ObjectID: %v", err)
		}

		user, err := s.store.GetUserByID(objectID)
		if err != nil {
			return err
		}

		return helpers.WriteJSON(w, http.StatusOK, user)
	case "PUT":
		return s.handleUpdateUser(w, r)
	case "DELETE":
		return s.handleDeleteUser(w, r)
	default:
		return fmt.Errorf("method not allowed %s", r.Method)
	}
}

// POST a user on the database
func (s *APIServer) handleCreateUser(w http.ResponseWriter, r *http.Request) error {
	var createUserRequest types.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&createUserRequest); err != nil {
		return fmt.Errorf("error decoding %v", err)
	}

	user := &types.User{
		ID:          primitive.NewObjectID(),
		FirstName:   createUserRequest.FirstName,
		LastName:    createUserRequest.LastName,
		Company:     createUserRequest.Company,
		PhoneNumber: createUserRequest.PhoneNumber,
		CreatedAt:   time.Now().UTC(),
	}

	if err := s.store.CreateUser(user); err != nil {
		return fmt.Errorf("error creating account %v", err)
	}

	return helpers.WriteJSON(w, http.StatusOK, user)
}

// DELETE a user by id
func (s *APIServer) handleDeleteUser(w http.ResponseWriter, r *http.Request) error {
	id, err := helpers.GetID(r)
	if err != nil {
		return err
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("error creating objectid %v", err)
	}

	if err := s.store.DeleteUser(objectID); err != nil {
		return fmt.Errorf("error deleting user %v", err)
	}

	return helpers.WriteJSON(w, http.StatusOK, map[string]string{"message": "User Deleted"})
}

// UPDATE user by id
func (s *APIServer) handleUpdateUser(w http.ResponseWriter, r *http.Request) error {
	id, err := helpers.GetID(r)
	if err != nil {
		return err
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("error creating ObjectID: %v", err)
	}

	var updatedUser types.User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		return fmt.Errorf("error decoding request body: %v", err)
	}

	updatedUser.ID = objectID
	if err := s.store.UpdateUser(objectID, &updatedUser); err != nil {
		return fmt.Errorf("error updating user: %v", err)
	}

	return helpers.WriteJSON(w, http.StatusOK, map[string]string{"message": "User updated successfully"})
}
