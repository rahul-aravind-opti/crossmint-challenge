package domain

// GoalMap represents the target state of the megaverse.
type GoalMap struct {
	Goal [][]string `json:"goal"`
}

// ParseObjectType converts a goal map string to an object type.
func ParseObjectType(goalString string) (objectType string, attributes map[string]string) {
	attributes = make(map[string]string)

	switch goalString {
	case "SPACE":
		return "", nil
	case "POLYANET":
		return "POLYANET", nil
	case "BLUE_SOLOON":
		return "SOLOON", map[string]string{"color": "blue"}
	case "RED_SOLOON":
		return "SOLOON", map[string]string{"color": "red"}
	case "PURPLE_SOLOON":
		return "SOLOON", map[string]string{"color": "purple"}
	case "WHITE_SOLOON":
		return "SOLOON", map[string]string{"color": "white"}
	case "UP_COMETH":
		return "COMETH", map[string]string{"direction": "up"}
	case "DOWN_COMETH":
		return "COMETH", map[string]string{"direction": "down"}
	case "LEFT_COMETH":
		return "COMETH", map[string]string{"direction": "left"}
	case "RIGHT_COMETH":
		return "COMETH", map[string]string{"direction": "right"}
	default:
		return "", nil
	}
}
