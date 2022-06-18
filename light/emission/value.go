// Copyright (c) 2021-2022 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package emission

// Value represents a light emission by its color and luminance.
//
// Every type that represents some sort of light emission should implement this interface.
type Value interface {
	IntoDCS(ColorProfile) DCSVector // IntoDCS returns the value transformed into the device color space.
}

// ValueReceiver can receive light emissions via light.GetColors(...).
//
// This means it can take a DCS vector, and transform it into its color space with the help of the color profile.
type ValueReceiver interface {
	FromDCS(ColorProfile, DCSVector) error // FromDCS transforms the device color space vector into the color space of Value.
}
