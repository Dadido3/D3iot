// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package emission

import (
	"math"
	"testing"
)

// Scaled sRGB primaries with a whitepoint of D65.
var (
	standardRGBRed   = CIE1931XYZAbs{0.4124564, 0.2126729, 0.0193339}
	standardRGBGreen = CIE1931XYZAbs{0.3575761, 0.7151522, 0.1191920}
	standardRGBBlue  = CIE1931XYZAbs{0.1804375, 0.0721750, 0.9503041}
)

func TestDCSToXYZ1(t *testing.T) {
	// Color profile with sRGB primaries, but linear transfer function.
	colorProfile := &ColorProfileGeneral{
		PrimaryColors: TransformationLinDCSToXYZ{standardRGBRed, standardRGBGreen, standardRGBBlue},
	}
	colorProfile.MustInit()

	color, err := colorProfile.DCSToXYZ([]float64{1, 0, 0})
	if err != nil {
		t.Fatalf("DCSToXYZ() failed: %v", err)
	}

	if color != standardRGBRed {
		t.Errorf("DCSToXYZ() returned wrong color result. Got %v, want %v", color, standardRGBRed)
	}
}

func TestXYZToDCS1(t *testing.T) {
	// Color profile with only one primary, and linear transfer function.
	colorProfile := &ColorProfileGeneral{
		PrimaryColors: TransformationLinDCSToXYZ{standardRGBRed},
	}
	colorProfile.MustInit()

	dcsVector := colorProfile.XYZToDCS(standardRGBRed.Scaled(0.4))

	want := DCSVector{0.4}
	if diff, err := want.Difference(dcsVector); err != nil || math.Abs(diff.ComponentSum()) > 0.000001 {
		t.Errorf("XYZToDCS() returned wrong device color space vector. Got %v, want %v", dcsVector, want)
	}
}

func TestXYZToDCS2(t *testing.T) {
	// Color profile with only two primaries, and linear transfer function.
	colorProfile := &ColorProfileGeneral{
		PrimaryColors: TransformationLinDCSToXYZ{standardRGBRed, standardRGBGreen},
	}
	colorProfile.MustInit()

	dcsVector := colorProfile.XYZToDCS(standardRGBGreen.Sum(standardRGBRed.Scaled(0.5)))

	want := DCSVector{0.5, 1}
	if diff, err := want.Difference(dcsVector); err != nil || math.Abs(diff.ComponentSum()) > 0.000001 {
		t.Errorf("XYZToDCS() returned wrong device color space vector. Got %v, want %v", dcsVector, want)
	}
}

func TestXYZToDCS3(t *testing.T) {
	// Color profile with sRGB primaries.
	colorProfile := &ColorProfileGeneral{
		WhitePointColor: CIE1931XYZAbs{}.Sum(standardRGBRed, standardRGBGreen, standardRGBBlue),
		PrimaryColors:   TransformationLinDCSToXYZ{standardRGBRed, standardRGBGreen, standardRGBBlue},
		OutputLimiter:   OutputLimiterSum{3},
		TransferFunc:    TransferFunctionStandardRGB,
	}
	colorProfile.MustInit()

	dcsVector := colorProfile.XYZToDCS(CIE1931XYZAbs{0.5, 0.4, 0.3})

	want := DCSVector{0.933728, 0.564098, 0.550101}
	if diff, err := want.Difference(dcsVector); err != nil || math.Abs(diff.ComponentSum()) > 0.000001 {
		t.Errorf("XYZToDCS() returned wrong device color space vector. Got %v, want %v", dcsVector, want)
	}
}
