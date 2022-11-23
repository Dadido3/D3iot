// Copyright (c) 2021-2022 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package wiz

import "github.com/Dadido3/D3iot/light/emission"

// products contains known WiZ light products.
// This can't be a map, as otherwise matching would get nondeterministic.
var products = []Product{
	{ // Tested: Yes, Profiled: Yes (GretagMacbeth eye-one Pro 42.17.79, CIE 2012 2°).
		moduleName:  "ESP03_SHRGB1W_01",
		deviceClass: deviceClassRGBTW,
		colorProfile: (&emission.ColorProfileGeneral{
			WhitePointColor: emission.CIE1931XYZAbs{X: 669.8132, Y: 750.2355, Z: 678.1854}.
				Sum(emission.CIE1931XYZAbs{X: 851.3526, Y: 752.0830, Z: 218.1395}),
			PrimaryColors: emission.TransformationLinDCSToXYZ{
				{X: 157.1970, Y: 70.2669, Z: -0.0390},
				{X: 48.3486, Y: 191.7574, Z: 25.0336},
				{X: 110.3764, Y: 65.9633, Z: 666.7974},
			},
			WhiteColors: emission.TransformationLinDCSToXYZ{
				{X: 669.8132, Y: 750.2355, Z: 678.1854},
				{X: 851.3526, Y: 752.0830, Z: 218.1395},
			},
			OutputLimiter: emission.OutputLimiterSum{Limit: 2},
		}).MustInit(),
		minTemp: newPtrUInt(2200), maxTemp: newPtrUInt(6500),
	},

	{ // Tested: No, Profiled: No.
		moduleName:  "ESP01_SHDW_01",
		deviceClass: deviceClassDW,
		colorProfile: (&emission.ColorProfileGeneral{
			WhitePointColor: emission.CIE1931XYZAbs{X: 0.5750302833988727, Y: 0.5063118604976793, Z: 0.16587493725962105}.Scaled(810),               // Scale by lumen.
			WhiteColors:     emission.TransformationLinDCSToXYZ{{X: 0.5750302833988727, Y: 0.5063118604976793, Z: 0.16587493725962105}}.Scaled(810), // Scale by lumen.
			OutputLimiter:   emission.OutputLimiterSum{Limit: 1},
		}).MustInit(),
	},

	{ // Tested: No, Profiled: No.
		moduleName:  "ESP56_SHTW3_01",
		deviceClass: deviceClassTW,
		colorProfile: (&emission.ColorProfileGeneral{
			WhitePointColor: emission.CIE1931XYZAbs{X: 0.4453145432034966, Y: 0.4936881395023207, Z: 0.42839706385865284}.
				Sum(emission.CIE1931XYZAbs{X: 0.5750302833988727, Y: 0.5063118604976793, Z: 0.16587493725962105}).Scaled(720), // Scale by lumen.
			WhiteColors: emission.TransformationLinDCSToXYZ{
				{X: 0.4453145432034966, Y: 0.4936881395023207, Z: 0.42839706385865284},
				{X: 0.5750302833988727, Y: 0.5063118604976793, Z: 0.16587493725962105},
			}.Scaled(720), // Scale by lumen.
			OutputLimiter: emission.OutputLimiterSum{Limit: 2},
		}).MustInit(),
		minTemp: newPtrUInt(2200), maxTemp: newPtrUInt(5500),
	},
}
