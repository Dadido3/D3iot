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

// SetColor sets the color of the light device.
func (l *Light) SetColor(c color.Color) error {
	// TODO: Convert colors
	return fmt.Errorf("not implemented yet")
}
