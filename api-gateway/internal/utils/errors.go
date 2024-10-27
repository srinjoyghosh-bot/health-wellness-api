package utils

import "fmt"

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.Message)
}

// Helper functions to create AppError

func NewBadRequestError(message string) *AppError {
	return &AppError{Code: 400, Message: message}
}

func NewUnauthorizedError(message string) *AppError {
	return &AppError{Code: 401, Message: message}
}

func NewNotFoundError(message string) *AppError {
	return &AppError{Code: 404, Message: message}
}

func NewInternalServerError(message string) *AppError {
	return &AppError{Code: 500, Message: message}
}
