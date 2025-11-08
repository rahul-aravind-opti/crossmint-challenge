package entities

import "errors"

var (
	errInvalidPosition        = errors.New("invalid position: coordinates must be non-negative")
	errInvalidSoloonColor     = errors.New("invalid soloon color: must be blue, red, purple, or white")
	errInvalidComethDirection = errors.New("invalid cometh direction: must be up, down, left, or right")
	errOutOfBounds            = errors.New("position out of bounds")
)

// RegisterValidationErrors allows the domain layer to inject canonical validation errors.
func RegisterValidationErrors(invalidPosition, invalidSoloonColor, invalidComethDirection error) {
	if invalidPosition != nil {
		errInvalidPosition = invalidPosition
	}
	if invalidSoloonColor != nil {
		errInvalidSoloonColor = invalidSoloonColor
	}
	if invalidComethDirection != nil {
		errInvalidComethDirection = invalidComethDirection
	}
}

// RegisterMegaverseErrors allows the domain layer to inject megaverse-related errors.
func RegisterMegaverseErrors(outOfBounds error) {
	if outOfBounds != nil {
		errOutOfBounds = outOfBounds
	}
}

func invalidPositionError() error {
	return errInvalidPosition
}

func invalidSoloonColorError() error {
	return errInvalidSoloonColor
}

func invalidComethDirectionError() error {
	return errInvalidComethDirection
}

func outOfBoundsError() error {
	return errOutOfBounds
}
