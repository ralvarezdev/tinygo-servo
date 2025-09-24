package tinygo_servo

import (
	"machine"

	tinygoerrors "github.com/ralvarezdev/tinygo-errors"
	tinygologger "github.com/ralvarezdev/tinygo-logger"
	tinygopwm "github.com/ralvarezdev/tinygo-pwm"
)

type (
	// DefaultHandler is the default implementation of the Servo interface
	DefaultHandler struct {
		afterSetAngleFunc   func(angle uint16)
		isMovementEnabled   func() bool
		isDirectionInverted bool
		frequency           uint16
		minPulseWidth       uint32
		maxPulseWidth       uint32
		centerAngle         uint16
		actuationRange    uint16
		leftLimitAngle    uint16
		rightLimitAngle   uint16
		angle               uint16
		logger              tinygologger.Logger
		pwm 			  tinygopwm.PWM
		channel 		  uint8
		period 				  uint32
	}
)

var (
	// setPeriodPrefix is the prefix for the log message when setting the PWM period
	setPeriodPrefix = []byte("Set Servo PWM period to:")

	// setAnglePrefix is the prefix message for new angle setting
	setAnglePrefix = []byte("Set servo angle degrees to:")
)

// NewDefaultHandler creates a new instance of DefaultHandler
//
// Parameters:
//
// pwm: The PWM interface to control the servo
// pin: The pin connected to the servo
// afterSetAngleFunc: A callback function to be called after setting the angle
// isMovementEnabled: A function to check if movement is enabled
// frequency: The frequency of the PWM signal
// minPulseWidth: The minimum pulse width for the servo motor
// maxPulseWidth: The maximum pulse width for the servo motor
// centerAngle: The center angle of the servo motor
// maxLeftAngle: The maximum left angle from the center
// maxRightAngle: The maximum right angle from the center
// isDirectionInverted: Whether the direction of the servo motor is inverted
// logger: The logger instance for logging messages
//
// Returns:
//
// An instance of DefaultHandler and an error if any occurred during initialization
func NewDefaultHandler(
	pwm tinygopwm.PWM,
	pin machine.Pin,
	afterSetAngleFunc func(angle uint16),
	isMovementEnabled func() bool,
	frequency uint16,
	minPulseWidth uint32,
	maxPulseWidth uint32,
	actuationRange uint16,
	centerAngle uint16,
	maxLeftAngle uint16,
	maxRightAngle uint16,
	isDirectionInverted bool,
	logger tinygologger.Logger,
) (*DefaultHandler, tinygoerrors.ErrorCode) {
	// Check if the frequency is zero
	if frequency == 0 {
		return nil, ErrorCodeServoZeroFrequency
	}

	// Configure the PWM
	period := 1e9 / float64(frequency)
	if err := pwm.Configure(
		machine.PWMConfig{
			Period: uint64(period),
		},
	); err != nil {
		return nil, ErrorCodeServoFailedToConfigurePWM
	}

	// Log the configured period
	if logger != nil {
		logger.AddMessageWithUint32(
			setPeriodPrefix,
			uint32(period),
			true,
			true,
			false,
		)
		logger.Debug()
	}

	// Get the channel from the pin
	channel, err := pwm.Channel(pin)
	if err != nil {
		return nil, ErrorCodeServoFailedToGetPWMChannel
	}

	// Check if the min pulse width is valid
	if minPulseWidth == 0 || minPulseWidth >= uint32(period) {
		return nil, ErrorCodeServoInvalidMinPulseWidth
	}

	// Check if the max pulse width is valid
	if maxPulseWidth == 0 || maxPulseWidth >= uint32(period) {
		return nil, ErrorCodeServoInvalidMaxPulseWidth
	}

	// Check if the actuation range is valid
	if actuationRange == 0 || actuationRange > 360 {
		return nil, ErrorCodeServoInvalidActuationRange
	}

	// Check if the center angle is valid
	if centerAngle < 0 || centerAngle > actuationRange {
		return nil, ErrorCodeServoInvalidCenterAngle
	}

	// Calculate the left and right limit angles
	leftLimitAngle := centerAngle - maxLeftAngle
	rightLimitAngle := centerAngle + maxRightAngle

	// Check if the left limit angle is valid
	if leftLimitAngle < 0 {
		leftLimitAngle = 0
	}

	// Check if the right limit angle is valid
	if rightLimitAngle > actuationRange {
		rightLimitAngle = actuationRange
	}

	// Initialize the servo with the provided parameters
	handler := &DefaultHandler{
		afterSetAngleFunc:   afterSetAngleFunc,
		isMovementEnabled:   isMovementEnabled,
		isDirectionInverted: isDirectionInverted,
		frequency:           frequency,
		minPulseWidth:       minPulseWidth,
		maxPulseWidth:       maxPulseWidth,
		angle:               centerAngle,
		centerAngle:         centerAngle,
		actuationRange:    actuationRange,
		logger:              logger,
		pwm: 			  pwm,
		channel: 		  channel,
		leftLimitAngle:    leftLimitAngle,
		rightLimitAngle:   rightLimitAngle,
		period: 				  uint32(period),

	}

	// Center the servo on initialization
	_ = handler.SetAngleToCenter()
	return handler, tinygoerrors.ErrorCodeNil
}

// GetAngle returns the current angle of the servo motor
//
// Returns:
//
// The current angle of the servo motor
func (h *DefaultHandler) GetAngle() uint16 {
	return h.angle
}

// SetAngle sets the angle of the servo motor
//
// Parameters:
//
// angle: The angle to set the servo motor to, must be between 0 and the actuation range
func (h *DefaultHandler) SetAngle(angle uint16) tinygoerrors.ErrorCode {
	// Check if the direction is inverted
	if h.isDirectionInverted {
		angle = h.actuationRange - angle
	}

	// Check if the angle is within the valid range
	if angle < h.centerAngle-h.leftLimitAngle || angle > h.centerAngle+h.rightLimitAngle {
		return ErrorCodeServoAngleOutOfRange
	}

	// Check if the angle is the same as the current angle
	if angle == h.angle {
		return tinygoerrors.ErrorCodeNil
	}

	// Update the current angle
	h.angle = angle

	// Calculate the pulse
	pulse := uint32(h.minPulseWidth) + uint32(float64(h.maxPulseWidth-h.minPulseWidth) * float64(angle) / float64(h.actuationRange))


	// Set the servo angle
	if h.isMovementEnabled == nil || h.isMovementEnabled() {
		tinygopwm.SetDuty(
			h.pwm,
			h.channel,
			pulse,
			h.period,
		)
	}

	// Log the new angle if logger is provided
	if h.logger != nil {
		h.logger.AddMessageWithUint16(setAnglePrefix, angle, true, true, false)
		h.logger.Debug()
	}

	// Call the after set angle function if provided
	if h.afterSetAngleFunc != nil {
		h.afterSetAngleFunc(angle)
	}

	return tinygoerrors.ErrorCodeNil
}

// IsAngleCentered checks if the servo motor angle is centered
//
// Returns:
//
// True if the servo motor is centered, false otherwise
func (h *DefaultHandler) IsAngleCentered() bool {
	return h.angle == h.centerAngle
}

// SetAngleToCenter centers the servo motor to the middle position
//
// Returns:
//
// An error if the servo motor could not be centered
func (h *DefaultHandler) SetAngleToCenter() tinygoerrors.ErrorCode {
	return h.SetAngle(h.centerAngle)
}

// SetAngleRelativeToCenter sets the angle of the servo motor relative to the center position
//
// Parameters:
//
// relativeAngle: The relative angle value between -90 and 90 degrees
//
// Returns:
//
// An error if the relative angle is not within the left and right limits
func (h *DefaultHandler) SetAngleRelativeToCenter(relativeAngle int16) tinygoerrors.ErrorCode {
	// Calculate the absolute angle based on the center angle and relative angle
	absoluteAngle := int16(h.centerAngle) + relativeAngle

	// Check if the absolute angle is within the left and right limits
	if absoluteAngle < int16(h.leftLimitAngle) || absoluteAngle > int16(h.rightLimitAngle) {
		return ErrorCodeServoAngleOutOfRange
	}

	// Set the servo angle
	return h.SetAngle(uint16(absoluteAngle))
}

// SetAngleToRight sets the servo motor to the right by a specified angle
//
// Parameters:
//
// angle: The angle value to move the servo to the right, must be between 0 and the right limit
//
// Returns:
//
// An error if the angle is not within the right limit
func (h *DefaultHandler) SetAngleToRight(angle uint16) tinygoerrors.ErrorCode {
	// Check if the angle is negative
	if angle < 0 {
		angle = 0
	}

	// Check if the angle is within the right limit
	if angle > h.rightLimitAngle-h.centerAngle {
		angle = h.rightLimitAngle - h.centerAngle
	}
	return h.SetAngleRelativeToCenter(int16(angle))
}

// SetAngleToLeft sets the servo motor to the left by a specified angle
//
// Parameters:
//
// angle: The angle value to move the servo to the left, must be between 0 and the left limit
//
// Returns:
//
// An error if the angle is not within the left limit
func (h *DefaultHandler) SetAngleToLeft(angle uint16) tinygoerrors.ErrorCode {
	// Check if the angle is negative
	if angle < 0 {
		angle = 0
	}

	// Check if the angle is within the left limit
	if angle > h.centerAngle-h.leftLimitAngle {
		angle = h.centerAngle - h.leftLimitAngle
	}
	return h.SetAngleRelativeToCenter(-int16(angle))
}