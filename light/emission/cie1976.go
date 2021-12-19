// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package emission

import "math"

// CIE1976LAB represents a color in the L*a*b* color space defined by the CIE in 1976.
type CIE1976LAB struct {
	L float64 // Perceptual lightness L* in the range of [0, 100].
	A float64 // Redness a*.
	B float64 // Blueness b*.

	WhitePoint CIE1931XYZRel // White point in CIE 1931 XYZ coordinates with relative luminance.
}

// TODO: Implement emission.Value interface.

// CIE1931XYZRel returns the color in the CIE 1931 XYZ color space with relative luminance.
func (c CIE1976LAB) CIE1931XYZRel() CIE1931XYZRel {
	// Using the intent of the CIE standard, not the numbers published by the CIE. See http://www.brucelindbloom.com/LContinuity.html.
	const delta = 6.0 / 29

	fInv := func(t float64) float64 {
		if t > delta {
			return t * t * t
		} else {
			return 3 * delta * delta * (t - 4.0/29)
		}
	}

	lScaled := (c.L + 16) / 116

	return CIE1931XYZRel{
		X: c.WhitePoint.X * fInv(lScaled+c.A/500),
		Y: c.WhitePoint.Y * fInv(lScaled),
		Z: c.WhitePoint.Z * fInv(lScaled-c.B/200),
	}
}

// Distance returns the euclidean distance between c and c2.
// This is a good measure for perceptual difference of two colors.
//
// The distance is officially called ΔE*ab (with ab being in subscript) or just ΔE*.
// A value of about 2.3 is just noticeable.
//
// Both L*a*b* colors need to be in a L*a*b* space with the same white point.
// If they don't, the result will be meaningless.
func (c CIE1976LAB) Distance(c2 CIE1976LAB) float64 {
	lDiff := c.L - c2.L
	aDiff := c.A - c2.A
	bDiff := c.B - c2.B

	return math.Sqrt(lDiff*lDiff + aDiff*aDiff + bDiff*bDiff)
}

// DistanceSqr returns the squared euclidean distance between c and c2.
//
// See c.Distance() for details.
func (c CIE1976LAB) DistanceSqr(c2 CIE1976LAB) float64 {
	lDiff := c.L - c2.L
	aDiff := c.A - c2.A
	bDiff := c.B - c2.B

	return lDiff*lDiff + aDiff*aDiff + bDiff*bDiff
}
