package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/users", makeHTTPHandleFunc(s.handleUsers))
	router.HandleFunc("/users/{id}", makeHTTPHandleFunc(s.handleGetUserById))

	log.Println("JSON API server running on port:", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

// MAIN function for handling user inquiries to the database
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

	return WriteJSON(w, http.StatusOK, users)
}

// GET user by using an id
func (s *APIServer) handleGetUserById(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		id, err := getID(r)
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

		return WriteJSON(w, http.StatusOK, user)
	case "DELETE":
		return s.handleDeleteUser(w, r)
	default:
		return fmt.Errorf("method not allowed %s", r.Method)
	}
}

// POST a user on the database
func (s *APIServer) handleCreateUser(w http.ResponseWriter, r *http.Request) error {
	var createUserRequest CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&createUserRequest); err != nil {
		return fmt.Errorf("error decoding %v", err)
	}

	user := &User{
		ID:          primitive.NewObjectID(),
		FirstName:   createUserRequest.FirstName,
		LastName:    createUserRequest.LastName,
		Company:     createUserRequest.Company,
		PhoneNumber: createUserRequest.PhoneNumber,
	}

	if err := s.store.CreateUser(user); err != nil {
		return fmt.Errorf("error creating account %v", err)
	}

	return WriteJSON(w, http.StatusOK, user)
}

// UPDATE a user using id
// func (s *APIServer) handleUpdateUser(w http.ResponseWriter, r *http.Request) error {
// 	return nil
// }

// DELETE a user by id
func (s *APIServer) handleDeleteUser(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
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

	return WriteJSON(w, http.StatusOK, map[string]string{"message": "User Deleted"})
}
