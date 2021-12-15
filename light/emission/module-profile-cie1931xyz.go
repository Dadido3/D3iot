// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package emission

import (
	"fmt"
)

// ModuleProfileCIE1931XYZ describes a module that is already in CIE 1931 XYZ color space.
// This only converts the color space into one with a relative luminance.
// No chromatic adaptation is done.
//
// This assumes that the device does the correct color transformation.
type ModuleProfileCIE1931XYZ struct {
	WhitePointColor CIE1931XYZAbs // The white point of the module.

	// Transfer function to convert from a linear device color space into a non linear device color space, and vice versa.
	// Set to nil if your DCS is linear.
	TransferFunc TransferFunction
}

// Check if it implements ModuleProfile.
var _ ModuleProfile = &ModuleProfileCIE1931XYZ{}

// Channels returns the dimensionality of the device color space.
func (e *ModuleProfileCIE1931XYZ) Channels() int {
	return 3
}

// WhitePoint returns the white point as a CIE 1931 XYZ color.
// This is also the brightest color a module can output.
func (e *ModuleProfileCIE1931XYZ) WhitePoint() CIE1931XYZAbs {
	return e.WhitePointColor
}

// ChannelPoints returns a list of channel colors.
// Depending on the module type, this could be the colors for:
//
//	- Single white emitter.
//	- RGB emitters.
//	- RGB + white emitters.
//	- RGB + cold white + warm white emitters.
func (e *ModuleProfileCIE1931XYZ) ChannelPoints() []CIE1931XYZAbs {
	return []CIE1931XYZAbs{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}}
}

// XYZToDCS takes a color and returns a vector in the device color space that reproduces the given color as closely as possible.
//
// Short: XYZ --> device color space.
func (e *ModuleProfileCIE1931XYZ) XYZToDCS(color CIE1931XYZAbs) DCSVector {
	v := LinDCSVector{color.X, color.Y, color.Z}

	// Scale so that the white point would result in Y = 1.0
	v = v.Scaled(1 / e.WhitePointColor.Y)

	return v.ClampedAndDeLinearized(e.TransferFunc)
}

// DCSToXYZ takes a vector from the device color space and returns the color it represents.
//
// Short: Device color space --> XYZ.
func (e *ModuleProfileCIE1931XYZ) DCSToXYZ(v DCSVector) (CIE1931XYZAbs, error) {
	if v.Channels() != e.Channels() {
		return CIE1931XYZAbs{}, fmt.Errorf("unexpected amount of channels. Got %d, want %d", v.Channels(), e.Channels())
	}

	linV := v.ClampedAndLinearized(e.TransferFunc)

	// Scale it up.
	linV = linV.Scaled(e.WhitePointColor.Y)

	return CIE1931XYZAbs{linV[0], linV[1], linV[2]}, nil
}

func (e *ModuleProfileCIE1931XYZ) TransferFunction() TransferFunction {
	return e.TransferFunc
}
