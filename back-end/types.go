package main

import (
	"back-end/storage"
	"net/http"
)

type APIServer struct {
	listenAddr string
	store      storage.Storage
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}
