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
// TODO: Consider refactoring this into an interface. The implementation should be inside of each device driver. This could be a general implementation that devices could use
type ModuleDescriptor struct {
	// The XYZ values of the primary channels that span the gamut for this module.
	// All colors that increase the gamut have to go into here.
	// This usually describes the 3 primary colors.
	PrimaryColors TransformationLinDCSToXYZ

	// Precalculated "inverse primary colors".
	// It's used to transform CIE 1931 XYZ colors into the device color space.
	//
	// This can be generated by
	//InversePrimaryColors TransformationXYZToLinDCS

	// The XYZ values of channels used to increase the CRI and light output of a lamp.
	// This is achieved by maximizing these channels in a way that doesn't change the color or luminance of the resulting color.
	// This usually describes the different white LEDs of a module.
	CRIColors TransformationLinDCSToXYZ

	// Examples:
	// - Single color emitter (E.g. White or some other color): PrimaryColors = ColorPrimariesDCS{{x,y,z}}, CRIColors = ColorPrimariesDCS{}
	// - 2 color emitter (E.g. tunable white): PrimaryColors = ColorPrimariesDCS{{x,y,z}, {x,y,z}}, CRIColors = ColorPrimariesDCS{}
	// - 3 color emitter (E.g. RGB): PrimaryColors = ColorPrimariesDCS{{x,y,z}, {x,y,z}, {x,y,z}}, CRIColors = ColorPrimariesDCS{}
	// - 4 color emitter (E.g. RGBW): PrimaryColors = ColorPrimariesDCS{{x,y,z}, {x,y,z}, {x,y,z}}, CRIColors = ColorPrimariesDCS{{x,y,z}}
	// - 5 color emitter (E.g. RGBCW): PrimaryColors = ColorPrimariesDCS{{x,y,z}, {x,y,z}, {x,y,z}}, CRIColors = ColorPrimariesDCS{{x,y,z}, {x,y,z}}

	// The limit of the sum of all values in the linearized device color space.
	// For other devices this may need some generalization, maybe a custom function/interface that can enforce this limit.
	LinearDCSSumLimit float64

	// TODO: Add transfer functions. They are not needed right now, but may be needed for future devices
}

// Channels returns the dimensionality of the device color space.
func (e *ModuleDescriptor) Channels() int {
	return len(e.PrimaryColors) + len(e.CRIColors)
}

// XYZToDCS takes a color and returns a vector in the device color space that reproduces the given color as closely as possible.
//
// Short: XYZ --> device color space.
func (e *ModuleDescriptor) XYZToDCS(color CIE1931XYZColor) ([]float64, error) {
	transMatrix, err := e.PrimaryColors.Inverted()
	if err != nil {
		return nil, fmt.Errorf("failed to invert primary color matrix: %w", err)
	}

	dcsValues := transMatrix.Multiplied(color)

	// TODO: Add stuff to calculate values for the CRIColors

	// Apply transfer function, clamp values.
	for i, value := range dcsValues {
		// For future: Transfer function math has to be put into here.
		dcsValues[i] = clamp01(value) // TODO: Maybe don't clamp here, but inside the device driver itself
	}

	return dcsValues, nil
}

// DCSToXYZ takes a vector from the device color space and returns the color it represents.
//
// Short: Device color space --> XYZ.
func (e *ModuleDescriptor) DCSToXYZ(dcsValues []float64) (CIE1931XYZColor, error) {
	if len(dcsValues) != e.Channels() {
		return CIE1931XYZColor{}, fmt.Errorf("unexpected amount of channel values. Got %d, want %d", len(dcsValues), e.Channels())
	}

	var sum float64
	dcsLinValues := make([]float64, 0, len(dcsValues))
	for _, dcsValue := range dcsValues {
		dcsValue = clamp01(dcsValue)
		// For future: Transfer function math has to be put here.
		dcsLinValues = append(dcsLinValues, dcsValue)
		sum += dcsValue
	}

	// Scale it so that the sum of all values doesn't exceed the limit.
	if sum > e.LinearDCSSumLimit {
		scale := float64(e.LinearDCSSumLimit) / float64(sum)
		for i := range dcsLinValues {
			dcsLinValues[i] *= scale
		}
	}

	// Calculate resulting color.
	// This is just the linear combination of all column vectors of the transformation matrix.

	result := CIE1931XYZColor{}

	if color, err := e.PrimaryColors.Multiplied(dcsLinValues[:len(e.PrimaryColors)]); err != nil {
		return CIE1931XYZColor{}, fmt.Errorf("failed to multiply transformation matrix with a linear device color space vector: %w", err)
	} else {
		result = result.Add(color)
	}

	if color, err := e.CRIColors.Multiplied(dcsLinValues[len(e.PrimaryColors):]); err != nil {
		return CIE1931XYZColor{}, fmt.Errorf("failed to multiply transformation matrix with a linear device color space vector: %w", err)
	} else {
		result = result.Add(color)
	}

	return result, nil
}
