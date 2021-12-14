// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package emission

var (
	// StandardIlluminantA represents a CIE standard illuminant.
	// Use the Scaled(s) method to change the luminosity:
	//	col := StandardIlluminantA.Scaled(789) // Returns the illuminant with a luminosity of 789 lumens.
	StandardIlluminantA = CIE1931xyYColor{0.44757, 0.40745, 1}

	// StandardIlluminantB represents a CIE standard illuminant.
	// Use the Scaled(s) method to change the luminosity:
	//	col := StandardIlluminantB.Scaled(789) // Returns the illuminant with a luminosity of 789 lumens.
	StandardIlluminantB = CIE1931xyYColor{0.34842, 0.35161, 1}

	// StandardIlluminantC represents a CIE standard illuminant.
	// Use the Scaled(s) method to change the luminosity:
	//	col := StandardIlluminantC.Scaled(789) // Returns the illuminant with a luminosity of 789 lumens.
	StandardIlluminantC = CIE1931xyYColor{0.31006, 0.31616, 1}

	// StandardIlluminantD50 represents a CIE standard illuminant.
	// Use the Scaled(s) method to change the luminosity:
	//	col := StandardIlluminantD50.Scaled(789) // Returns the illuminant with a luminosity of 789 lumens.
	StandardIlluminantD50 = CIE1931xyYColor{0.34567, 0.35850, 1}

	// StandardIlluminantD55 represents a CIE standard illuminant.
	// Use the Scaled(s) method to change the luminosity:
	//	col := StandardIlluminantD55.Scaled(789) // Returns the illuminant with a luminosity of 789 lumens.
	StandardIlluminantD55 = CIE1931xyYColor{0.33242, 0.35850, 1}

	// StandardIlluminantD65 represents a CIE standard illuminant.
	// Use the Scaled(s) method to change the luminosity:
	//	col := StandardIlluminantD65.Scaled(789) // Returns the illuminant with a luminosity of 789 lumens.
	StandardIlluminantD65 = CIE1931xyYColor{0.31271, 0.32902, 1}

	// StandardIlluminantD75 represents a CIE standard illuminant.
	// Use the Scaled(s) method to change the luminosity:
	//	col := StandardIlluminantD75.Scaled(789) // Returns the illuminant with a luminosity of 789 lumens.
	StandardIlluminantD75 = CIE1931xyYColor{0.29902, 0.31485, 1}

	// StandardIlluminantD93 represents a CIE standard illuminant.
	// Use the Scaled(s) method to change the luminosity:
	//	col := StandardIlluminantD93.Scaled(789) // Returns the illuminant with a luminosity of 789 lumens.
	StandardIlluminantD93 = CIE1931xyYColor{0.28315, 0.29711, 1}

	// StandardIlluminantE represents a CIE standard illuminant.
	// Use the Scaled(s) method to change the luminosity:
	//	col := StandardIlluminantE.Scaled(789) // Returns the illuminant with a luminosity of 789 lumens.
	StandardIlluminantE = CIE1931XYZColor{1, 1, 1}
)

// StandardIlluminantDSeries represents a daylight color according to the CIE standard illuminant series D.
//
// The valid temperature range is 4000K to 25000K.
type StandardIlluminantDSeries struct {
	Temperature float64 // Temperature in K.
	Luminance   float64 // Luminance in lumen.
}

var _ ValueIntoDCS = &StandardIlluminantDSeries{} // TODO: Implement transformation from DCS

// IntoDCS implements the Value interface.
func (si StandardIlluminantDSeries) IntoDCS(mp ModuleProfile) DCSColor {
	return mp.XYZToDCS(si.CIE1931XYZColor())
}

// FromDCS implements the Value interface.
/*func (si *StandardIlluminantDSeries) FromDCS(mp ModuleProfile, c DCSColor) error {
	// Calculate CCT.
	return fmt.Errorf("conversion from DCS to %T not implemented yet", si)
}*/

func (si StandardIlluminantDSeries) CIE1931xyYColor() CIE1931xyYColor {
	// Source: https://en.wikipedia.org/wiki/Standard_illuminant.

	t := si.Temperature

	var x float64
	switch {
	case t >= 4000 && t <= 7000:
		x = 0.244063 + 0.09911e3/t + 2.9678e6/t/t - 4.6070e9/t/t/t
	case t >= 7000 && t <= 25000:
		x = 0.237040 + 0.24748e3/t + 1.9018e6/t/t - 2.0064e9/t/t/t
	default:
		return CIE1931xyYColor{0.3, 0.3, 0}
	}

	y := -3.000*x*x + 2.870*x - 0.275

	return CIE1931xyYColor{x, y, si.Luminance}
}

func (si StandardIlluminantDSeries) CIE1931XYZColor() CIE1931XYZColor {
	return si.CIE1931xyYColor().CIE1931XYZColor()
}
