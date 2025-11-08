package entities

// ComethDirection represents the possible directions for a Cometh.
type ComethDirection string

const (
	UpCometh    ComethDirection = "up"
	DownCometh  ComethDirection = "down"
	LeftCometh  ComethDirection = "left"
	RightCometh ComethDirection = "right"
)

// Cometh represents a comet object with direction.
type Cometh struct {
	Position  Position
	Direction ComethDirection `json:"direction"`
}

func (c *Cometh) GetPosition() Position {
	return c.Position
}

func (c *Cometh) GetType() string {
	return "COMETH"
}

func (c *Cometh) Validate() error {
	if c.Position.Row < 0 || c.Position.Column < 0 {
		return invalidPositionError()
	}
	switch c.Direction {
	case UpCometh, DownCometh, LeftCometh, RightCometh:
		return nil
	default:
		return invalidComethDirectionError()
	}
}
