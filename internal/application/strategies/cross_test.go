package strategies

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/crossmint/megaverse-challenge/internal/domain/entities"
)

func TestCrossPatternGeneratePlan(t *testing.T) {
	strategy := NewCrossPatternStrategy()

	plan, err := strategy.GeneratePlan(context.Background())
	require.NoError(t, err)

	require.Equal(t, OrderParallel, plan.Order)
	require.Len(t, plan.Objects, 13)

	expected := map[entities.Position]struct{}{
		{Row: 2, Column: 2}: {},
		{Row: 2, Column: 8}: {},
		{Row: 3, Column: 3}: {},
		{Row: 3, Column: 7}: {},
		{Row: 4, Column: 4}: {},
		{Row: 4, Column: 6}: {},
		{Row: 5, Column: 5}: {},
		{Row: 6, Column: 4}: {},
		{Row: 6, Column: 6}: {},
		{Row: 7, Column: 3}: {},
		{Row: 7, Column: 7}: {},
		{Row: 8, Column: 2}: {},
		{Row: 8, Column: 8}: {},
	}

	for _, obj := range plan.Objects {
		require.Equal(t, "POLYANET", obj.GetType())
		pos := obj.GetPosition()
		if _, ok := expected[pos]; !ok {
			t.Fatalf("unexpected position %+v", pos)
		}
		delete(expected, pos)
	}

	require.Empty(t, expected, "not all expected positions were created")
}
