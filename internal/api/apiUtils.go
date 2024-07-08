package api

import (
	"encoding/json"
	"net/http"
)

type FuncAPI func(w http.ResponseWriter, r *http.Request) error

func MakeHTTPHandleFunc(fn FuncAPI) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, err.Error())
		}
	}
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}