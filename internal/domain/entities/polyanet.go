package entities

// Polyanet represents a basic planet object.
type Polyanet struct {
	Position Position
}

func (p *Polyanet) GetPosition() Position {
	return p.Position
}

func (p *Polyanet) GetType() string {
	return "POLYANET"
}

func (p *Polyanet) Validate() error {
	if p.Position.Row < 0 || p.Position.Column < 0 {
		return invalidPositionError()
	}
	return nil
}
