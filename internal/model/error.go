package model

type ErrorResponse struct {
	Message string `json:"message" example:"Internal Server Error"`
}

// HTTPError is the body handlers send on a 4xx or 5xx response, e.g. {"error": "Invalid Input"}.
type HTTPError struct {
	Error string `json:"error" example:"Invalid Input"`
}

// The types below exist only for the Swagger docs. They have the same shape as
// HTTPError, but each one carries the example message for its status code, so
// the docs show "Not Found" for 404 instead of one shared example.

type BadRequestError struct {
	Error string `json:"error" example:"Invalid Input"`
}

type NotFoundError struct {
	Error string `json:"error" example:"Not Found"`
}

type InternalServerError struct {
	Error string `json:"error" example:"Internal Server Error"`
}
