// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package emission

import (
	"fmt"
)

// ModuleDescriptor describes a group of LEDs that together generate a single color impression.
//
// This contains everything that is necessary to convert from the CIE 1931 XYZ color space into the device color space, and vice versa.
//
// Most devices contain just one module with some number of channels (E.g. RGBW light bulbs).
// But there can be multiple modules per device, e.g. multi headed lamps, addressable LED strips.
type ModuleDescriptor interface {
	// DCSChannels returns the dimensionality of the device color space.
	DCSChannels() int

	// XYZToDCS takes a color and returns a vector in the device color space that reproduces the given color as closely as possible.
	//
	// Short: XYZ --> device color space.
	XYZToDCS(color CIE1931XYZColor) (DCSColor, error)

	// DCSToXYZ takes a vector from the device color space and returns the color it represents.
	//
	// Short: Device color space --> XYZ.
	DCSToXYZ(v DCSColor) (CIE1931XYZColor, error)
}

// ModuleDescriptorGeneral is a general implementation of a ModuleDescriptor.
// It supports:
//
//	- Up to 3 primary colored emitters.
//	- Up to 3 white colored emitters, so a total of 6 colored emitters.
//	- Custom transfer functions.
//	- Limit of the sum of DCS values.
type ModuleDescriptorGeneral struct {
	// The XYZ values of the primary channels that span the gamut for this module.
	// All colors that increase the gamut have to go into here.
	// This usually describes the 3 primary colors.
	PrimaryColors TransformationLinDCSToXYZ

	// TODO: Precalculate inverse transformations

	// The XYZ values of channels used to increase the CRI and light output of a lamp.
	// This is achieved by maximizing these channels in a way that doesn't change the color or luminance of the resulting color.
	// Normally these are the white LEDs which span a color gamut that lies inside that of the PrimaryColors.
	WhiteColors TransformationLinDCSToXYZ

	// The limit of the sum of all values in the linearized device color space.
	// For other devices this may need some generalization, maybe a custom function/interface that can enforce this limit.
	LinearDCSSumLimit float64

	// Transfer function to convert from a linear device color space into a non linear device color space, and vice versa.
	// Set to nil if your DCS is linear.
	TransferFunction TransferFunction
}

// Check if it implements ModuleDescriptor.
var _ ModuleDescriptor = &ModuleDescriptorGeneral{}

// Channels returns the dimensionality of the device color space.
func (e *ModuleDescriptorGeneral) DCSChannels() int {
	return len(e.PrimaryColors) + len(e.WhiteColors)
}

// AllChannels returns a transformation that contains all channels (A list of all colors).
func (e *ModuleDescriptorGeneral) AllChannels() TransformationLinDCSToXYZ {
	result := make(TransformationLinDCSToXYZ, 0, e.DCSChannels())
	return append(append(result, e.PrimaryColors...), e.WhiteColors...)
}

// XYZToDCS takes a color and returns a vector in the device color space that reproduces the given color as closely as possible.
//
// Short: XYZ --> device color space.
func (e *ModuleDescriptorGeneral) XYZToDCS(color CIE1931XYZColor) (DCSColor, error) {
	// TODO: Precalculate inverse transformations
	primaryTransMatrix, err := e.PrimaryColors.Inverted()
	if err != nil {
		return nil, fmt.Errorf("failed to invert primary color matrix: %w", err)
	}
	whiteTransMatrix, err := e.WhiteColors.Inverted()
	if err != nil {
		return nil, fmt.Errorf("failed to invert white color matrix: %w", err)
	}

	// Determine the DCS values of the primary colors.
	primaryV := primaryTransMatrix.Multiplied(color)

	// Determine the closest possible DCS values of the white colors.
	whiteV := whiteTransMatrix.Multiplied(color).ClampedIndividually()
	// Get the color of whiteValues
	whiteColor, err := e.WhiteColors.Multiplied(whiteV)
	if err != nil {
		return nil, err // Shouldn't happen.
	}
	// Get the proportions of the primary colors that represent whiteColor.
	// Can contain negative values if the white is outside the gamut of the primaries.
	whiteVInPrimary := primaryTransMatrix.Multiplied(whiteColor)

	// The following is just a question of how to weight the different values.
	// We want to increase the CRI and light output.
	// The primary values are decreased as the whites are increased, in a way that doesn't change the total luminance or color output.

	whiteScaling, err := primaryV.ScaledToPositiveDifference(whiteVInPrimary)
	if err != nil {
		return nil, fmt.Errorf("failed to find white color scaling factor: %w", err)
	}

	primaryV, err = primaryV.Sum(whiteVInPrimary.Scaled(-whiteScaling))
	if err != nil {
		return nil, err // Shouldn't happen.
	}
	whiteV = whiteV.Scaled(whiteScaling)

	// Clamp values, apply transfer function.
	resultV := make(DCSColor, 0, primaryV.Channels()+whiteV.Channels())
	resultV = append(resultV, primaryV.DeLinearized(e.TransferFunction)...)
	resultV = append(resultV, whiteV.DeLinearized(e.TransferFunction)...)

	return resultV, nil
}

// DCSToXYZ takes a vector from the device color space and returns the color it represents.
//
// Short: Device color space --> XYZ.
func (e *ModuleDescriptorGeneral) DCSToXYZ(v DCSColor) (CIE1931XYZColor, error) {
	if v.Channels() != e.DCSChannels() {
		return CIE1931XYZColor{}, fmt.Errorf("unexpected amount of channels. Got %d, want %d", v.Channels(), e.DCSChannels())
	}

	linV := v.Linearized(e.TransferFunction)
	linSum := linV.ComponentSum()

	// Scale it so that the sum of all values doesn't exceed the limit.
	if linSum > e.LinearDCSSumLimit {
		linV = linV.Scaled(float64(e.LinearDCSSumLimit) / float64(linSum))
	}

	// Calculate resulting color.
	// This is just the linear combination of all column vectors of the transformation matrix.

	result := CIE1931XYZColor{}

	if color, err := e.AllChannels().Multiplied(linV); err != nil {
		return CIE1931XYZColor{}, fmt.Errorf("failed to multiply transformation matrix with a linear device color space vector: %w", err)
	} else {
		result = result.Sum(color)
	}

	return result, nil
}
