// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package light

import "github.com/Dadido3/D3iot/light/emission"

type Light interface {
	// SetColors sets the colors of all the modules in the light device.
	SetColors(colors ...emission.Value) error

	// Colors returns the current colors of all the modules in the light device.
	Colors() ([]emission.CIE1931XYZColor, error)

	// ModuleProfiles returns the profiles that describe each module.
	ModuleProfiles() ([]emission.ModuleProfile, error)
}
