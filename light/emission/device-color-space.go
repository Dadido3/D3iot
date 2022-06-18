// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package emission

import (
	"fmt"
	"math"
)

// DCSVector represents a color in a device color space.
// This is more or less what is sent to the light device.
//
// If clamped, the values are in the range [0, 1].
// They may or may not be clamped, that depends on the context where they are used.
//
// Example: 5 channels could represent RGB + cold white + warm white.
type DCSVector []float64

var _ Value = &DCSVector{}

// Copy returns a copy of v.
func (v DCSVector) Copy() DCSVector {
	vCopy := make(DCSVector, v.Channels())
	copy(vCopy, v)
	return vCopy
}

// IntoDCS implements the Value interface.
func (v DCSVector) IntoDCS(cp ColorProfile) DCSVector {
	return v
}

// FromDCS implements the Value interface.
func (v *DCSVector) FromDCS(cp ColorProfile, v2 DCSVector) error {
	*v = v2.Copy()
	return nil
}

// Channels returns the number of channels.
// This is the dimensionality of the DCS.
func (v DCSVector) Channels() int {
	return len(v)
}

// ClampedIndividually returns all channels individually clamped into the range [0, 1].
//
//	DCSVector{1.1, 0.9, -0.1} --> DCSVector{1.0, 0.9, 0.0}
func (v DCSVector) ClampedIndividually() DCSVector {
	result := v.Copy()
	for i, channel := range result {
		result[i] = clamp01(channel)
	}
	return result
}

// ClampedAndLinearized returns v transformed into linear device color space by the given transfer function tf.
// This clamps the values before transforming them.
func (v DCSVector) ClampedAndLinearized(tf TransferFunction) LinDCSVector {
	if tf != nil {
		return tf.Linearize(v.ClampedIndividually())
	}

	// Transfer function is linear.
	return LinDCSVector(v.ClampedIndividually())
}

// Difference returns the difference v - v2.
//
// This may or may not make sense to use, as this is not a linear space.
func (v DCSVector) Difference(v2 DCSVector) (DCSVector, error) {
	if v.Channels() != v2.Channels() {
		return nil, fmt.Errorf("mismatching number of channels %d and %d", v.Channels(), v2.Channels())
	}

	result := v.Copy()
	for i, channel := range result {
		result[i] = channel - v2[i]
	}

	return result, nil
}

// ComponentSum returns the sum of all components.
// No clamp is applied.
//
// This may or may not make sense to use, as this is not a linear space.
func (v DCSVector) ComponentSum() float64 {
	var result float64
	for _, channel := range v {
		result += channel
	}
	return result
}

// LinDCSVector represents a color in a linear device color space.
// This is more or less what is sent to the light device, but linearized.
//
// If clamped, the values are in the range [0, 1].
// They may or may not be clamped, that depends on the context where they are used.
type LinDCSVector []float64

var _ Value = &LinDCSVector{}

// Copy returns a copy of v.
func (v LinDCSVector) Copy() LinDCSVector {
	vCopy := make(LinDCSVector, v.Channels())
	copy(vCopy, v)
	return vCopy
}

// IntoDCS implements the Value interface.
func (v LinDCSVector) IntoDCS(cp ColorProfile) DCSVector {
	return v.ClampedAndDeLinearized(cp.TransferFunction())
}

// FromDCS implements the Value interface.
func (v *LinDCSVector) FromDCS(cp ColorProfile, v2 DCSVector) error {
	*v = v2.ClampedAndLinearized(cp.TransferFunction())
	return nil
}

// Channels returns the number of channels.
// This is the dimensionality of the DCS.
func (v LinDCSVector) Channels() int {
	return len(v)
}

// ClampedIndividually returns all channels individually clamped into the range [0, 1].
//
//	LinDCSVector{1.1, 0.9, -0.1} --> LinDCSVector{1.0, 0.9, 0.0}
func (v LinDCSVector) ClampedIndividually() LinDCSVector {
	result := make(LinDCSVector, 0, v.Channels())
	for _, channel := range v {
		result = append(result, clamp01(channel))
	}
	return result
}

// ClampedUniform returns all channels scaled by a single scaling factor in a way so that every channel is below or equal to one.
//
//	LinDCSVector{1.1, 0.9, -0.1} --> LinDCSVector{1.0, 0.818, -0.091}
func (v LinDCSVector) ClampedUniform() LinDCSVector {
	scale := 1.0
	for _, channel := range v {
		neededScale := 1.0 / channel
		if channel > 1.0 && scale > neededScale {
			scale = neededScale
		}
	}
	return v.Scaled(scale)
}

// ClampedToPositive returns all channels individually clamped into the range [0, +inf].
//
//	LinDCSVector{1.1, 0.9, -0.1} --> LinDCSVector{1.1, 0.9, 0.0}
func (v LinDCSVector) ClampedToPositive() LinDCSVector {
	result := make(LinDCSVector, 0, v.Channels())
	for _, channel := range v {
		if channel >= 0 {
			result = append(result, channel)
		} else {
			result = append(result, 0)
		}
	}
	return result
}

// ClampedAndDeLinearized returns v transformed into device color space by the given transfer function tf.
// This clamps the values before transforming them.
func (v LinDCSVector) ClampedAndDeLinearized(tf TransferFunction) DCSVector {
	if tf != nil {
		return tf.DeLinearize(v.ClampedIndividually())
	}

	// Transfer function is linear.
	return DCSVector(v.ClampedIndividually())
}

// ComponentSum returns the sum of all components.
// No clamp is applied.
func (v LinDCSVector) ComponentSum() float64 {
	var result float64
	for _, channel := range v {
		result += channel
	}
	return result
}

// Sum returns the sum of v and all other vectors.
func (v LinDCSVector) Sum(vectors ...LinDCSVector) (LinDCSVector, error) {
	result := v.Copy()

	for _, vector := range vectors {
		if v.Channels() != vector.Channels() {
			return nil, fmt.Errorf("mismatching number of channels %d and %d", v.Channels(), vector.Channels())
		}
		for i, channel := range vector {
			result[i] += channel
		}
	}

	return result, nil
}

// Difference returns the difference v - v2.
//
// This may or may not make sense to use, as this is not a linear space.
func (v LinDCSVector) Difference(v2 LinDCSVector) (LinDCSVector, error) {
	if v.Channels() != v2.Channels() {
		return nil, fmt.Errorf("mismatching number of channels %d and %d", v.Channels(), v2.Channels())
	}

	result := make(LinDCSVector, 0, v.Channels())
	for i, channel := range v {
		result = append(result, channel-v2[i])
	}

	return result, nil
}

// Scaled returns v scaled by the scalar s.
func (v LinDCSVector) Scaled(s float64) LinDCSVector {
	result := make(LinDCSVector, 0, v.Channels())
	for _, channel := range v {
		result = append(result, channel*s)
	}
	return result
}

// ScaledToPositiveDifference returns the largest scaling factor s in a way so that v - v2*s doesn't result in any negative channel values.
// The result is clamped to [0, 1].
// TODO: Find better name, there must be some mathematical concept that describes this
func (v LinDCSVector) ScaledToPositiveDifference(v2 LinDCSVector) (float64, error) {
	if v.Channels() != v2.Channels() {
		return 0, fmt.Errorf("mismatching number of channels %d and %d", v.Channels(), v2.Channels())
	}

	sMin := 1.0

	for i, channel := range v {
		s := channel / v2[i]
		if sMin > s && !math.IsNaN(s) {
			sMin = s
		}
	}

	return clamp01(sMin), nil
}
