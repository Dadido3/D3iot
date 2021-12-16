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

// SetColors sets the emission values of all the modules in the light device.
// Values which are not set are assumed to equal a turned off module.
// This must return an error if there are more values than there are modules in a light device.
func (l *Light) SetColors(emissionValues ...emission.ValueIntoDCS) error {
	switch len(emissionValues) {
	case 0:
		pilot := NewPilot(false)
		return l.SetPilot(pilot)

	case 1:
		colorProfile := l.ColorProfiles()[0]

		// Transform emission value into DCS.
		vector := emissionValues[0].IntoDCS(colorProfile)

		switch dc := l.product.deviceClass; dc {
		case deviceClassDW:
			if vector.Channels() == 1 {
				dimming := uint(normFloatToInt(vector[0], 100))
				return l.SetPilot(NewPilot(true).WithScene(SceneCoolWhite, 100).WithDimming(dimming))
			} else {
				return fmt.Errorf("unexpected number of channels. Got %d, want %d", vector.Channels(), 1)
			}

		case deviceClassTW:
			if vector.Channels() == 2 {
				cw, ww := normFloatToUint8(vector[0]), normFloatToUint8(vector[1])
				return l.SetPilot(NewPilotWithWhite(100, cw, ww))
			} else {
				return fmt.Errorf("unexpected number of channels. Got %d, want %d", vector.Channels(), 2)
			}

		case deviceClassRGBTW:
			if vector.Channels() == 5 {
				r, g, b, cw, ww := normFloatToUint8(vector[0]), normFloatToUint8(vector[1]), normFloatToUint8(vector[2]), normFloatToUint8(vector[3]), normFloatToUint8(vector[4])
				return l.SetPilot(NewPilotWithRGBW(100, r, g, b, cw, ww))
			} else {
				return fmt.Errorf("unexpected number of channels. Got %d, want %d", vector.Channels(), 5)
			}

		default:
			return fmt.Errorf("unsupported device class %q", dc)

		}

	default:
		return fmt.Errorf("got %d emission values, this device has only 1 module", len(emissionValues))
	}
}

// GetColors queries the light device for all emission values of its modules and writes them back into the given list emissionValues.
// This must return an error if there are more values than there are modules in a light device.
func (l *Light) GetColors(emissionValues ...emission.ValueFromDCS) error {
	// Check number of emission values.
	switch len(emissionValues) {
	case 0:
		return nil

	case 1:
		// Continue.

	default:
		return fmt.Errorf("got %d emission values, this device has only 1 module", len(emissionValues))
	}

	pilot, err := l.GetPilot()
	if err != nil {
		return fmt.Errorf("couldn't read pilot: %w", err)
	}

	if pilot.Scene != nil {
		return fmt.Errorf("current state can't be represented by a color")
	}

	// Generate DCS color/vector.
	var vector emission.DCSVector
	switch dc := l.product.deviceClass; dc {
	case deviceClassDW:
		if pilot.State && pilot.HasDimming() {
			vector = emission.DCSVector{float64(*pilot.Dimming) / 100}
		} else {
			vector = emission.DCSVector{0}
		}

	case deviceClassTW:
		if pilot.State && pilot.HasWhite() {
			vector = emission.DCSVector{float64(*pilot.CW) / 255, float64(*pilot.WW) / 255}
		} else {
			vector = emission.DCSVector{0, 0}
		}

	case deviceClassRGBTW:
		if pilot.State && pilot.HasRGBW() {
			vector = emission.DCSVector{float64(*pilot.R) / 255, float64(*pilot.G) / 255, float64(*pilot.B) / 255, float64(*pilot.CW) / 255, float64(*pilot.WW) / 255}
		} else {
			vector = emission.DCSVector{0, 0, 0, 0, 0}
		}

	default:
		return fmt.Errorf("unsupported device class %q", dc)

	}

	colorProfile := l.ColorProfiles()[0]
	return emissionValues[0].FromDCS(colorProfile, vector)
}

// Modules returns the number of modules.
// All devices must at least have one module.
func (l *Light) Modules() int {
	return 1 // Most (or all?) WiZ lights have one module.
}

// ColorProfiles returns the color profiles of every module in this device.
// The length of the resulting list must always be equal to the number of modules for this device.
func (l *Light) ColorProfiles() []emission.ColorProfile {
	return []emission.ColorProfile{l.Product().colorProfile}
}
