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

func (s *LogoPatternStrategy) GeneratePlan(ctx context.Context) (CreationPlan, error) {
	// Always fetch the latest goal map so we mirror Crossmint's target layout even if it changes between runs.
	goalMap, err := s.repository.GetGoalMap(ctx)
	if err != nil {
		return CreationPlan{}, fmt.Errorf("failed to fetch goal map: %w", err)
	}
	s.goalMap = goalMap

	if goalMap.Goal == nil || len(goalMap.Goal) == 0 {
		return CreationPlan{}, fmt.Errorf("goal map is empty")
	}

	var objects []entities.AstralObject

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
		// There can be 100+ objects; running in parallel keeps the overall runtime acceptable.
		Order: OrderParallel,
	}, nil
}

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
		// The API occasionally introduces new tokens; log them so we can extend support without failing the build.
		fmt.Printf("Warning: Unknown cell value '%s' at position (%d, %d)\n", cellValue, row, col)
		return nil
	}
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
