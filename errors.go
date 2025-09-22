package tinygo_servo

import (
	tinygoerrors "github.com/ralvarezdev/tinygo-errors"
)

const (
	// ErrorCodeServoStartNumber is the starting number for servo-related error codes.
	ErrorCodeServoStartNumber uint16 = 5230
)

const (
	ErrorCodeServoFailedToConfigurePWM tinygoerrors.ErrorCode = tinygoerrors.ErrorCode(iota + ErrorCodeServoStartNumber)
	ErrorCodeServoZeroFrequency
	ErrorCodeServoAngleOutOfRange
	ErrorCodeServoInvalidMinPulseWidth
	ErrorCodeServoInvalidMaxPulseWidth
	ErrorCodeServoNilHandler
	ErrorCodeServoUnknownDirection
	ErrorCodeServoFailedToGetPWMChannel
	ErrorCodeServoInvalidActuationRange
	ErrorCodeServoInvalidCenterAngle
)
