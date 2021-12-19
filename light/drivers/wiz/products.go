// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package wiz

import "github.com/Dadido3/D3iot/light/emission"

// products contains known WiZ light products.
// This can't be a map, as otherwise matching would get nondeterministic.
var products = []Product{
	{ // Tested: Yes, Profiled: Yes.
		moduleName:  "ESP03_SHRGB1W_01",
		deviceClass: deviceClassRGBTW,
		colorProfile: (&emission.ColorProfileGeneral{
			WhitePointColor: emission.CIE1931XYZAbs{X: 678.3330441265266, Y: 770.694363496923, Z: 628.5470196992967}.
				Sum(emission.CIE1931XYZAbs{X: 859.9603595688454, Y: 761.5085215026879, Z: 222.80239948974832}),
			PrimaryColors: emission.TransformationLinDCSToXYZ{
				{X: 206.8752389686704, Y: 93.13069968324798, Z: -0.2512365368135359},
				{X: 37.089570775907255, Y: 182.40345692892237, Z: 22.10488337648135},
				{X: 90.55843100384949, Y: 42.003319100660164, Z: 532.7732464538633},
			},
			WhiteColors: emission.TransformationLinDCSToXYZ{
				{X: 678.3330441265266, Y: 770.694363496923, Z: 628.5470196992967},
				{X: 859.9603595688454, Y: 761.5085215026879, Z: 222.80239948974832},
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
