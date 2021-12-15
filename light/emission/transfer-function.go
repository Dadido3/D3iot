// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package emission

import "math"

// TransferFunction is used to transform from device color spaces into linear device colors spaces, and vice versa.
type TransferFunction interface {
	Linearize(v DCSVector) LinDCSVector
	DeLinearize(v LinDCSVector) DCSVector
}

// transferFunctionStandardRGB implements the sRGB transfer function.
type transferFunctionStandardRGB struct{}

// TransferFunctionStandardRGB implements the sRGB transfer function.
var TransferFunctionStandardRGB = transferFunctionStandardRGB{}

func (tf transferFunctionStandardRGB) Linearize(vector DCSVector) LinDCSVector {
	result := make(LinDCSVector, 0, vector.Channels())
	for _, channel := range vector {
		var trans float64

		if channel <= 0.04045 {
			trans = channel / 12.92
		} else {
			trans = math.Pow((channel+0.055)/1.055, 2.4)
		}

		result = append(result, trans)
	}
	return result
}

func (tf transferFunctionStandardRGB) DeLinearize(vector LinDCSVector) DCSVector {
	result := make(DCSVector, 0, vector.Channels())
	for _, channel := range vector {
		var trans float64

		if channel <= 0.0031308 {
			trans = 12.92 * channel
		} else {
			trans = 1.055*math.Pow(channel, 1/2.4) - 0.055
		}

		result = append(result, trans)
	}
	return result
}

// TransferFunctionGamma implements the a gamma transfer function.
type TransferFunctionGamma struct {
	Gamma float64
}

func (tf TransferFunctionGamma) Linearize(vector DCSVector) LinDCSVector {
	result := make(LinDCSVector, 0, vector.Channels())
	for _, channel := range vector {
		result = append(result, math.Pow(channel, tf.Gamma))
	}
	return result
}

func (tf TransferFunctionGamma) DeLinearize(vector LinDCSVector) DCSVector {
	result := make(DCSVector, 0, vector.Channels())
	for _, channel := range vector {
		result = append(result, math.Pow(channel, 1/tf.Gamma))
	}
	return result
}
