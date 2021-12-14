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
		moduleProfile: (&emission.ModuleProfileGeneral{
			WhitePointColor: emission.CIE1931XYZColor{X: 0.4453145432034966, Y: 0.4936881395023207, Z: 0.42839706385865284}.
				Sum(emission.CIE1931XYZColor{X: 0.5750302833988727, Y: 0.5063118604976793, Z: 0.16587493725962105}).Scaled(1521), // Scale by lumen.
			PrimaryColors: emission.TransformationLinDCSToXYZ{
				{X: 0.13404038256313233, Y: 0.05773272383434033, Z: 0.0},
				{X: 0.027205814312668924, Y: 0.1303357517167337, Z: 0.006074024335587602},
				{X: 0.07715380598919405, Y: 0.04729943996432647, Z: 0.3951593127343209},
			}.Scaled(1521), // Scale by lumen.
			WhiteColors: emission.TransformationLinDCSToXYZ{
				{X: 0.4453145432034966, Y: 0.4936881395023207, Z: 0.42839706385865284},
				{X: 0.5750302833988727, Y: 0.5063118604976793, Z: 0.16587493725962105},
			}.Scaled(1521), // Scale by lumen.
			OutputLimiter: emission.OutputLimiterSum{Limit: 2},
		}).MustInit(),
		minTemp: newPtrUInt(2200), maxTemp: newPtrUInt(6500),
	},

	{ // Tested: No, Profiled: No.
		moduleName:  "ESP01_SHDW_01",
		deviceClass: deviceClassDW,
		moduleProfile: (&emission.ModuleProfileGeneral{
			WhitePointColor: emission.CIE1931XYZColor{X: 0.5750302833988727, Y: 0.5063118604976793, Z: 0.16587493725962105}.Scaled(810),             // Scale by lumen.
			WhiteColors:     emission.TransformationLinDCSToXYZ{{X: 0.5750302833988727, Y: 0.5063118604976793, Z: 0.16587493725962105}}.Scaled(810), // Scale by lumen.
			OutputLimiter:   emission.OutputLimiterSum{Limit: 1},
		}).MustInit(),
	},

	{ // Tested: No, Profiled: No.
		moduleName:  "ESP56_SHTW3_01",
		deviceClass: deviceClassTW,
		moduleProfile: (&emission.ModuleProfileGeneral{
			WhitePointColor: emission.CIE1931XYZColor{X: 0.4453145432034966, Y: 0.4936881395023207, Z: 0.42839706385865284}.
				Sum(emission.CIE1931XYZColor{X: 0.5750302833988727, Y: 0.5063118604976793, Z: 0.16587493725962105}).Scaled(720), // Scale by lumen.
			WhiteColors: emission.TransformationLinDCSToXYZ{
				{X: 0.4453145432034966, Y: 0.4936881395023207, Z: 0.42839706385865284},
				{X: 0.5750302833988727, Y: 0.5063118604976793, Z: 0.16587493725962105},
			}.Scaled(720), // Scale by lumen.
			OutputLimiter: emission.OutputLimiterSum{Limit: 2},
		}).MustInit(),
		minTemp: newPtrUInt(2200), maxTemp: newPtrUInt(5500),
	},
}
