// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package emission

import (
	"fmt"
	"math"
)

// CIE1931XYZColor represents a color according to the CIE 1931 XYZ standard.
//
// The luminance Y is not relative, but absolute.
// And contrary to the normal definition, it uses lumen as its unit.
//
// An equal-energy radiator would result in X == Y == Z.
type CIE1931XYZColor struct {
	X, Y, Z float64
}

// Sum returns the sum of c and all colors.
func (c CIE1931XYZColor) Sum(colors ...CIE1931XYZColor) CIE1931XYZColor {
	result := c
	for _, color := range colors {
		result.X += color.X
		result.Y += color.Y
		result.Z += color.Z
	}
	return result
}

// Scaled returns c scaled by the scalar s.
func (c CIE1931XYZColor) Scaled(s float64) CIE1931XYZColor {
	return CIE1931XYZColor{c.X * s, c.Y * s, c.Z * s}
}

// CrossProd returns the cross product between two color vectors.
func (c CIE1931XYZColor) CrossProd(c2 CIE1931XYZColor) CIE1931XYZColor {
	return CIE1931XYZColor{c.Y*c2.Z - c.Z*c2.Y, c.Z*c2.X - c.X*c2.Z, c.X*c2.Y - c.Y*c2.X}
}

// DCSColor represents a color in a device color space.
// This is more or less what is sent to the light device.
//
// If clamped, the values are in the range [0, 1].
// They can be unclamped, that depends on the context where they are used.
//
// Example: 5 channels could represent RGB + cold white + warm white.
type DCSColor []float64

// Channels returns the amount of channels.
// This is the dimensionality of the DCS.
func (c DCSColor) Channels() int {
	return len(c)
}

// ClampedIndividually returns all channels individually clamped into the range [0, 1].
//
//	DCSColor{1.1, 0.9} --> DCSColor{1.0, 0.9}
func (c DCSColor) ClampedIndividually() DCSColor {
	result := make(DCSColor, 0, c.Channels())
	for _, channel := range c {
		result = append(result, clamp01(channel))
	}
	return result
}

// Linearized returns c transformed into linear device color space by the given transfer function tf.
// This clamps the values before transforming them.
func (c DCSColor) Linearized(tf TransferFunction) LinDCSColor {
	if tf != nil {
		return tf.Linearize(c.ClampedIndividually())
	}

	// Transfer function is linear.
	return LinDCSColor(c.ClampedIndividually())
}

// LinDCSColor represents a color in a linear device color space.
// This is more or less what is sent to the light device, but linearized.
//
// If clamped, the values are in the range [0, 1].
// They can be unclamped, that depends on the context where they are used.
type LinDCSColor []float64

// Channels returns the amount of channels.
// This is the dimensionality of the DCS.
func (c LinDCSColor) Channels() int {
	return len(c)
}

// ClampedIndividually returns all channels individually clamped into the range [0, 1].
//
//	LinDCSColor{1.1, 0.9} --> LinDCSColor{1.0, 0.9}
func (c LinDCSColor) ClampedIndividually() LinDCSColor {
	result := make(LinDCSColor, 0, c.Channels())
	for _, channel := range c {
		result = append(result, clamp01(channel))
	}
	return result
}

// ClampedToPositive returns all channels individually clamped into the range [0, +inf].
//
//	LinDCSColor{1.1, 0.9, -0.1} --> LinDCSColor{1.1, 0.9, 0.0}
func (c LinDCSColor) ClampedToPositive() LinDCSColor {
	result := make(LinDCSColor, 0, c.Channels())
	for _, channel := range c {
		if channel >= 0 {
			result = append(result, channel)
		} else {
			result = append(result, 0)
		}
	}
	return result
}

// DeLinearized returns c transformed into device color space by the given transfer function tf.
// This clamps the values before transforming them.
func (c LinDCSColor) DeLinearized(tf TransferFunction) DCSColor {
	if tf != nil {
		return tf.DeLinearize(c.ClampedIndividually())
	}

	// Transfer function is linear.
	return DCSColor(c.ClampedIndividually())
}

// ComponentSum returns the sum of all components.
// No clamp is applied.
func (c LinDCSColor) ComponentSum() float64 {
	var result float64
	for _, channel := range c {
		result += channel
	}
	return result
}

// Sum returns the sum of c and all other colors.
func (c LinDCSColor) Sum(colors ...LinDCSColor) (LinDCSColor, error) {
	result := make(LinDCSColor, c.Channels())
	copy(result, c)

	for _, color := range colors {
		if c.Channels() != color.Channels() {
			return nil, fmt.Errorf("mismatching amount of channels %d and %d", c.Channels(), color.Channels())
		}
		for i, channel := range color {
			result[i] += channel
		}
	}

	return result, nil
}

// Scaled returns c scaled by the scalar s.
func (c LinDCSColor) Scaled(s float64) LinDCSColor {
	result := make(LinDCSColor, 0, c.Channels())
	for _, channel := range c {
		result = append(result, channel*s)
	}
	return result
}

// ScaledToPositiveDifference returns a scaling factor s in a way so that c - c2*s doesn't result in any negative channel values.
// The result is clamped to [0, 1]
// TODO: Find better name, there must be some mathematical concept that describes this
func (c LinDCSColor) ScaledToPositiveDifference(c2 LinDCSColor) (float64, error) {
	if c.Channels() != c2.Channels() {
		return 0, fmt.Errorf("mismatching amount of channels %d and %d", c.Channels(), c2.Channels())
	}

	sMin := 1.0

	for i, channel := range c {
		s := channel / c2[i]
		if sMin > s && !math.IsNaN(s) {
			sMin = s
		}
	}

	return clamp01(sMin), nil
}
