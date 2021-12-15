// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package emission

// BlackBodyFixed represents the emission of black body radiator with a fixed luminance in lumen.
//
// This doesn't return a color of a daylight temperature.
//
// The valid temperature range is 1667K to 25000K.
type BlackBodyFixed struct {
	Temperature float64 // Temperature in K.
	Luminance   float64 // Luminance in lumen.
}

var _ ValueIntoDCS = &BlackBodyFixed{} // TODO: Implement transformation from DCS

// IntoDCS implements the Value interface.
func (b BlackBodyFixed) IntoDCS(mp ModuleProfile) DCSVector {
	return mp.XYZToDCS(b.CIE1931XYZAbs())
}

// FromDCS implements the Value interface.
/*func (b *BlackBodyFixed) FromDCS(mp ModuleProfile, v DCSVector) error {
	// Calculate CCT.
	return fmt.Errorf("conversion from DCS to %T not implemented yet", b)
}*/

func (b BlackBodyFixed) CIE1931xyYAbs() CIE1931xyYAbs {
	// Source: https://en.wikipedia.org/wiki/Planckian_locus.

	t := b.Temperature

	var x float64
	switch {
	case t >= 1667 && t <= 4000:
		x = -0.2661239e9/t/t/t - 0.2343589e6/t/t + 0.8776956e3/t + 0.179910
	case t >= 4000 && t <= 25000:
		x = -3.0258469e9/t/t/t + 2.1070379e6/t/t + 0.2226347e3/t + 0.240390
	default:
		return CIE1931xyYAbs{0.3, 0.3, 0}
	}

	var y float64
	switch {
	case t >= 1667 && t <= 2222:
		y = -1.1063814*x*x*x - 1.34811020*x*x + 2.18555832*x - 0.20219683
	case t >= 2222 && t <= 4000:
		y = -0.9549476*x*x*x - 1.37418593*x*x + 2.09137015*x - 0.16748867
	case t >= 4000 && t <= 25000:
		y = +3.0817580*x*x*x - 5.87338670*x*x + 3.75112997*x - 0.37001483
	default:
		return CIE1931xyYAbs{0.3, 0.3, 0}
	}

	return CIE1931xyYAbs{x, y, b.Luminance}
}

func (b BlackBodyFixed) CIE1931XYZAbs() CIE1931XYZAbs {
	return b.CIE1931xyYAbs().CIE1931XYZAbs()
}
