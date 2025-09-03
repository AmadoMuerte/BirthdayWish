package response

import (
	"net/http"

	"github.com/go-chi/render"
)

type MessageResponse struct {
	Message string `json:"message,omitempty"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

type ErrorResponse struct {
	Message string `json:"message,omitempty"`
	Details string `json:"details,omitempty"`
}

func Error(msg string) MessageResponse {
	return MessageResponse{
		Message: msg,
	}
}

func ErrorResponseJSON(w http.ResponseWriter, r *http.Request, status int, message string) {
	render.Status(r, status)
	render.JSON(w, r, ErrorResponse{
		Message: message,
	})
}

func ErrorResponseWithDetails(w http.ResponseWriter, r *http.Request, status int, message, details string) {
	render.Status(r, status)
	render.JSON(w, r, ErrorResponse{
		Message: message,
		Details: details,
	})
}

func Success(msg string) MessageResponse {
	return MessageResponse{
		Message: msg,
	}
}

func SuccessResponse(w http.ResponseWriter, r *http.Request, status int, data any) {
	render.Status(r, status)
	render.JSON(w, r, data)
}
