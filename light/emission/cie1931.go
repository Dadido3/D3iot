// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package emission

import (
	"fmt"
)

// CIE1931xyYAbs represents a color in the CIE 1931 XYZ color space with an absolute luminance in lumen.
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

// Relative returns the color in the CIE 1931 XYZ color space with relative luminance.
//
// maxLuminance defines the highest possible luminance in lumen.
func (c CIE1931XYZAbs) Relative(maxLuminance float64) CIE1931XYZRel {
	return CIE1931XYZRel{c.X / maxLuminance, c.Y / maxLuminance, c.Z / maxLuminance}
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

// CIE1931XYZRel represents a color in the CIE 1931 XYZ color space with a relative luminance.
//
// Luminance Y is in the range of [0, 1].
//
// The absolute brightness depends on the device/module that this emission value is rendered on.
// A luminance value Y of 1 corresponds with the full light output.
//
// An equal-energy radiator would result in X == Y == Z.
type CIE1931XYZRel struct {
	X, Y, Z float64
}

var _ Value = &CIE1931XYZRel{}

// IntoDCS implements the Value interface.
func (c CIE1931XYZRel) IntoDCS(mp ModuleProfile) DCSVector {
	maxLuminance := mp.WhitePoint().Y
	return c.Absolute(maxLuminance).IntoDCS(mp)
}

// FromDCS implements the Value interface.
func (c *CIE1931XYZRel) FromDCS(mp ModuleProfile, v DCSVector) error {
	var res CIE1931XYZAbs
	if err := res.FromDCS(mp, v); err != nil {
		return err
	}

	maxLuminance := mp.WhitePoint().Y
	*c = res.Relative(maxLuminance)
	return nil
}

// Absolute returns the color in the CIE 1931 XYZ color space with absolute luminance in lumen.
//
// maxLuminance defines the highest possible luminance in lumen.
func (c CIE1931XYZRel) Absolute(maxLuminance float64) CIE1931XYZAbs {
	return CIE1931XYZAbs{c.X * maxLuminance, c.Y * maxLuminance, c.Z * maxLuminance}
}

// Scaled returns c scaled by the scalar s.
func (c CIE1931XYZRel) Scaled(s float64) CIE1931XYZRel {
	return CIE1931XYZRel{c.X * s, c.Y * s, c.Z * s}
}

// CIE1931xyYAbs represents a color in the CIE 1931 xyY color space with an absolute luminance in lumen.
//
// An equal-energy radiator would result in x == y == 1/3.
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

// Relative returns the color in the CIE 1931 xyY color space with relative luminance.
//
// maxLuminance defines the highest possible luminance in lumen.
func (c CIE1931xyYAbs) Relative(maxLuminance float64) CIE1931xyYRel {
	return CIE1931xyYRel{c.X, c.Y, c.LuminanceY / maxLuminance}
}

// CIE1931XYZAbs returns the color in the CIE 1931 XYZ color space with absolute luminance in lumens.
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

// CIE1931xyYRel represents a color in the CIE 1931 xyY color space with a relative luminance.
//
// LuminanceY is in the range of [0, 1].
//
// The absolute brightness depends on the device/module that this emission value is rendered on.
// A luminance value Y of 1 corresponds with the full light output.
//
// An equal-energy radiator would result in x == y == 1/3.
type CIE1931xyYRel struct {
	X, Y       float64 // x, y in the range of [0, 1]
	LuminanceY float64 // Relative luminance Y in the range of [0, 1].
}

var _ Value = &CIE1931xyYRel{}

// IntoDCS implements the Value interface.
func (c CIE1931xyYRel) IntoDCS(mp ModuleProfile) DCSVector {
	maxLuminance := mp.WhitePoint().Y
	return c.Absolute(maxLuminance).IntoDCS(mp)
}

// FromDCS implements the Value interface.
func (c *CIE1931xyYRel) FromDCS(mp ModuleProfile, v DCSVector) error {
	var res CIE1931xyYAbs
	if err := res.FromDCS(mp, v); err != nil {
		return err
	}

	maxLuminance := mp.WhitePoint().Y
	*c = res.Relative(maxLuminance)
	return nil
}

// Absolute returns the color in the CIE 1931 xyY color space with absolute luminance in lumen.
//
// maxLuminance defines the highest possible luminance in lumen.
func (c CIE1931xyYRel) Absolute(maxLuminance float64) CIE1931xyYAbs {
	return CIE1931xyYAbs{c.X, c.Y, c.LuminanceY * maxLuminance}
}

// CIE1931XYZAbs returns the color in the CIE 1931 XYZ color space with relative luminance.
func (c CIE1931xyYRel) CIE1931XYZRel() CIE1931XYZRel {
	return CIE1931XYZRel{
		(c.X * c.LuminanceY) / c.Y,
		c.LuminanceY,
		(1 - c.X - c.Y) * c.LuminanceY / c.Y,
	}
}

// Scaled returns c scaled by the scalar s.
func (c CIE1931xyYRel) Scaled(s float64) CIE1931xyYRel {
	return CIE1931xyYRel{c.X, c.Y, c.LuminanceY * s}
}
