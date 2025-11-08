package entities

// Megaverse represents a 2D grid of the space.
type Megaverse struct {
	Grid   [][]AstralObject
	Width  int
	Height int
}

// NewMegaverse creates a new megaverse with the specified dimensions.
func NewMegaverse(width, height int) *Megaverse {
	grid := make([][]AstralObject, height)
	for i := range grid {
		grid[i] = make([]AstralObject, width)
	}
	return &Megaverse{
		Grid:   grid,
		Width:  width,
		Height: height,
	}
}

// PlaceObject places an astral object at the specified position.
func (m *Megaverse) PlaceObject(obj AstralObject) error {
	pos := obj.GetPosition()
	if pos.Row < 0 || pos.Row >= m.Height || pos.Column < 0 || pos.Column >= m.Width {
		return outOfBoundsError()
	}
	m.Grid[pos.Row][pos.Column] = obj
	return nil
}

// GetObject returns the object at the specified position.
func (m *Megaverse) GetObject(row, column int) (AstralObject, error) {
	if row < 0 || row >= m.Height || column < 0 || column >= m.Width {
		return nil, outOfBoundsError()
	}
	return m.Grid[row][column], nil
}

// Clear removes all objects from the megaverse.
func (m *Megaverse) Clear() {
	for i := range m.Grid {
		for j := range m.Grid[i] {
			m.Grid[i][j] = nil
		}
	}
}
