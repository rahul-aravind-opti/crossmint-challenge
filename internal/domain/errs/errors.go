package errs

import "errors"

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

