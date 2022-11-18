// Copyright (c) 2021 David Vogel
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
			WhitePointColor: emission.CIE1931XYZAbs{X: 664.9241, Y: 752.0420, Z: 688.7433}.
				Sum(emission.CIE1931XYZAbs{X: 864.1765, Y: 768.9580, Z: 224.2333}),
			PrimaryColors: emission.TransformationLinDCSToXYZ{
				{X: 185.6709, Y: 83.9461, Z: 0.0320},
				{X: 46.3062, Y: 189.2580, Z: 24.4028},
				{X: 115.8685, Y: 67.5708, Z: 696.3857},
			},
			WhiteColors: emission.TransformationLinDCSToXYZ{
				{X: 664.9241, Y: 752.0420, Z: 688.7433},
				{X: 864.1765, Y: 768.9580, Z: 224.2333},
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
