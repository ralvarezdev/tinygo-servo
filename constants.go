//go:build tinygo && (rp2040 || rp2350)

package tinygo_servo

const (
	// LeftLimitAngle is the angle that represents the left limit position of the servo motors.
	LeftLimitAngle uint16 = 0

	// RightLimitAngle is the angle that represents the right limit position of the servo motors.
	RightLimitAngle uint16 = 180
)