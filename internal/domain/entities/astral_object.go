package entities

// AstralObject is the base interface for all megaverse objects.
type AstralObject interface {
	GetPosition() Position
	GetType() string
	Validate() error
}
