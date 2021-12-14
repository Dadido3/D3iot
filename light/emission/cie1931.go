// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package emission

import (
	"fmt"
)

// CIE1931XYZColor represents a color according to the CIE 1931 XYZ standard.
//
// The luminance Y is not relative, but absolute.
// And contrary to the normal definition, it uses lumen as its unit.
//
// An equal-energy radiator would result in X == Y == Z.
type CIE1931XYZColor struct {
	X, Y, Z float64
}

var _ Value = &CIE1931XYZColor{}

// IntoDCS implements the Value interface.
func (c CIE1931XYZColor) IntoDCS(mp ModuleProfile) DCSColor {
	return mp.XYZToDCS(c)
}

// FromDCS implements the Value interface.
func (c *CIE1931XYZColor) FromDCS(mp ModuleProfile, dcsColor DCSColor) error {
	var err error
	if *c, err = mp.DCSToXYZ(dcsColor); err != nil {
		return fmt.Errorf("failed to convert from DCS to %T: %w", c, err)
	}
	return nil
}

// Sum returns the sum of c and all colors.
func (c CIE1931XYZColor) Sum(colors ...CIE1931XYZColor) CIE1931XYZColor {
	result := c
	for _, color := range colors {
		result.X += color.X
		result.Y += color.Y
		result.Z += color.Z
	}
	return result
}

// Scaled returns c scaled by the scalar s.
func (c CIE1931XYZColor) Scaled(s float64) CIE1931XYZColor {
	return CIE1931XYZColor{c.X * s, c.Y * s, c.Z * s}
}

// CrossProd returns the cross product between two color vectors.
func (c CIE1931XYZColor) CrossProd(c2 CIE1931XYZColor) CIE1931XYZColor {
	return CIE1931XYZColor{c.Y*c2.Z - c.Z*c2.Y, c.Z*c2.X - c.X*c2.Z, c.X*c2.Y - c.Y*c2.X}
}

// CIE1931xyYColor represents a color in the CIE 1931 xyY color space.
//
// LuminanceY is not relative, but absolute.
// And contrary to the normal definition, it uses lumen as its unit.
type CIE1931xyYColor struct {
	X, Y       float64 // x, y in the range of [0, 1]
	LuminanceY float64 // Luminance Y in lumen.
}

var _ Value = &CIE1931xyYColor{}

// IntoDCS implements the Value interface.
func (c CIE1931xyYColor) IntoDCS(mp ModuleProfile) DCSColor {
	return mp.XYZToDCS(c.CIE1931XYZColor())
}

// FromDCS implements the Value interface.
func (c *CIE1931xyYColor) FromDCS(mp ModuleProfile, dcsColor DCSColor) error {
	if xyzColor, err := mp.DCSToXYZ(dcsColor); err != nil {
		return fmt.Errorf("failed to convert from DCS to %T: %w", c, err)
	} else {
		sum := xyzColor.X + xyzColor.Y + xyzColor.Z
		*c = CIE1931xyYColor{xyzColor.X / sum, xyzColor.Y / sum, xyzColor.Y}
	}
	return nil
}

func (c CIE1931xyYColor) CIE1931XYZColor() CIE1931XYZColor {
	return CIE1931XYZColor{
		(c.X * c.LuminanceY) / c.Y,
		c.LuminanceY,
		(1 - c.X - c.Y) * c.LuminanceY / c.Y,
	}
}

// Scaled returns c scaled by the scalar s.
func (c CIE1931xyYColor) Scaled(s float64) CIE1931xyYColor {
	return CIE1931xyYColor{c.X, c.Y, c.LuminanceY * s}
}
