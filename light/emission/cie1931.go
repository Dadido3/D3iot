// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package emission

import (
	"fmt"
)

// CIE1931XYZAbs represents a color according to the CIE 1931 XYZ standard.
//
// The luminance Y is not relative, but absolute.
// And contrary to the normal definition, it uses lumen as its unit.
//
// An equal-energy radiator would result in X == Y == Z.
type CIE1931XYZAbs struct {
	X, Y, Z float64
}

var _ Value = &CIE1931XYZAbs{}

// IntoDCS implements the Value interface.
func (c CIE1931XYZAbs) IntoDCS(mp ModuleProfile) DCSVector {
	return mp.XYZToDCS(c)
}

// FromDCS implements the Value interface.
func (c *CIE1931XYZAbs) FromDCS(mp ModuleProfile, v DCSVector) error {
	var err error
	if *c, err = mp.DCSToXYZ(v); err != nil {
		return fmt.Errorf("failed to convert from DCS to %T: %w", c, err)
	}
	return nil
}

// Sum returns the sum of c and all colors.
func (c CIE1931XYZAbs) Sum(colors ...CIE1931XYZAbs) CIE1931XYZAbs {
	result := c
	for _, color := range colors {
		result.X += color.X
		result.Y += color.Y
		result.Z += color.Z
	}
	return result
}

// Scaled returns c scaled by the scalar s.
func (c CIE1931XYZAbs) Scaled(s float64) CIE1931XYZAbs {
	return CIE1931XYZAbs{c.X * s, c.Y * s, c.Z * s}
}

// CrossProd returns the cross product between two color vectors.
func (c CIE1931XYZAbs) CrossProd(c2 CIE1931XYZAbs) CIE1931XYZAbs {
	return CIE1931XYZAbs{c.Y*c2.Z - c.Z*c2.Y, c.Z*c2.X - c.X*c2.Z, c.X*c2.Y - c.Y*c2.X}
}

// CIE1931xyYAbs represents a color in the CIE 1931 xyY color space.
//
// LuminanceY is not relative, but absolute.
// And contrary to the normal definition, it uses lumen as its unit.
type CIE1931xyYAbs struct {
	X, Y       float64 // x, y in the range of [0, 1]
	LuminanceY float64 // Luminance Y in lumen.
}

var _ Value = &CIE1931xyYAbs{}

// IntoDCS implements the Value interface.
func (c CIE1931xyYAbs) IntoDCS(mp ModuleProfile) DCSVector {
	return mp.XYZToDCS(c.CIE1931XYZAbs())
}

// FromDCS implements the Value interface.
func (c *CIE1931xyYAbs) FromDCS(mp ModuleProfile, v DCSVector) error {
	if xyzColor, err := mp.DCSToXYZ(v); err != nil {
		return fmt.Errorf("failed to convert from DCS to %T: %w", c, err)
	} else {
		sum := xyzColor.X + xyzColor.Y + xyzColor.Z
		*c = CIE1931xyYAbs{xyzColor.X / sum, xyzColor.Y / sum, xyzColor.Y}
	}
	return nil
}

func (c CIE1931xyYAbs) CIE1931XYZAbs() CIE1931XYZAbs {
	return CIE1931XYZAbs{
		(c.X * c.LuminanceY) / c.Y,
		c.LuminanceY,
		(1 - c.X - c.Y) * c.LuminanceY / c.Y,
	}
}

// Scaled returns c scaled by the scalar s.
func (c CIE1931xyYAbs) Scaled(s float64) CIE1931xyYAbs {
	return CIE1931xyYAbs{c.X, c.Y, c.LuminanceY * s}
}
