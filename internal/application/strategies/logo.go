package strategies

import (
	"context"
	"fmt"
	"strings"

	"github.com/crossmint/megaverse-challenge/internal/domain"
	"github.com/crossmint/megaverse-challenge/internal/domain/entities"
)

// LogoPatternStrategy implements the pattern based on the goal map for Phase 2
type LogoPatternStrategy struct {
	repository domain.MegaverseRepository
	goalMap    *domain.GoalMap
}

// NewLogoPatternStrategy creates a new logo pattern strategy
func NewLogoPatternStrategy(repository domain.MegaverseRepository) *LogoPatternStrategy {
	return &LogoPatternStrategy{
		repository: repository,
	}
}

// GetName returns the name of the strategy
func (s *LogoPatternStrategy) GetName() string {
	return "Logo Pattern (Phase 2)"
}

// GeneratePlan creates a plan based on the goal map
func (s *LogoPatternStrategy) GeneratePlan(ctx context.Context) (CreationPlan, error) {
	// Fetch the goal map from the API
	goalMap, err := s.repository.GetGoalMap(ctx)
	if err != nil {
		return CreationPlan{}, fmt.Errorf("failed to fetch goal map: %w", err)
	}
	s.goalMap = goalMap

	if goalMap.Goal == nil || len(goalMap.Goal) == 0 {
		return CreationPlan{}, fmt.Errorf("goal map is empty")
	}

	var objects []entities.AstralObject

	// Parse the goal map and create objects
	for row, rowData := range goalMap.Goal {
		for col, cellValue := range rowData {
			obj := s.parseGoalCell(cellValue, row, col)
			if obj != nil {
				objects = append(objects, obj)
			}
		}
	}

	return CreationPlan{
		Objects: objects,
		Order:   OrderParallel,
	}, nil
}

// parseGoalCell converts a goal map cell value to an AstralObject
func (s *LogoPatternStrategy) parseGoalCell(cellValue string, row, col int) entities.AstralObject {
	cellValue = strings.ToUpper(strings.TrimSpace(cellValue))
	position := entities.Position{Row: row, Column: col}

	switch cellValue {
	case "SPACE", "":
		return nil

	case "POLYANET":
		return &entities.Polyanet{Position: position}

	case "BLUE_SOLOON":
		return &entities.Soloon{
			Position: position,
			Color:    entities.BlueSoloon,
		}

	case "RED_SOLOON":
		return &entities.Soloon{
			Position: position,
			Color:    entities.RedSoloon,
		}

	case "PURPLE_SOLOON":
		return &entities.Soloon{
			Position: position,
			Color:    entities.PurpleSoloon,
		}

	case "WHITE_SOLOON":
		return &entities.Soloon{
			Position: position,
			Color:    entities.WhiteSoloon,
		}

	case "UP_COMETH":
		return &entities.Cometh{
			Position:  position,
			Direction: entities.UpCometh,
		}

	case "DOWN_COMETH":
		return &entities.Cometh{
			Position:  position,
			Direction: entities.DownCometh,
		}

	case "LEFT_COMETH":
		return &entities.Cometh{
			Position:  position,
			Direction: entities.LeftCometh,
		}

	case "RIGHT_COMETH":
		return &entities.Cometh{
			Position:  position,
			Direction: entities.RightCometh,
		}

	default:
		// Unknown cell value, log and skip
		fmt.Printf("Warning: Unknown cell value '%s' at position (%d, %d)\n", cellValue, row, col)
		return nil
	}
}

// Validate checks if the megaverse matches the goal map
func (s *LogoPatternStrategy) Validate(megaverse *entities.Megaverse) error {
	if megaverse == nil {
		return fmt.Errorf("megaverse is nil")
	}

	if s.goalMap == nil {
		return fmt.Errorf("goal map not loaded")
	}

	goalHeight := len(s.goalMap.Goal)
	goalWidth := 0
	if goalHeight > 0 {
		goalWidth = len(s.goalMap.Goal[0])
	}

	if megaverse.Width != goalWidth || megaverse.Height != goalHeight {
		return fmt.Errorf("megaverse size mismatch: expected %dx%d, got %dx%d",
			goalWidth, goalHeight, megaverse.Width, megaverse.Height)
	}

	// Check each position against the goal
	for row, rowData := range s.goalMap.Goal {
		for col, expectedValue := range rowData {
			obj, err := megaverse.GetObject(row, col)
			if err != nil {
				return fmt.Errorf("failed to retrieve object at (%d, %d): %w", row, col, err)
			}

			expectedObj := s.parseGoalCell(expectedValue, row, col)

			if expectedObj == nil && obj == nil {
				continue // Both empty, that's correct
			}

			if expectedObj == nil && obj != nil {
				return fmt.Errorf("unexpected object at (%d, %d): expected empty, got %s",
					row, col, obj.GetType())
			}

			if expectedObj != nil && obj == nil {
				return fmt.Errorf("missing object at (%d, %d): expected %s",
					row, col, expectedObj.GetType())
			}

			// Check if types match
			if expectedObj.GetType() != obj.GetType() {
				return fmt.Errorf("type mismatch at (%d, %d): expected %s, got %s",
					row, col, expectedObj.GetType(), obj.GetType())
			}

			// Check specific attributes
			switch expected := expectedObj.(type) {
			case *entities.Soloon:
				if actual, ok := obj.(*entities.Soloon); ok {
					if expected.Color != actual.Color {
						return fmt.Errorf("soloon color mismatch at (%d, %d): expected %s, got %s",
							row, col, expected.Color, actual.Color)
					}
				}
			case *entities.Cometh:
				if actual, ok := obj.(*entities.Cometh); ok {
					if expected.Direction != actual.Direction {
						return fmt.Errorf("cometh direction mismatch at (%d, %d): expected %s, got %s",
							row, col, expected.Direction, actual.Direction)
					}
				}
			}
		}
	}

	return nil
}

// GetGridSize returns the dimensions based on the goal map
func (s *LogoPatternStrategy) GetGridSize() (width, height int) {
	if s.goalMap == nil || len(s.goalMap.Goal) == 0 {
		return 0, 0
	}

	height = len(s.goalMap.Goal)
	if height > 0 {
		width = len(s.goalMap.Goal[0])
	}

	return width, height
}

// GetExecutionOrder returns the recommended execution order for this pattern
func (s *LogoPatternStrategy) GetExecutionOrder() ExecutionOrder {
	// Logo pattern might have many objects, parallel execution is recommended
	// but we could also use batched if we group by object type
	return OrderParallel
}
