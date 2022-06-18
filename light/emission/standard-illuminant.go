// Copyright (c) 2021-2022 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package emission

var (
	// StandardIlluminantA represents a CIE standard illuminant.
	//	col := StandardIlluminantA.Absolute(789) // Returns the illuminant with an absolute luminance of 789 lumens.
	//	col := StandardIlluminantA.Scaled(0.5)   // Returns the illuminant at half the max. brightness of the device/module.
	StandardIlluminantA = CIE1931xyYRel{0.44757, 0.40745, 1}.CIE1931XYZRel()

	// StandardIlluminantB represents a CIE standard illuminant.
	//	col := StandardIlluminantB.Absolute(789) // Returns the illuminant with an absolute luminance of 789 lumens.
	//	col := StandardIlluminantB.Scaled(0.5)   // Returns the illuminant at half the max. brightness of the device/module.
	StandardIlluminantB = CIE1931xyYRel{0.34842, 0.35161, 1}.CIE1931XYZRel()

	// StandardIlluminantC represents a CIE standard illuminant.
	//	col := StandardIlluminantC.Absolute(789) // Returns the illuminant with an absolute luminance of 789 lumens.
	//	col := StandardIlluminantC.Scaled(0.5)   // Returns the illuminant at half the max. brightness of the device/module.
	StandardIlluminantC = CIE1931xyYRel{0.31006, 0.31616, 1}.CIE1931XYZRel()

	// StandardIlluminantD50 represents a CIE standard illuminant.
	//	col := StandardIlluminantD50.Absolute(789) // Returns the illuminant with an absolute luminance of 789 lumens.
	//	col := StandardIlluminantD50.Scaled(0.5)   // Returns the illuminant at half the max. brightness of the device/module.
	StandardIlluminantD50 = CIE1931xyYRel{0.34567, 0.35850, 1}.CIE1931XYZRel()

	// StandardIlluminantD55 represents a CIE standard illuminant.
	//	col := StandardIlluminantD55.Absolute(789) // Returns the illuminant with an absolute luminance of 789 lumens.
	//	col := StandardIlluminantD55.Scaled(0.5)   // Returns the illuminant at half the max. brightness of the device/module.
	StandardIlluminantD55 = CIE1931xyYRel{0.33242, 0.35850, 1}.CIE1931XYZRel()

	// StandardIlluminantD65 represents a CIE standard illuminant.
	//	col := StandardIlluminantD65.Absolute(789) // Returns the illuminant with an absolute luminance of 789 lumens.
	//	col := StandardIlluminantD65.Scaled(0.5)   // Returns the illuminant at half the max. brightness of the device/module.
	StandardIlluminantD65 = CIE1931xyYRel{0.31271, 0.32902, 1}.CIE1931XYZRel()

	// StandardIlluminantD75 represents a CIE standard illuminant.
	//	col := StandardIlluminantD75.Absolute(789) // Returns the illuminant with an absolute luminance of 789 lumens.
	//	col := StandardIlluminantD75.Scaled(0.5)   // Returns the illuminant at half the max. brightness of the device/module.
	StandardIlluminantD75 = CIE1931xyYRel{0.29902, 0.31485, 1}.CIE1931XYZRel()

	// StandardIlluminantD93 represents a CIE standard illuminant.
	//	col := StandardIlluminantD93.Absolute(789) // Returns the illuminant with an absolute luminance of 789 lumens.
	//	col := StandardIlluminantD93.Scaled(0.5)   // Returns the illuminant at half the max. brightness of the device/module.
	StandardIlluminantD93 = CIE1931xyYRel{0.28315, 0.29711, 1}.CIE1931XYZRel()

	// StandardIlluminantE represents a CIE standard illuminant.
	//	col := StandardIlluminantE.Absolute(789) // Returns the illuminant with an absolute luminance of 789 lumens.
	//	col := StandardIlluminantE.Scaled(0.5)   // Returns the illuminant at half the max. brightness of the device/module.
	StandardIlluminantE = CIE1931XYZRel{1, 1, 1}
)

// StandardIlluminantDSeries represents a daylight color according to the CIE standard illuminant series D.
//
// The valid temperature range is 4000K to 25000K.
type StandardIlluminantDSeries struct {
	Temperature float64 // Temperature in K.
	Luminance   float64 // Luminance in lumen.
}

var _ Value = &StandardIlluminantDSeries{} // TODO: Implement transformation from DCS

// IntoDCS implements the Value interface.
func (si StandardIlluminantDSeries) IntoDCS(cp ColorProfile) DCSVector {
	return cp.XYZToDCS(si.CIE1931XYZAbs())
}

// FromDCS implements the ValueReceiver interface.
/*func (si *StandardIlluminantDSeries) FromDCS(cp ColorProfile, v DCSVector) error {
	// Calculate CCT.
	return fmt.Errorf("conversion from DCS to %T not implemented yet", si)
}*/

func (si StandardIlluminantDSeries) CIE1931xyYAbs() CIE1931xyYAbs {
	// Source: https://en.wikipedia.org/wiki/Standard_illuminant.

	t := si.Temperature

	var x float64
	switch {
	case t >= 4000 && t <= 7000:
		x = 0.244063 + 0.09911e3/t + 2.9678e6/t/t - 4.6070e9/t/t/t
	case t >= 7000 && t <= 25000:
		x = 0.237040 + 0.24748e3/t + 1.9018e6/t/t - 2.0064e9/t/t/t
	default:
		return CIE1931xyYAbs{0.3, 0.3, 0}
	}

	y := -3.000*x*x + 2.870*x - 0.275

	return CIE1931xyYAbs{x, y, si.Luminance}
}

func (si StandardIlluminantDSeries) CIE1931XYZAbs() CIE1931XYZAbs {
	return si.CIE1931xyYAbs().CIE1931XYZAbs()
}
