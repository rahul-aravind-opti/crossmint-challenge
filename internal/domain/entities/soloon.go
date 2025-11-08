package entities

// SoloonColor represents the possible colors for a Soloon.
type SoloonColor string

const (
	BlueSoloon   SoloonColor = "blue"
	RedSoloon    SoloonColor = "red"
	PurpleSoloon SoloonColor = "purple"
	WhiteSoloon  SoloonColor = "white"
)

// Soloon represents a moon object with color.
type Soloon struct {
	Position Position
	Color    SoloonColor `json:"color"`
}

func (s *Soloon) GetPosition() Position {
	return s.Position
}

func (s *Soloon) GetType() string {
	return "SOLOON"
}

func (s *Soloon) Validate() error {
	if s.Position.Row < 0 || s.Position.Column < 0 {
		return invalidPositionError()
	}
	switch s.Color {
	case BlueSoloon, RedSoloon, PurpleSoloon, WhiteSoloon:
		return nil
	default:
		return invalidSoloonColorError()
	}
}
