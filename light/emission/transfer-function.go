// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package emission

import "math"

// TransferFunction is used to transform from device color spaces into linear device colors spaces, and vice versa.
type TransferFunction interface {
	Linearize(v DCSColor) LinDCSColor
	DeLinearize(v LinDCSColor) DCSColor
}

// transferFunctionStandardRGB implements the sRGB transfer function.
type transferFunctionStandardRGB struct{}

// TransferFunctionStandardRGB implements the sRGB transfer function.
var TransferFunctionStandardRGB = &transferFunctionStandardRGB{}

func (tf *transferFunctionStandardRGB) Linearize(values DCSColor) LinDCSColor {
	result := make(LinDCSColor, 0, values.Channels())
	for _, value := range values {
		var trans float64

		if value <= 0.04045 {
			trans = value / 12.92
		} else {
			trans = math.Pow((value+0.055)/1.055, 2.4)
		}

		result = append(result, trans)
	}
	return result
}

func (tf *transferFunctionStandardRGB) DeLinearize(values LinDCSColor) DCSColor {
	result := make(DCSColor, 0, values.Channels())
	for _, value := range values {
		var trans float64

		if value <= 0.0031308 {
			trans = 12.92 * value
		} else {
			trans = 1.055*math.Pow(value, 1/2.4) - 0.055
		}

		result = append(result, trans)
	}
	return result
}
