// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

// CIE1931XYZAbs represents a color according to the CIE 1931 XYZ standard.
// The component Y is in the range [0, 1].
// The X and Z components are somewhat in the range [0, 1].
type CIE1931XYZAbs struct {
	X, Y, Z float64
}

// Add returns the sum of c and c2.
func (c CIE1931XYZAbs) Add(c2 CIE1931XYZAbs) CIE1931XYZAbs {
	return CIE1931XYZAbs{c.X + c2.X, c.Y + c2.Y, c.Z + c2.Z}
}

// Scale returns c scaled by the scalar s
func (c CIE1931XYZAbs) Scale(s float64) CIE1931XYZAbs {
	return CIE1931XYZAbs{c.X * s, c.Y * s, c.Z * s}
}

// RGBWValue represents a color in the color space of a specific WiZ light device.
// Together with a set of primaries (RGBWPrimaries), this can be transformed into another color space.
//
// This space is linear, with one exception:
// The sum of all values can't exceed a maximum value (E.g. 512).
// If it does, the values will be normalized in the device so that the sum fits within the allowed range.
//
// We will assume that the space is fully linear for profiling.
type RGBWValue struct {
	R, G, B, CW, WW uint8
}

// CIE1931XYZAbs combines the values with the given primaries, and returns a color in the CIE 1931 XYZ space.
func (v RGBWValue) CIE1931XYZAbs(primaries RGBWPrimaries) CIE1931XYZAbs {
	c := CIE1931XYZAbs{}

	r := primaries.R.Scale(float64(v.R) / 255)
	g := primaries.G.Scale(float64(v.G) / 255)
	b := primaries.B.Scale(float64(v.B) / 255)
	cw := primaries.CW.Scale(float64(v.CW) / 255)
	ww := primaries.WW.Scale(float64(v.WW) / 255)

	return c.Add(r).Add(g).Add(b).Add(cw).Add(ww)
}
