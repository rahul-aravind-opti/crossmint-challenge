package strategies

import (
	"context"

	"github.com/crossmint/megaverse-challenge/internal/domain/entities"
)

// PatternStrategy defines the interface for different megaverse creation patterns
type PatternStrategy interface {
	// GetName returns the name of the strategy
	GetName() string

	// GeneratePlan creates a plan of objects to be placed in the megaverse
	GeneratePlan(ctx context.Context) (CreationPlan, error)

	// Validate checks if the current megaverse matches the expected pattern
	Validate(megaverse *entities.Megaverse) error

	// GetGridSize returns the dimensions of the megaverse for this pattern
	GetGridSize() (width, height int)
}

// ExecutionOrder defines the order in which objects should be created
type ExecutionOrder int

const (
	// OrderSequential creates objects one by one in order
	OrderSequential ExecutionOrder = iota

	// OrderParallel allows objects to be created in parallel (with rate limiting)
	OrderParallel

	// OrderBatched creates objects in batches
	OrderBatched
)

// CreationPlan represents a plan for creating objects in the megaverse
type CreationPlan struct {
	Objects   []entities.AstralObject
	Order     ExecutionOrder
	BatchSize int // Used only for OrderBatched
}
