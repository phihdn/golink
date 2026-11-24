package model

type ErrorResponse struct {
	Message string `json:"message" example:"Internal Server Error"`
}
