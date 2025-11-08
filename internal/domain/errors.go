package domain

import (
	"errors"

	"github.com/crossmint/megaverse-challenge/internal/domain/entities"
)

var (
	// ErrInvalidPosition indicates that the position is invalid (negative coordinates)
	ErrInvalidPosition = errors.New("invalid position: coordinates must be non-negative")

	// ErrOutOfBounds indicates that the position is outside the megaverse boundaries
	ErrOutOfBounds = errors.New("position out of bounds")

	// ErrInvalidSoloonColor indicates an invalid color for a Soloon
	ErrInvalidSoloonColor = errors.New("invalid soloon color: must be blue, red, purple, or white")

	// ErrInvalidComethDirection indicates an invalid direction for a Cometh
	ErrInvalidComethDirection = errors.New("invalid cometh direction: must be up, down, left, or right")
)

func init() {
	entities.RegisterValidationErrors(ErrInvalidPosition, ErrInvalidSoloonColor, ErrInvalidComethDirection)
	entities.RegisterMegaverseErrors(ErrOutOfBounds)
}

// APIError represents a detailed API error with status code
type APIError struct {
	StatusCode int
	Message    string
	Endpoint   string
}

func (e *APIError) Error() string {
	return e.Message
}

// NewAPIError creates a new API error
func NewAPIError(statusCode int, message, endpoint string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Message:    message,
		Endpoint:   endpoint,
	}
}
