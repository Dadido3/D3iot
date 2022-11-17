// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package wiz

import "github.com/Dadido3/D3iot/light/emission"

// products contains known WiZ light products.
// This can't be a map, as otherwise matching would get nondeterministic.
var products = []Product{
	{ // Tested: Yes, Profiled: Yes (GretagMacbeth eye-one Pro 42.17.79).
		moduleName:  "ESP03_SHRGB1W_01",
		deviceClass: deviceClassRGBTW,
		colorProfile: (&emission.ColorProfileGeneral{
			WhitePointColor: emission.CIE1931XYZAbs{X: 680.061978810243, Y: 761.238333698754, Z: 746.14571679305}.
				Sum(emission.CIE1931XYZAbs{X: 844.960958613623, Y: 759.761666301247, Z: 230.895041995431}),
			PrimaryColors: emission.TransformationLinDCSToXYZ{
				{X: 198.136406299187, Y: 86.6027529080741, Z: 0.0479301103007406},
				{X: 43.7905927586912, Y: 195.166857944595, Z: 31.6765054168812},
				{X: 129.149462685368, Y: 51.1132970868513, Z: 731.546582530092},
			},
			WhiteColors: emission.TransformationLinDCSToXYZ{
				{X: 680.061978810243, Y: 761.238333698754, Z: 746.14571679305},
				{X: 844.960958613623, Y: 759.761666301247, Z: 230.895041995431},
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
