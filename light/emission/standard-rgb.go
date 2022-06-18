// Copyright (c) 2021-2022 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package emission

// Reuse TransformationLinDCSToXYZ type for transformation from color space with 3 channels.
var standardRGBTransformation = TransformationLinDCSToXYZ{
	CIE1931XYZRel{0.4124, 0.2126, 0.0193}.Absolute(1),
	CIE1931XYZRel{0.3576, 0.7152, 0.1192}.Absolute(1),
	CIE1931XYZRel{0.1805, 0.0722, 0.9505}.Absolute(1),
}

// Reuse TransformationXYZToLinDCS type for transformation into color space with 3 channels.
var standardRGBInvTransformation = standardRGBTransformation.MustInverted()

// StandardRGB represents a color according to IEC 61966-2-1:1999.
// Commonly known as sRGB.
type StandardRGB struct {
	R, G, B float64 // R, G and B are in the range of [0, 1].
}

// Check if type implements the RGB interface.
var _ RGB = &StandardRGB{}

// CIE1931XYZRel transforms the RGB color space into a CIE 1931 XYZ color space with relative luminance.
// This implements the RGB interface.
func (c StandardRGB) CIE1931XYZRel() CIE1931XYZRel {
	// Reuse a DCSVector and the DCS transfer function.
	linRGB := DCSVector{c.R, c.G, c.B}.ClampedAndLinearized(TransferFunctionStandardRGB)

	// Reuse DCS transformation. Result is in absolute luminance, has to be converted into relative afterwards.
	transformed, _ := standardRGBTransformation.Multiplied(linRGB) // As the matrix exists and is 3x3, this shouldn't cause any error.

	return transformed.Relative(1)
}

// FromCIE1931XYZRel transforms the CIE 1931 XYZ color space with relative luminance into an RGB color space.
// This implements the RGB interface.
func (c *StandardRGB) FromCIE1931XYZRel(xyzColor CIE1931XYZRel) {
	// Reuse DCS transformation.
	transformed := standardRGBInvTransformation.Multiplied(xyzColor.Absolute(1))

	// Reuse DCS transfer function.
	linRGB := transformed.ClampedAndDeLinearized(TransferFunctionStandardRGB)
	c.R, c.G, c.B = linRGB[0], linRGB[1], linRGB[2] // As the inv matrix exists and is 3x3, these elements should exist.
}

// IntoDCS implements the Value interface.
func (c StandardRGB) IntoDCS(cp ColorProfile) DCSVector {
	return c.CIE1931XYZRel().IntoDCS(cp)
}

// FromDCS implements the ValueReceiver interface.
func (c *StandardRGB) FromDCS(cp ColorProfile, v DCSVector) error {
	var xyzColor CIE1931XYZRel
	if err := xyzColor.FromDCS(cp, v); err != nil {
		return err
	}

	c.FromCIE1931XYZRel(xyzColor)
	return nil
}
