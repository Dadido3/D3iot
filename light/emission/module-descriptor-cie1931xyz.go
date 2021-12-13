// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package emission

import (
	"fmt"
)

// ModuleDescriptorCIE1931XYZ describes a module that is already in CIE 1931 XYZ color space.
// This only converts the color space into one with a relative luminance.
//
// This assumes that the device does the correct color transformation.
type ModuleDescriptorCIE1931XYZ struct {
	WhitePointColor CIE1931XYZColor // The white point of the module.

	// Transfer function to convert from a linear device color space into a non linear device color space, and vice versa.
	// Set to nil if your DCS is linear.
	TransferFunction TransferFunction
}

// Check if it implements ModuleDescriptor.
var _ ModuleDescriptor = &ModuleDescriptorCIE1931XYZ{}

// Channels returns the dimensionality of the device color space.
func (e *ModuleDescriptorCIE1931XYZ) Channels() int {
	return 3
}

// WhitePoint returns the white point as a CIE 1931 XYZ color.
// This is also the brightest color a module can output.
func (e *ModuleDescriptorCIE1931XYZ) WhitePoint() CIE1931XYZColor {
	return e.WhitePointColor
}

// ChannelPoints returns a list of channel colors.
// Depending on the module type, this could be the colors for:
//
//	- Single white emitter.
//	- RGB emitters.
//	- RGB + white emitters.
//	- RGB + cold white + warm white emitters.
func (e *ModuleDescriptorCIE1931XYZ) ChannelPoints() []CIE1931XYZColor {
	return []CIE1931XYZColor{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}}
}

// XYZToDCS takes a color and returns a vector in the device color space that reproduces the given color as closely as possible.
//
// Short: XYZ --> device color space.
func (e *ModuleDescriptorCIE1931XYZ) XYZToDCS(color CIE1931XYZColor) (DCSColor, error) {
	v := LinDCSColor{color.X, color.Y, color.Z}

	// Scale so that the white point would result in Y = 1.0
	v = v.Scaled(1 / e.WhitePointColor.Y)

	return v.ClampedAndDeLinearized(e.TransferFunction), nil
}

// DCSToXYZ takes a vector from the device color space and returns the color it represents.
//
// Short: Device color space --> XYZ.
func (e *ModuleDescriptorCIE1931XYZ) DCSToXYZ(v DCSColor) (CIE1931XYZColor, error) {
	if v.Channels() != e.Channels() {
		return CIE1931XYZColor{}, fmt.Errorf("unexpected amount of channels. Got %d, want %d", v.Channels(), e.Channels())
	}

	linV := v.ClampedAndLinearized(e.TransferFunction)

	// Scale it up.
	linV = linV.Scaled(e.WhitePointColor.Y)

	return CIE1931XYZColor{linV[0], linV[1], linV[2]}, nil
}
