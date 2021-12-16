// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package light

import "github.com/Dadido3/D3iot/light/emission"

// Light is the common interface for all light devices.
// It defines the basic methods to set or get colors, and query the device for its abilities.
//
// A light device contains a number of modules, but at least 1.
// A module is a set of light emitting things that together create a single color impression.
//
// If you need more control over the device, you can use a type assertion to get an implementation of the Light interface.
type Light interface {
	// SetColors sets the emission values of all the modules in the light device.
	// Values which are not set are assumed to equal a turned off module.
	// This will return an error if you try to set more values than there are modules in a light device.
	SetColors(emissionValues ...emission.ValueIntoDCS) error

	// GetColors queries the light device for all emission values of its modules and writes them back into the given list emissionValues.
	// This will return an error if you try to get more values than there are modules in a light device.
	GetColors(emissionValues ...emission.ValueFromDCS) error

	// Modules returns the amount of modules.
	// All devices have at least one module, but most devices have just one.
	Modules() int

	// ColorProfiles returns the color profiles of every module in this device.
	// The length of the resulting list is always equal to the number of modules for this device.
	ColorProfiles() []emission.ColorProfile
}
