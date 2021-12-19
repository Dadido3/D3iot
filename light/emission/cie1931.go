// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package emission

import (
	"fmt"
	"math"
)

// CIE1931xyYAbs represents a color in the CIE 1931 XYZ color space with an absolute luminance in lumen.
//
// An equal-energy radiator would result in X == Y == Z.
type CIE1931XYZAbs struct {
	X, Y, Z float64
}

var _ Value = &CIE1931XYZAbs{}

// IntoDCS implements the Value interface.
func (c CIE1931XYZAbs) IntoDCS(cp ColorProfile) DCSVector {
	return cp.XYZToDCS(c)
}

// FromDCS implements the Value interface.
func (c *CIE1931XYZAbs) FromDCS(cp ColorProfile, v DCSVector) error {
	var err error
	if *c, err = cp.DCSToXYZ(v); err != nil {
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

// CIE1931xyYAbs returns the color in the CIE 1931 xyY color space with an absolute luminance in lumen.
func (c CIE1931XYZAbs) CIE1931xyYAbs() CIE1931xyYAbs {
	sum := c.X + c.Y + c.Z

	return CIE1931xyYAbs{
		X:          c.X / sum,
		Y:          c.Y / sum,
		LuminanceY: c.Y,
	}
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
func (c CIE1931XYZRel) IntoDCS(cp ColorProfile) DCSVector {
	maxLuminance := cp.WhitePoint().Y
	return c.Absolute(maxLuminance).IntoDCS(cp)
}

// FromDCS implements the Value interface.
func (c *CIE1931XYZRel) FromDCS(cp ColorProfile, v DCSVector) error {
	var res CIE1931XYZAbs
	if err := res.FromDCS(cp, v); err != nil {
		return err
	}

	maxLuminance := cp.WhitePoint().Y
	*c = res.Relative(maxLuminance)
	return nil
}

// TransformRGB writes the color into the given RGB color space.
//
//	var rgbColor StandardRGB
//	c.TransformRGB(&rgbColor) // Writes result into rgbColor.
func (c CIE1931XYZRel) TransformRGB(rgbColor RGB) {
	rgbColor.FromCIE1931XYZRel(c)
}

// Absolute returns the color in the CIE 1931 XYZ color space with absolute luminance in lumen.
//
// maxLuminance defines the highest possible luminance in lumen.
func (c CIE1931XYZRel) Absolute(maxLuminance float64) CIE1931XYZAbs {
	return CIE1931XYZAbs{c.X * maxLuminance, c.Y * maxLuminance, c.Z * maxLuminance}
}

// CIE1931xyYRel returns the color in the CIE 1931 xyY color space with a relative luminance.
func (c CIE1931XYZRel) CIE1931xyYRel() CIE1931xyYRel {
	sum := c.X + c.Y + c.Z

	return CIE1931xyYRel{
		X:          c.X / sum,
		Y:          c.Y / sum,
		LuminanceY: c.Y,
	}
}

// Scaled returns c scaled by the scalar s.
func (c CIE1931XYZRel) Scaled(s float64) CIE1931XYZRel {
	return CIE1931XYZRel{c.X * s, c.Y * s, c.Z * s}
}

// Distance returns the euclidean distance between c and c2.
//
// This distance doesn't represent perceptual difference of two colors.
// See the method CIE1976LABDistance() for a better metric.
func (c CIE1931XYZRel) Distance(c2 CIE1931XYZRel) float64 {
	xDiff := c.X - c2.X
	yDiff := c.Y - c2.Y
	zDiff := c.Z - c2.Z

	return math.Sqrt(xDiff*xDiff + yDiff*yDiff + zDiff*zDiff)
}

// CIE1976LABDistance returns the euclidean distance between c and c2 in the CIE 1976 L*a*b* color space.
// This is a good measure for perceptual difference of two colors.
//
// The distance is officially called ΔE*ab (with ab being in subscript) or just ΔE*.
// A value of about 2.3 is just noticeable.
func (c CIE1931XYZRel) CIE1976LABDistance(c2 CIE1931XYZRel, whitePoint CIE1931XYZRel) float64 {
	cLAB := c.CIE1976LAB(whitePoint)
	c2LAB := c2.CIE1976LAB(whitePoint)

	return cLAB.Distance(c2LAB)
}

// CIE1976LABDistanceSqr returns the squared euclidean distance between c and c2 in the CIE 1976 L*a*b* color space.
//
// See c.CIE1976LABDistance() for details.
func (c CIE1931XYZRel) CIE1976LABDistanceSqr(c2 CIE1931XYZRel, whitePoint CIE1931XYZRel) float64 {
	cLAB := c.CIE1976LAB(whitePoint)
	c2LAB := c2.CIE1976LAB(whitePoint)

	return cLAB.DistanceSqr(c2LAB)
}

// CIE1976LAB returns the color transformed into the CIE 1976 L*a*b* color space with the given white point.
func (c CIE1931XYZRel) CIE1976LAB(whitePoint CIE1931XYZRel) CIE1976LAB {
	// Using the intent of the CIE standard, not the numbers published by the CIE. See http://www.brucelindbloom.com/LContinuity.html.
	const delta = 6.0 / 29

	f := func(t float64) float64 {
		if t > delta*delta*delta {
			return math.Pow(t, 1.0/3)
		} else {
			return t/(3*delta*delta) + 4.0/29
		}
	}

	return CIE1976LAB{
		L:          116*f(c.Y/whitePoint.Y) - 16,
		A:          500 * (f(c.X/whitePoint.X) - f(c.Y/whitePoint.Y)),
		B:          200 * (f(c.Y/whitePoint.Y) - f(c.Z/whitePoint.Z)),
		WhitePoint: whitePoint,
	}
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
func (c CIE1931xyYAbs) IntoDCS(cp ColorProfile) DCSVector {
	return cp.XYZToDCS(c.CIE1931XYZAbs())
}

// FromDCS implements the Value interface.
func (c *CIE1931xyYAbs) FromDCS(cp ColorProfile, v DCSVector) error {
	if xyzColor, err := cp.DCSToXYZ(v); err != nil {
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
func (c CIE1931xyYRel) IntoDCS(cp ColorProfile) DCSVector {
	maxLuminance := cp.WhitePoint().Y
	return c.Absolute(maxLuminance).IntoDCS(cp)
}

// FromDCS implements the Value interface.
func (c *CIE1931xyYRel) FromDCS(cp ColorProfile, v DCSVector) error {
	var res CIE1931xyYAbs
	if err := res.FromDCS(cp, v); err != nil {
		return err
	}

	maxLuminance := cp.WhitePoint().Y
	*c = res.Relative(maxLuminance)
	return nil
}

// Absolute returns the color in the CIE 1931 xyY color space with absolute luminance in lumen.
//
// maxLuminance defines the highest possible luminance in lumen.
func (c CIE1931xyYRel) Absolute(maxLuminance float64) CIE1931xyYAbs {
	return CIE1931xyYAbs{c.X, c.Y, c.LuminanceY * maxLuminance}
}

// CIE1931XYZRel returns the color in the CIE 1931 XYZ color space with relative luminance.
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
