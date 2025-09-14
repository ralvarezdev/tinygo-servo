//go:build tinygo && (rp2040 || rp2350)

package tinygo_servo

import (
	tinygotypes "github.com/ralvarezdev/tinygo-types"
)

type (
	// Handler is the interface to handle servo operations
	Handler interface {
		SetAngle(angle uint16) tinygotypes.ErrorCode
		GetAngle() uint16
		SetAngleRelativeToCenter(relativeAngle int16) tinygotypes.ErrorCode
		IsAngleCentered() bool
		SetAngleToCenter() tinygotypes.ErrorCode
		SetAngleToRight(angle uint16) tinygotypes.ErrorCode
		SetAngleToLeft(angle uint16) tinygotypes.ErrorCode
		SetDirectionToCenter() tinygotypes.ErrorCode
		SetDirectionToRight(angle uint16) tinygotypes.ErrorCode
		SetDirectionToLeft(angle uint16) tinygotypes.ErrorCode
	}
)
