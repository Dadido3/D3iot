package light

import "image/color"

type Light interface {
	// SetColor sets the color of the light device.
	SetColor(c color.Color) error
}
