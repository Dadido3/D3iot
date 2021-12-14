// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package light

import "github.com/Dadido3/D3iot/light/emission"

// Light is the common interface for all light devices.
// It defines the basic methods to set or get colors, and query some of the devices its abilities.
//
// If you more control over the device, you can use a type assertion to get a less general device.
type Light interface {
	// SetColors sets the colors of all the modules in the light device.
	// Colors which are not set are assumed to be turned off.
	// This will return an error if you try to set more colors than there are modules in a light device.
	SetColors(colors ...emission.Value) error

	// GetColors queries the light device for all colors of its modules and returns them.
	// This always returns as much elements as the device as modules, even in case of an error.
	GetColors() ([]emission.CIE1931XYZColor, error)

	// Modules returns the amount of modules.
	// All devices have at least one module, but most devices have just one.
	Modules() int

	// ModuleProfiles returns the profiles that describe the modules of a device.
	// The length of the resulting list is equal to the number of modules for this device.
	// This always returns as much elements as the device as modules.
	ModuleProfiles() []emission.ModuleProfile
}
