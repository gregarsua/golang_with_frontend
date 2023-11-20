package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
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
	// router.HandleFunc("/user/{id}", makeHTTPHandleFunc(s.handleGetUserById))

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
// func (s *APIServer) handleGetUserById(w http.ResponseWriter, r *http.Request) error {
// 	switch r.Method {
// 	case "GET":
// 		id, err := getID(r)
// 		if err != nil {
// 			return err
// 		}

// 		user, err := s.store.GetUsersById(id)
// 		if err != nil {
// 			return nil
// 		}

// 		return WriteJSON(w, http.StatusOK, user)
// 	default:
// 		return fmt.Errorf("method not allowed %s", r.Method)
// 	}
// }

// POST a user on the database
func (s *APIServer) handleCreateUser(w http.ResponseWriter, r *http.Request) error {
	return nil
}
