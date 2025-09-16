package tinygo_servo

import (
	tinygoerrors "github.com/ralvarezdev/tinygo-errors"
)

type (
	// Handler is the interface to handle servo operations
	Handler interface {
		SetAngle(angle uint16) tinygoerrors.ErrorCode
		GetAngle() uint16
		SetAngleRelativeToCenter(relativeAngle int16) tinygoerrors.ErrorCode
		IsAngleCentered() bool
		SetAngleToCenter() tinygoerrors.ErrorCode
		SetAngleToRight(angle uint16) tinygoerrors.ErrorCode
		SafeSetAngleToRight(angle uint16) tinygoerrors.ErrorCode
		SetAngleToLeft(angle uint16) tinygoerrors.ErrorCode
		SafeSetAngleToLeft(angle uint16) tinygoerrors.ErrorCode
	}
)
