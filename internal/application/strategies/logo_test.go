package strategies

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/crossmint/megaverse-challenge/internal/domain"
	"github.com/crossmint/megaverse-challenge/internal/domain/entities"
)

type stubRepository struct {
	goal *domain.GoalMap
}

func (s *stubRepository) CreatePolyanet(context.Context, entities.Position) error {
	panic("not implemented")
}
func (s *stubRepository) CreateSoloon(context.Context, entities.Position, entities.SoloonColor) error {
	panic("not implemented")
}
func (s *stubRepository) CreateCometh(context.Context, entities.Position, entities.ComethDirection) error {
	panic("not implemented")
}
func (s *stubRepository) DeleteObject(context.Context, string, entities.Position) error {
	panic("not implemented")
}
func (s *stubRepository) GetCurrentMap(context.Context) (*entities.Megaverse, error) {
	panic("not implemented")
}
func (s *stubRepository) IsHealthy(context.Context) error                     { panic("not implemented") }
func (s *stubRepository) GetGoalMap(context.Context) (*domain.GoalMap, error) { return s.goal, nil }

func TestLogoPatternGeneratePlan(t *testing.T) {
	repo := &stubRepository{goal: &domain.GoalMap{Goal: [][]string{
		{"POLYANET", "SPACE", "BLUE_SOLOON"},
		{"SPACE", "RED_SOLOON", "LEFT_COMETH"},
	}}}

	strategy := NewLogoPatternStrategy(repo)
	plan, err := strategy.GeneratePlan(context.Background())
	require.NoError(t, err)

	require.Equal(t, OrderParallel, plan.Order)
	require.Len(t, plan.Objects, 4)

	seen := make(map[entities.Position]string)

	for _, obj := range plan.Objects {
		seen[obj.GetPosition()] = obj.GetType()
		switch o := obj.(type) {
		case *entities.Polyanet:
			require.Equal(t, entities.Position{Row: 0, Column: 0}, o.Position)
		case *entities.Soloon:
			switch o.Position {
			case (entities.Position{Row: 0, Column: 2}):
				require.Equal(t, entities.BlueSoloon, o.Color)
			case (entities.Position{Row: 1, Column: 1}):
				require.Equal(t, entities.RedSoloon, o.Color)
			default:
				t.Fatalf("unexpected soloon position %+v", o.Position)
			}
		case *entities.Cometh:
			require.Equal(t, entities.Position{Row: 1, Column: 2}, o.Position)
			require.Equal(t, entities.LeftCometh, o.Direction)
		default:
			t.Fatalf("unexpected object type %T", obj)
		}
	}

	require.Equal(t, map[entities.Position]string{
		{Row: 0, Column: 0}: "POLYANET",
		{Row: 0, Column: 2}: "SOLOON",
		{Row: 1, Column: 1}: "SOLOON",
		{Row: 1, Column: 2}: "COMETH",
	}, seen)
}
