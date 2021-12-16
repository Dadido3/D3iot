// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package emission

import (
	"fmt"
)

// ColorProfileGeneral is a general implementation of a ColorProfile.
// It supports:
//
//	- Up to 3 primary colored emitters.
//	- Up to 3 white emitters, so a total of 6 emitters.
//	- Custom transfer functions.
//	- Custom output limiter functions.
//
// You need to call the Init() or MustInit() method before this profile can be used.
type ColorProfileGeneral struct {
	// WhitePointColor is the brightest color that the module can output.
	// Usually it's the combination of all white emitters.
	// Or of all primary emitters, if there are no white ones, or if the whites are are less bright.
	WhitePointColor CIE1931XYZAbs

	// The XYZ values of the primary channels that span the gamut for this module.
	// All colors that increase the gamut have to go into here.
	// This usually describes the 3 primary colors.
	PrimaryColors TransformationLinDCSToXYZ
	// Precalculated inverse of the above.
	invPrimaryColors TransformationXYZToLinDCS

	// The XYZ values of channels used to increase the CRI and light output of a lamp.
	// This is achieved by maximizing these channels in a way that doesn't change the color or luminance of the resulting color.
	// Normally these are the white LEDs which span a color gamut that lies inside that of the PrimaryColors.
	WhiteColors TransformationLinDCSToXYZ
	// Precalculated inverse of the above.
	invWhiteColors TransformationXYZToLinDCS

	// Limit the output in some form.
	OutputLimiter OutputLimiter

	// Transfer function to convert from a linear device color space into a non linear device color space, and vice versa.
	// Set to nil if your DCS is linear.
	TransferFunc TransferFunction
}

// Check if it implements ColorProfile.
var _ ColorProfile = &ColorProfileGeneral{}

// Init precalculates some values.
func (e *ColorProfileGeneral) Init() error {
	var err error

	if e.invPrimaryColors, err = e.PrimaryColors.Inverted(); err != nil {
		return fmt.Errorf("failed to invert primary color matrix: %w", err)
	}

	if e.invWhiteColors, err = e.WhiteColors.Inverted(); err != nil {
		return fmt.Errorf("failed to invert white color matrix: %w", err)
	}

	return nil
}

// MustInit is the same as Init(), but panics on any error.
// As a small help, this returns the module itself.
func (e *ColorProfileGeneral) MustInit() *ColorProfileGeneral {
	if err := e.Init(); err != nil {
		panic(fmt.Sprintf("Failed to init module profile %v: %v", e, err))
	}

	return e
}

// Channels returns the dimensionality of the device color space.
func (e *ColorProfileGeneral) Channels() int {
	return len(e.PrimaryColors) + len(e.WhiteColors)
}

// WhitePoint returns the white point as a CIE 1931 XYZ color.
// This is also the brightest color a module can output.
func (e *ColorProfileGeneral) WhitePoint() CIE1931XYZAbs {
	return e.WhitePointColor
}

// ChannelPoints returns a list of channel colors.
// Depending on the module type, this could be the colors for:
//
//	- Single white emitter.
//	- RGB emitters.
//	- RGB + white emitters.
//	- RGB + cold white + warm white emitters.
func (e *ColorProfileGeneral) ChannelPoints() []CIE1931XYZAbs {
	return e.FullTransformation()
}

// FullTransformation returns a transformation (matrix) that contains all channels (A list of all colors).
func (e *ColorProfileGeneral) FullTransformation() TransformationLinDCSToXYZ {
	result := make(TransformationLinDCSToXYZ, 0, e.Channels())
	return append(append(result, e.PrimaryColors...), e.WhiteColors...)
}

// XYZToDCS takes a color and returns a vector in the device color space that reproduces the given color as closely as possible.
//
// Short: XYZ --> device color space.
func (e *ColorProfileGeneral) XYZToDCS(color CIE1931XYZAbs) DCSVector {
	// Determine the DCS values of the primary channels.
	primaryV := e.invPrimaryColors.Multiplied(color)

	// Determine the closest possible DCS values of the white channels.
	whiteV := e.invWhiteColors.Multiplied(color).ClampedToPositive()
	// Get the color of whiteValues
	whiteColor, err := e.WhiteColors.Multiplied(whiteV)
	if err != nil {
		panic(err) // Shouldn't happen.
	}
	// Get the proportions of the primary colors that represent whiteColor.
	// Can contain negative values if the white is outside the gamut of the primaries.
	whiteVInPrimary := e.invPrimaryColors.Multiplied(whiteColor)

	// The following is just a question of how to weight the different values.
	// We want to increase the CRI and light output.
	// The primary values are decreased as the whites are increased, in a way that doesn't change the total luminance or color output.

	whiteScaling, err := primaryV.ScaledToPositiveDifference(whiteVInPrimary)
	if err != nil {
		panic(err) // Shouldn't happen.
	}

	primaryV, err = primaryV.Sum(whiteVInPrimary.Scaled(-whiteScaling))
	if err != nil {
		panic(err) // Shouldn't happen.
	}
	whiteV = whiteV.Scaled(whiteScaling)

	// Put all values into one slice.
	linV := make(LinDCSVector, 0, primaryV.Channels()+whiteV.Channels())
	linV = append(append(linV, primaryV...), whiteV...)

	// Limit output.
	if e.OutputLimiter != nil {
		linV = e.OutputLimiter.LimitDCS(linV)
	}

	// Clamp values, apply transfer function.
	nonLinV := linV.ClampedAndDeLinearized(e.TransferFunc)

	return nonLinV
}

// DCSToXYZ takes a vector from the device color space and returns the color it represents.
//
// Short: Device color space --> XYZ.
func (e *ColorProfileGeneral) DCSToXYZ(v DCSVector) (CIE1931XYZAbs, error) {
	if v.Channels() != e.Channels() {
		return CIE1931XYZAbs{}, fmt.Errorf("unexpected amount of channels. Got %d, want %d", v.Channels(), e.Channels())
	}

	linV := v.ClampedAndLinearized(e.TransferFunc)

	if e.OutputLimiter != nil {
		linV = e.OutputLimiter.LimitDCS(linV)
	}

	// Calculate resulting color.
	// This is just the linear combination of all column vectors of the transformation matrix.

	result := CIE1931XYZAbs{}

	if color, err := e.FullTransformation().Multiplied(linV); err != nil {
		return CIE1931XYZAbs{}, fmt.Errorf("failed to multiply transformation matrix with a linear device color space vector: %w", err)
	} else {
		result = result.Sum(color)
	}

	return result, nil
}

func (e *ColorProfileGeneral) TransferFunction() TransferFunction {
	return e.TransferFunc
}
