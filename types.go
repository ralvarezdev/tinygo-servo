//go:build tinygo && (rp2040 || rp2350)

package tinygo_servo

import (
	"time"

	"machine"

	tinygotypes "github.com/ralvarezdev/tinygo-types"
	tinygodriversservo "tinygo.org/x/drivers/servo"
	tinygologger "github.com/ralvarezdev/tinygo-logger"
)

type (
	// DefaultHandler is the default implementation of the Servo interface
	DefaultHandler struct {
		afterSetAngleFunc  func(angle uint16)
		isMovementEnabled  func() bool
		isDirectionInverted bool
		frequency           uint16
		minPulseWidth       uint16
		halfPulseWidth      uint16
		maxPulseWidth       uint16
		rangePulseWidth     uint16
		centerAngle         uint16
		maxAngle            uint16
		servo               tinygodriversservo.Servo
		angle               uint16
		logger  		 tinygologger.Logger
	}
)

var (
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
// maxAngle: The maximum angle the servo can move from the center
// isDirectionInverted: Whether the direction of the servo motor is inverted
// logger: The logger instance for logging messages
//
// Returns:
//
// An instance of DefaultHandler and an error if any occurred during initialization
func NewDefaultHandler(
	pwm machine.PWM,
	pin machine.Pin,
	afterSetAngleFunc func(angle uint16),
	isMovementEnabled func() bool,
	frequency uint16,
	minPulseWidth uint16,
	maxPulseWidth uint16,
	centerAngle uint16,
	maxAngle uint16,
	isDirectionInverted bool,
	logger tinygologger.Logger,
) (*DefaultHandler, tinygotypes.ErrorCode) {
	// Configure the PWM
	if err := pwm.Configure(
		machine.PWMConfig{
			Period: uint64(time.Second / time.Duration(frequency)),
		},
	); err != nil {
		return nil, ErrorCodeServoFailedToConfigurePWM
	}

	// Create a new instance of the servo
	servo, err := tinygodriversservo.New(pwm, pin)
	if err != nil {
		return nil, ErrorCodeServoFailedToInitializeServo
	}

	// Calculate the half pulse and range pulse
	halfPulseWidth := (maxPulseWidth + minPulseWidth) / 2
	rangePulseWidth := maxPulseWidth - minPulseWidth

	// Initialize the servo with the provided parameters
	handler := &DefaultHandler{
		afterSetAngleFunc:  afterSetAngleFunc,
		isMovementEnabled:  isMovementEnabled,
		isDirectionInverted: isDirectionInverted,
		frequency:           frequency,
		minPulseWidth:       minPulseWidth,
		halfPulseWidth:      halfPulseWidth,
		maxPulseWidth:       maxPulseWidth,
		rangePulseWidth:     rangePulseWidth,
		servo:               servo,
		angle:               centerAngle,
		centerAngle:         centerAngle,
		logger:  		 logger,
	}

	// Center the servo on initialization
	_ = handler.SetAngleToCenter()
	return handler, tinygotypes.ErrorCodeNil
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
func (h *DefaultHandler) SetAngle(angle uint16) tinygotypes.ErrorCode {
	// Check if the angle is within the valid range
	if angle < h.centerAngle-h.maxAngle || angle > h.centerAngle+h.maxAngle {
		return ErrorCodeServoAngleOutOfRange
	}
	if angle < LeftLimitAngle || angle > RightLimitAngle {
		return ErrorCodeServoAngleOutOfRange
	}

	// Check if the angle is the same as the current angle
	if angle == h.angle {
		return tinygotypes.ErrorCodeNil
	}

	// Check if the direction is inverted
	if h.isDirectionInverted {
		angle = RightLimitAngle - (angle - LeftLimitAngle)
	}

	// Update the current angle
	h.angle = angle

	// Set the servo angle
	if h.isMovementEnabled == nil || h.isMovementEnabled() {
		if err := h.servo.SetAngleWithMicroseconds(
			int(angle),
			int(h.minPulseWidth),
			int(h.maxPulseWidth),
		); err != nil {
			return ErrorCodeServoFailedToSetServoAngle
		}
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

	return tinygotypes.ErrorCodeNil
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
func (h *DefaultHandler) SetAngleToCenter() tinygotypes.ErrorCode {
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
func (h *DefaultHandler) SetAngleRelativeToCenter(relativeAngle int16) tinygotypes.ErrorCode {
	// Calculate the absolute angle based on the center angle and relative angle
	absoluteAngle := int16(h.centerAngle) + relativeAngle

	// Check if the absolute angle is within the left and right limits
	if absoluteAngle < int16(LeftLimitAngle) || absoluteAngle > int16(RightLimitAngle) {
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
func (h *DefaultHandler) SetAngleToRight(angle uint16) tinygotypes.ErrorCode {
	return h.SetAngleRelativeToCenter(-int16(angle))
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
func (h *DefaultHandler) SetAngleToLeft(angle uint16) tinygotypes.ErrorCode {
	return h.SetAngleRelativeToCenter(int16(angle))
}

// SetDirectionToCenter sets the direction to center
func (h *DefaultHandler) SetDirectionToCenter() tinygotypes.ErrorCode {
	return h.SetAngleToCenter()
}

// SetDirectionToRight sets the direction to right
//
// Parameters:
//
// angle: The angle value to move the servo to the left, must be between 0 and the left limit
//
// Returns:
//
// An error if the angle is not within the left limit
func (h *DefaultHandler) SetDirectionToRight(angle uint16) tinygotypes.ErrorCode {
	return h.SetAngleToLeft(angle)
}

// SetDirectionToLeft sets the direction to left
//
// Parameters:
//
// angle: The angle value to move the servo to the right, must be between 0 and the right limit
//
// Returns:
//
// An error if the angle is not within the right limit
func (h *DefaultHandler) SetDirectionToLeft(angle uint16) tinygotypes.ErrorCode {
	return h.SetAngleToRight(angle)
}