package response

import (
	"encoding/json"
	"net/http"
)

type Success struct {
	Data any `json:"data"`
}

type Error struct {
	Error string `json:"error"`
}

func JsonResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Success{Data: data})
}

func ErrorResponse(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Error{Error: msg})
}
