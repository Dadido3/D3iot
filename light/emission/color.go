// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package emission

// CIE1931XYZColor represents a color according to the CIE 1931 XYZ standard.
//
// The luminance Y is not relative, but absolute.
// And contrary to the normal definition, it uses lumen as its unit.
//
// An equal-energy radiator would result in X == Y == Z.
type CIE1931XYZColor struct {
	X, Y, Z float64
}

// Add returns the sum of c and c2.
func (c CIE1931XYZColor) Add(c2 CIE1931XYZColor) CIE1931XYZColor {
	return CIE1931XYZColor{c.X + c2.X, c.Y + c2.Y, c.Z + c2.Z}
}

// Scale returns c scaled by the scalar s.
func (c CIE1931XYZColor) Scale(s float64) CIE1931XYZColor {
	return CIE1931XYZColor{c.X * s, c.Y * s, c.Z * s}
}

// CrossProd returns the cross product between two color vectors.
func (c CIE1931XYZColor) CrossProd(c2 CIE1931XYZColor) CIE1931XYZColor {
	return CIE1931XYZColor{c.Y*c2.Z - c.Z*c2.Y, c.Z*c2.X - c.X*c2.Z, c.X*c2.Y - c.Y*c2.X}
}
