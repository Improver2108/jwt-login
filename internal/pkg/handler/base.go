package handler

import (
	"encoding/json"
	"net/http"

	"github.com/improver2108/jwt-login/internal/pkg/response"
)

type Base[T any] struct {
	Service T
}

func (b *Base[T]) Decode(w http.ResponseWriter, r *http.Request, v any) bool {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return false
	}
	return true
}
func (b *Base[T]) OK(w http.ResponseWriter, data any) {
	response.JsonResponse(w, http.StatusOK, data)
}

func (b *Base[T]) Created(w http.ResponseWriter, data any) {
	response.JsonResponse(w, http.StatusCreated, data)
}

func (b *Base[T]) BadRequest(w http.ResponseWriter, msg string) {
	response.ErrorResponse(w, http.StatusBadRequest, msg)
}

func (b *Base[T]) Unauthorized(w http.ResponseWriter, msg string) {
	response.ErrorResponse(w, http.StatusUnauthorized, msg)
}

func (b *Base[T]) NotFound(w http.ResponseWriter, msg string) {
	response.ErrorResponse(w, http.StatusNotFound, msg)
}

func (b *Base[T]) Conflict(w http.ResponseWriter, msg string) {
	response.ErrorResponse(w, http.StatusConflict, msg)
}

func (b *Base[T]) InternalError(w http.ResponseWriter, msg string) {
	response.ErrorResponse(w, http.StatusInternalServerError, msg)
}

func (b *Base[T]) BadGateway(w http.ResponseWriter, msg string) {
	response.ErrorResponse(w, http.StatusBadGateway, msg)
}

func (b *Base[T]) ServiceUnavailable(w http.ResponseWriter, msg string) {
	response.ErrorResponse(w, http.StatusServiceUnavailable, msg)
}
