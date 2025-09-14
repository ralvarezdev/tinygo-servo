//go:build tinygo && (rp2040 || rp2350)

package tinygo_servo

import (
	tinygotypes "github.com/ralvarezdev/tinygo-types"
)

const (
	// ErrorCodeServoStartNumber is the starting number for servo-related error codes.
	ErrorCodeServoStartNumber uint16 = 5230
)

const (
	ErrorCodeServoFailedToConfigurePWM tinygotypes.ErrorCode = tinygotypes.ErrorCode(iota + ErrorCodeServoStartNumber)
	ErrorCodeServoFailedToInitializeServo
	ErrorCodeServoAngleOutOfRange
	ErrorCodeServoAngleBelowMinPulseWidth
	ErrorCodeServoAngleAboveMaxPulseWidth
	ErrorCodeServoFailedToSetServoAngle
	ErrorCodeServoNilHandler
)