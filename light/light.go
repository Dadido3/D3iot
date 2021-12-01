// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package light

import "image/color"

type Light interface {
	// SetColor sets the color of the light device.
	SetColor(c color.Color) error
}
