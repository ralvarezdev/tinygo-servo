package tinygo_servo

type (
	// Direction is an enum to represent the different servo directions for the vehicle.
	Direction uint8
)

const (
	DirectionNil Direction = iota
	DirectionLeft
	DirectionRight
	DirectionStraight
)

// InvertedDirection returns the inverted direction.
func (d Direction) InvertedDirection() Direction {
	switch d {
	case DirectionLeft:
		return DirectionRight
	case DirectionRight:
		return DirectionLeft
	case DirectionStraight:
		return DirectionStraight
	default:
		return DirectionNil
	}
}
