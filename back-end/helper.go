package main

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

// func getID(r *http.Request) (int, error) {
// 	idStr := mux.Vars(r)["id"]
// 	id, err := strconv.Atoi(idStr)

// 	if err != nil {
// 		return id, fmt.Errorf("invalid id %s", idStr)
// 	}

// 	return id, nil
// }

// func scanIntoAccount(rows *sql.Rows) (*Users, error) {
// 	account := new(Users)
// 	err := rows.Scan(
// 		&account.ID,
// 		&account.FirstName,
// 		&account.LastName,
// 		&account.Company)

// 	return account, err
// }
