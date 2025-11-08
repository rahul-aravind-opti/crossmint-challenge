package strategies

import (
	"context"
	"fmt"

	"github.com/crossmint/megaverse-challenge/internal/domain/entities"
)

// CrossPatternStrategy implements the X-shaped cross pattern for Phase 1
type CrossPatternStrategy struct {
	gridSize int
	startRow int
	endRow   int
}

// NewCrossPatternStrategy creates a new cross pattern strategy
func NewCrossPatternStrategy() *CrossPatternStrategy {
	return &CrossPatternStrategy{
		gridSize: 11, // 11x11 grid for Phase 1
		startRow: 2,  // Official challenge starts the cross at row 2
		endRow:   8,  // ...and ends it at row 8
	}
}

// GetName returns the name of the strategy
func (s *CrossPatternStrategy) GetName() string {
	return "Cross Pattern (Phase 1)"
}

// GeneratePlan creates a plan for the X-shaped cross pattern
func (s *CrossPatternStrategy) GeneratePlan(_ context.Context) (CreationPlan, error) {
	var objects []entities.AstralObject

	// Sanity check configuration
	if s.startRow < 0 || s.endRow >= s.gridSize || s.startRow > s.endRow {
		return CreationPlan{}, fmt.Errorf("invalid cross configuration: start=%d end=%d size=%d",
			s.startRow, s.endRow, s.gridSize)
	}

	// Create Polyanets along both diagonals
	for i := s.startRow; i <= s.endRow; i++ {
		// Main diagonal (top-left to bottom-right)
		objects = append(objects, &entities.Polyanet{
			Position: entities.Position{Row: i, Column: i},
		})

		// Anti-diagonal (top-right to bottom-left)
		// Skip the (single) center point as it's already added by the main diagonal
		if i != s.gridSize/2 {
			objects = append(objects, &entities.Polyanet{
				Position: entities.Position{Row: i, Column: s.gridSize - 1 - i},
			})
		}
	}

	return CreationPlan{
		Objects: objects,
		Order:   OrderParallel,
	}, nil
}

// GetGridSize returns the dimensions of the megaverse for this pattern
func (s *CrossPatternStrategy) GetGridSize() (width, height int) {
	return s.gridSize, s.gridSize
}
