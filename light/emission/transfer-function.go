// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package emission

// TransferFunction is used to transform from device color spaces into linear device colors spaces, and vice versa.
type TransferFunction interface {
	Linearize(v DCSColor) LinDCSColor
	DeLinearize(v LinDCSColor) DCSColor
}
