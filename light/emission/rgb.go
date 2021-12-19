// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package emission

// RGB represents any color space that contains RGB like primaries.
// Any type that implements this interface is able to convert between XYZ and its own RGB like color space.
type RGB interface {
	Value // Enforce that RGB color spaces implement the emission.Value interface.

	CIE1931XYZRel() CIE1931XYZRel             // CIE1931XYZRel transforms the RGB color space into a CIE 1931 XYZ color space with relative luminance.
	FromCIE1931XYZRel(xyzColor CIE1931XYZRel) // FromCIE1931XYZRel transforms the CIE 1931 XYZ color space with relative luminance into an RGB color space.
}
