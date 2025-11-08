package domain

import (
	"context"

	"github.com/crossmint/megaverse-challenge/internal/domain/entities"
)

// MegaverseRepository defines the interface for interacting with the megaverse API
type MegaverseRepository interface {
	// CreatePolyanet creates a new Polyanet at the specified position
	CreatePolyanet(ctx context.Context, position entities.Position) error

	// CreateSoloon creates a new Soloon with the specified color at the given position
	CreateSoloon(ctx context.Context, position entities.Position, color entities.SoloonColor) error

	// CreateCometh creates a new Cometh with the specified direction at the given position
	CreateCometh(ctx context.Context, position entities.Position, direction entities.ComethDirection) error

	// DeleteObject removes an astral object at the specified position
	DeleteObject(ctx context.Context, objectType string, position entities.Position) error

	// GetGoalMap retrieves the goal map for the current challenge phase
	GetGoalMap(ctx context.Context) (*GoalMap, error)

	// GetCurrentMap retrieves the current state of the megaverse (if available)
	GetCurrentMap(ctx context.Context) (*entities.Megaverse, error)
}

// HealthChecker defines an interface for checking service health
type HealthChecker interface {
	// IsHealthy checks if the API service is healthy and reachable
	IsHealthy(ctx context.Context) error
}
