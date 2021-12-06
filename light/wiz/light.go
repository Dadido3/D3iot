// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package wiz

import (
	"fmt"
	"image/color"
	"sync"
	"time"
)

type Light struct {
	address string

	product *Product // The cached product of the lamp. This will be queried from the device as soon as it is needed.

	deadline  time.Duration // Default deadline for a whole communcation action (sending and receiving).
	retries   uint          // Number of retries when the deadline got exceeded.
	connMutex sync.Mutex
}

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
	l.product = product
}

// SetColor sets the color of the light device.
func (l *Light) SetColor(c color.Color) error {
	// TODO: Convert colors
	return fmt.Errorf("not implemented yet")
}
