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

	product *Product // The cached product of the lamp. This will be queried from the device as soon as it is needed.

	deadline   time.Duration // Default deadline for a whole communcation action (sending and receiving).
	retries    uint          // Number of retries when the deadline got exceeded.
	connMutex  sync.Mutex    // Mutex preventing simultaneous connections to this device.
	paramMutex sync.Mutex    // Mutex protecting parameters of this object.

	DebugWriter io.Writer // Writer that can be used to debug network communication.
}

// Check implementation of light.Light.
var _ light.Light = &Light{}

// NewLight returns an object that represents a single WiZ light accessible by the given address.
//
//	light := NewLight("192.168.1.123:38899")
func NewLight(address string) *Light {
	return &Light{
		address:  address,
		deadline: 100 * time.Millisecond,
		retries:  10,
	}
}

// Product returns an exact or general product descriptor of the device's abilities and limits.
func (l *Light) Product() (*Product, error) {
	l.paramMutex.Lock()
	defer l.paramMutex.Unlock()

	// Use cached product if possible.
	if l.product != nil {
		return l.product, nil
	}

	// Query product from lamp.
	devInfo, err := l.GetDeviceInfo()
	if err != nil {
		return nil, err
	}

	product, err := determineProduct(devInfo.ModuleName)
	if err != nil {
		return nil, err
	}

	l.product = product
	return product, nil
}

// SetProduct allows to overwrite the product of the device.
// This is useful if users of this lib are caching the products of known devices on their own.
//
// If you want to cause a re-evaluation, set the product to nil.
func (l *Light) SetProduct(product *Product) {
	l.paramMutex.Lock()
	defer l.paramMutex.Unlock()

	l.product = product
}

// SetColors sets the colors of all the modules in the light device.
// Colors which are not set will be turned off.
func (l *Light) SetColors(colors ...emission.Value) error {
	switch len(colors) {
	case 0:
		pilot := NewPilot(false)
		return l.SetPilot(pilot)

	case 1:
		product, err := l.Product()
		if err != nil {
			return fmt.Errorf("failed to determine WiZ product: %w", err)
		}
		if product.moduleProfile == nil {
			return fmt.Errorf("WiZ product doesn't contain module profile")
		}
		color := colors[0].DCSColor(product.moduleProfile)

		switch color.Channels() {
		case 1:
			return l.SetPilot(NewPilot(true).WithDimming(uint(color[0] * 100)))
		case 2:
			return l.SetPilot(NewPilotWithRGBW(100, 0, 0, 0, uint8(color[0]*255), uint8(color[1]*255)))
		case 3:
			return l.SetPilot(NewPilotWithRGBW(100, uint8(color[0]*255), uint8(color[1]*255), uint8(color[2]*255), 0, 0))
		case 4:
			return l.SetPilot(NewPilotWithRGBW(100, uint8(color[0]*255), uint8(color[1]*255), uint8(color[2]*255), 0, uint8(color[3]*255)))
		case 5:
			return l.SetPilot(NewPilotWithRGBW(100, uint8(color[0]*255), uint8(color[1]*255), uint8(color[2]*255), uint8(color[3]*255), uint8(color[4]*255)))
		}
		return fmt.Errorf("unsupported channel amount of %d", color.Channels())

	default:
		return fmt.Errorf("supplied %d colors, this device has only 1", len(colors))
	}

}

// Colors returns the current colors of all the modules in the light device.
func (l *Light) Colors() ([]emission.CIE1931XYZColor, error) {
	pilot, err := l.GetPilot()
	if err != nil {
		return nil, fmt.Errorf("couldn't read pilot: %w", err)
	}

	product, err := l.Product()
	if err != nil {
		return nil, fmt.Errorf("failed to determine WiZ product: %w", err)
	}
	if product.moduleProfile == nil {
		return nil, fmt.Errorf("WiZ product doesn't contain module profile")
	}

	// TODO: Fix this crap
	dcsColor := make(emission.DCSColor, product.moduleProfile.Channels())
	for i := range dcsColor {
		switch i {
		case 0:
			if pilot.R != nil {
				dcsColor[i] = float64(*pilot.R) / 255
			}
		case 1:
			if pilot.G != nil {
				dcsColor[i] = float64(*pilot.G) / 255
			}
		case 2:
			if pilot.B != nil {
				dcsColor[i] = float64(*pilot.B) / 255
			}
		case 3:
			if pilot.CW != nil {
				dcsColor[i] = float64(*pilot.CW) / 255
			}
		case 4:
			if pilot.WW != nil {
				dcsColor[i] = float64(*pilot.WW) / 255
			}
		}
	}

	xyzColor, err := product.moduleProfile.DCSToXYZ(dcsColor)
	return []emission.CIE1931XYZColor{xyzColor}, err
}

// ModuleProfiles returns the profiles that describe each module.
func (l *Light) ModuleProfiles() ([]emission.ModuleProfile, error) {
	product, err := l.Product()
	if err != nil {
		return nil, fmt.Errorf("failed to determine WiZ product: %w", err)
	}
	if product.moduleProfile == nil {
		return nil, fmt.Errorf("WiZ product doesn't contain module profile")
	}

	return []emission.ModuleProfile{product.moduleProfile}, nil
}
