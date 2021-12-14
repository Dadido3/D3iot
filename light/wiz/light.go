// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package wiz

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/Dadido3/D3iot/light"
	"github.com/Dadido3/D3iot/light/emission"
)

type Light struct {
	address string

	// The product describing the device.
	// Either an exact match or a general product that may fit good enough.
	// This must not be nil.
	product *Product

	deadline  time.Duration // Default deadline for a whole communcation action (sending and receiving).
	retries   uint          // Number of retries when the deadline got exceeded.
	connMutex sync.Mutex    // Mutex preventing simultaneous connections to this device.
	//paramMutex sync.Mutex    // Mutex protecting parameters of this object.

	DebugWriter io.Writer // Writer that can be used to debug network communication.
}

// Check implementation of light.Light.
var _ light.Light = &Light{}

// NewLight returns an object that represents a single WiZ light accessible by the given address.
//
// This will query the product type, so it needs to be able to connect via the given address.
//
//	light, err := NewLight("192.168.1.123:38899")
func NewLight(address string) (*Light, error) {
	light := &Light{
		address:  address,
		deadline: 100 * time.Millisecond,
		retries:  10,
	}

	var err error
	if light.product, err = light.determineProduct(); err != nil {
		return nil, fmt.Errorf("couldn't determine WiZ product: %w", err)
	}

	return light, nil
}

// NewLight returns an object that represents a single WiZ light accessible by the given address.
//
// This will not query the device to determine the WiZ product, but use the one defined in the parameter.
// Therefore it will not make an attempt to communicate with the light.
func NewLightWithProduct(address string, product *Product) (*Light, error) { // TODO: Find a way to make a more generalized version of this. There need to be some way to create a light object without having to connect to determine the product
	if product == nil {
		return nil, fmt.Errorf("no product defined")
	}

	light := &Light{
		address:  address,
		deadline: 100 * time.Millisecond,
		retries:  10,
		product:  product,
	}

	return light, nil
}

// determineProduct queries and determines the product of the device.
func (l *Light) determineProduct() (*Product, error) {
	// Query device info from lamp.
	devInfo, err := l.GetDeviceInfo()
	if err != nil {
		return nil, err
	}

	product, err := determineProduct(devInfo.ModuleName)
	if err != nil {
		return nil, err
	}

	return product, nil
}

// Product returns an exact or general product descriptor of the device's abilities and limits.
func (l *Light) Product() *Product {
	return l.product
}

// SetColors sets the colors of all the modules in the light device.
// Colors which are not set are assumed to be turned off.
// This will return an error if you try to set more colors than there are modules in a light device.
func (l *Light) SetColors(colors ...emission.Value) error {
	switch len(colors) {
	case 0:
		pilot := NewPilot(false)
		return l.SetPilot(pilot)

	case 1:
		moduleProfile := l.ModuleProfiles()[0]

		// Transform dcsColor into DCS.
		dcsColor := colors[0].DCSColor(moduleProfile)

		switch dc := l.product.deviceClass; dc {
		case deviceClassDW:
			if dcsColor.Channels() == 1 {
				return l.SetPilot(NewPilot(true).WithDimming(uint(dcsColor[0] * 100)))
			} else {
				return fmt.Errorf("unexpected number of channels. Got %d, want %d", dcsColor.Channels(), 1)
			}

		case deviceClassTW:
			if dcsColor.Channels() == 2 {
				return l.SetPilot(NewPilotWithRGBW(100, 0, 0, 0, uint8(dcsColor[0]*255), uint8(dcsColor[1]*255)))
			} else {
				return fmt.Errorf("unexpected number of channels. Got %d, want %d", dcsColor.Channels(), 2)
			}

		case deviceClassRGBTW:
			if dcsColor.Channels() == 5 {
				return l.SetPilot(NewPilotWithRGBW(100, uint8(dcsColor[0]*255), uint8(dcsColor[1]*255), uint8(dcsColor[2]*255), uint8(dcsColor[3]*255), uint8(dcsColor[4]*255)))
			} else {
				return fmt.Errorf("unexpected number of channels. Got %d, want %d", dcsColor.Channels(), 5)
			}

		default:
			return fmt.Errorf("unsupported device class %q", dc)

		}

	default:
		return fmt.Errorf("supplied %d colors, this device has only 1 module", len(colors))
	}
}

// GetColors queries the light device for all colors of its modules and returns them.
// This must always returns as much elements as the device has modules, even in case of an error.
func (l *Light) GetColors() ([]emission.CIE1931XYZColor, error) {
	pilot, err := l.GetPilot()
	if err != nil {
		return []emission.CIE1931XYZColor{{}}, fmt.Errorf("couldn't read pilot: %w", err)
	}

	moduleProfile := l.ModuleProfiles()[0]

	// Generate DCS color/vector.
	var dcsColor emission.DCSColor
	switch dc := l.product.deviceClass; dc {
	case deviceClassDW:
		if pilot.State && pilot.Dimming != nil {
			dcsColor = emission.DCSColor{float64(*pilot.Dimming) / 100}
		}

	case deviceClassTW:
		if pilot.State && pilot.CW != nil && pilot.WW != nil {
			dcsColor = emission.DCSColor{float64(*pilot.CW) / 255, float64(*pilot.WW) / 255}
		}

	case deviceClassRGBTW:
		if pilot.State && pilot.R != nil && pilot.G != nil && pilot.B != nil && pilot.CW != nil && pilot.WW != nil {
			dcsColor = emission.DCSColor{float64(*pilot.R) / 255, float64(*pilot.G) / 255, float64(*pilot.B) / 255, float64(*pilot.CW) / 255, float64(*pilot.WW) / 255}
		}

	default:
		return []emission.CIE1931XYZColor{{}}, fmt.Errorf("unsupported device class %q", dc)

	}

	xyzColor, err := moduleProfile.DCSToXYZ(dcsColor)
	return []emission.CIE1931XYZColor{xyzColor}, err
}

// Modules returns the amount of modules.
// All devices must at least have one module.
func (l *Light) Modules() int {
	return 1 // Most (or all?) WiZ lights have one module.
}

// ModuleProfiles returns the profiles that describe the modules of a device.
// The length of the resulting list is equal to the number of modules for this device.
// This must always returns as much elements as the device has modules.
func (l *Light) ModuleProfiles() []emission.ModuleProfile {
	return []emission.ModuleProfile{l.Product().moduleProfile}
}
