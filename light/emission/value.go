// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package emission

// Value describes a light emission by its color and luminance.
//
// Every type that represents some sort of light emission should implement this interface.
type Value interface {
	IntoDCS(ModuleProfile) DCSColor        // IntoDCS returns the value transformed into the device color space.
	FromDCS(ModuleProfile, DCSColor) error // FromDCS transforms the device color space vector into the color space of Value.
}

// ValueIntoDCS is the Value into DCS transformation part of the Value interface.
type ValueIntoDCS interface {
	IntoDCS(ModuleProfile) DCSColor // IntoDCS returns the value transformed into the device color space.
}

// ValueIntoDCS is the DCS into Value transformation part of the Value interface.
type ValueFromDCS interface {
	FromDCS(ModuleProfile, DCSColor) error // FromDCS transforms the device color space vector into the color space of Value.
}
