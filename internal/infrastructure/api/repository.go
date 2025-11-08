package api

import (
	"context"
	"errors"
	"fmt"

	"github.com/crossmint/megaverse-challenge/internal/domain"
	"github.com/crossmint/megaverse-challenge/internal/domain/entities"
)

// Repository implements the MegaverseRepository interface using the HTTP API
type Repository struct {
	client *Client
}

// NewRepository creates a new API repository
func NewRepository(client *Client) *Repository {
	return &Repository{
		client: client,
	}
}

// CreatePolyanet creates a new Polyanet at the specified position
func (r *Repository) CreatePolyanet(ctx context.Context, position entities.Position) error {
	req := CreatePolyanetRequest{
		Row:         position.Row,
		Column:      position.Column,
		CandidateID: r.client.GetCandidateID(),
	}

	return r.client.Post(ctx, "/polyanets", req)
}

// CreateSoloon creates a new Soloon with the specified color at the given position
func (r *Repository) CreateSoloon(ctx context.Context, position entities.Position, color entities.SoloonColor) error {
	req := CreateSoloonRequest{
		Row:         position.Row,
		Column:      position.Column,
		Color:       string(color),
		CandidateID: r.client.GetCandidateID(),
	}

	return r.client.Post(ctx, "/soloons", req)
}

// CreateCometh creates a new Cometh with the specified direction at the given position
func (r *Repository) CreateCometh(ctx context.Context, position entities.Position, direction entities.ComethDirection) error {
	req := CreateComethRequest{
		Row:         position.Row,
		Column:      position.Column,
		Direction:   string(direction),
		CandidateID: r.client.GetCandidateID(),
	}

	return r.client.Post(ctx, "/comeths", req)
}

// DeleteObject removes an astral object at the specified position
func (r *Repository) DeleteObject(ctx context.Context, objectType string, position entities.Position) error {
	req := DeleteRequest{
		Row:         position.Row,
		Column:      position.Column,
		CandidateID: r.client.GetCandidateID(),
	}

	var endpoint string
	switch objectType {
	case "POLYANET":
		endpoint = "/polyanets"
	case "SOLOON":
		endpoint = "/soloons"
	case "COMETH":
		endpoint = "/comeths"
	default:
		return fmt.Errorf("unknown object type: %s", objectType)
	}

	return r.client.Delete(ctx, endpoint, req)
}

// GetGoalMap retrieves the goal map for the current challenge phase
func (r *Repository) GetGoalMap(ctx context.Context) (*domain.GoalMap, error) {
	endpoint := fmt.Sprintf("/map/%s/goal", r.client.GetCandidateID())

	var goalMap domain.GoalMap
	if err := r.client.Get(ctx, endpoint, &goalMap); err != nil {
		return nil, fmt.Errorf("failed to get goal map: %w", err)
	}

	return &goalMap, nil
}

// GetCurrentMap retrieves the current state of the megaverse
func (r *Repository) GetCurrentMap(ctx context.Context) (*entities.Megaverse, error) {
	// callers should treat a 404 as “not supported”.
	endpoint := fmt.Sprintf("/map/%s", r.client.GetCandidateID())

	type apiCell struct {
		Type      *int   `json:"type,omitempty"`
		Color     string `json:"color,omitempty"`
		Direction string `json:"direction,omitempty"`
	}

	var response struct {
		Map struct {
			Content [][]*apiCell `json:"content"`
		} `json:"map"`
	}

	if err := r.client.Get(ctx, endpoint, &response); err != nil {
		// If the endpoint doesn't exist, return nil
		var apiErr *domain.APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == 404 {
			return nil, fmt.Errorf("current map endpoint not available")
		}
		return nil, fmt.Errorf("failed to get current map: %w", err)
	}

	content := response.Map.Content
	height := len(content)
	if height == 0 {
		return entities.NewMegaverse(0, 0), nil
	}
	width := len(content[0])

	megaverse := entities.NewMegaverse(width, height)

	for row, rowData := range content {
		for col, cell := range rowData {
			if cell == nil || cell.Type == nil {
				continue
			}

			pos := entities.Position{Row: row, Column: col}

			// The API uses numeric type codes: 0 = polyanet, 1 = soloon, 2 = cometh.
			switch *cell.Type {
			case 0:
				megaverse.PlaceObject(&entities.Polyanet{Position: pos})
			case 1:
				megaverse.PlaceObject(&entities.Soloon{
					Position: pos,
					Color:    entities.SoloonColor(cell.Color),
				})
			case 2:
				megaverse.PlaceObject(&entities.Cometh{
					Position:  pos,
					Direction: entities.ComethDirection(cell.Direction),
				})
			}
		}
	}

	return megaverse, nil
}

// IsHealthy checks if the API service is healthy and reachable
func (r *Repository) IsHealthy(ctx context.Context) error {
	// Try to fetch the goal map as a health check
	_, err := r.GetGoalMap(ctx)
	return err
}
