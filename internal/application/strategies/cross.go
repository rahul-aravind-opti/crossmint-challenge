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

// Validate checks if the megaverse has the correct cross pattern
func (s *CrossPatternStrategy) Validate(megaverse *entities.Megaverse) error {
	if megaverse == nil {
		return fmt.Errorf("megaverse is nil")
	}

	if megaverse.Width != s.gridSize || megaverse.Height != s.gridSize {
		return fmt.Errorf("invalid grid size: expected %dx%d, got %dx%d",
			s.gridSize, s.gridSize, megaverse.Width, megaverse.Height)
	}

	// Check each position in the grid
	for row := 0; row < s.gridSize; row++ {
		for col := 0; col < s.gridSize; col++ {
			obj, _ := megaverse.GetObject(row, col)

			// Check if this position should have a Polyanet
			shouldHavePolyanet := row >= s.startRow && row <= s.endRow &&
				(col == row || col == s.gridSize-1-row)

			if shouldHavePolyanet {
				if obj == nil {
					return fmt.Errorf("missing Polyanet at position (%d, %d)", row, col)
				}
				if _, ok := obj.(*entities.Polyanet); !ok {
					return fmt.Errorf("expected Polyanet at position (%d, %d), got %s",
						row, col, obj.GetType())
				}
			} else {
				if obj != nil {
					return fmt.Errorf("unexpected object at position (%d, %d): %s",
						row, col, obj.GetType())
				}
			}
		}
	}

	return nil
}

// GetGridSize returns the dimensions of the megaverse for this pattern
func (s *CrossPatternStrategy) GetGridSize() (width, height int) {
	return s.gridSize, s.gridSize
}
