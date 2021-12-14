// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package emission

// Value describes a light emission by its color and luminance.
//
// Every type that implements this value can be used to control the output of light devices.
type Value interface {
	DCSColor(ModuleProfile) DCSColor // DCSColor returns the value transformed into the device color space.
}
