// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package emission

import (
	"testing"
)

// Scaled sRGB primaries with a whitepoint of D65.
var (
	standardRGBRed   = CIE1931XYZColor{0.4124564, 0.2126729, 0.0193339}
	standardRGBGreen = CIE1931XYZColor{0.3575761, 0.7151522, 0.1191920}
	standardRGBBlue  = CIE1931XYZColor{0.1804375, 0.0721750, 0.9503041}
)

// Module descriptor with sRGB primaries, but linear transfer function.
var moduleDescriptorLinearStandardRGB = &ModuleDescriptorGeneral{
	PrimaryColors:     TransformationLinDCSToXYZ{standardRGBRed, standardRGBGreen, standardRGBBlue},
	LinearDCSSumLimit: 3,
}

// Module descriptor with only two primaries, and linear transfer function.
var moduleDescriptorLinearRG = &ModuleDescriptorGeneral{
	PrimaryColors:     TransformationLinDCSToXYZ{standardRGBRed, standardRGBGreen},
	LinearDCSSumLimit: 2,
}

// Module descriptor with only one primary, and linear transfer function.
var moduleDescriptorLinearR = &ModuleDescriptorGeneral{
	PrimaryColors:     TransformationLinDCSToXYZ{standardRGBRed},
	LinearDCSSumLimit: 1,
}

func TestDCSToXYZ(t *testing.T) {
	color, err := moduleDescriptorLinearStandardRGB.DCSToXYZ([]float64{1, 0, 0})
	if err != nil {
		t.Fatalf("DCSToXYZ() failed: %v", err)
	}

	if color != standardRGBRed {
		t.Errorf("DCSToXYZ() returned wrong color result. Got %v, want %v", color, standardRGBRed)
	}
}

func TestXYZToDCS1(t *testing.T) {
	dcsValue, err := moduleDescriptorLinearR.XYZToDCS(standardRGBRed.Scaled(0.4))
	if err != nil {
		t.Fatalf("XYZToDCS() failed: %v", err)
	}

	if dcsValue.Channels() != 1 || compareFloat64(dcsValue[0], 0.4) {
		t.Errorf("XYZToDCS() returned wrong device color space vector. Got %v, want %v", dcsValue, []float64{0.4})
	}
}

func TestXYZToDCS2(t *testing.T) {
	dcsValue, err := moduleDescriptorLinearRG.XYZToDCS(standardRGBGreen.Sum(standardRGBRed.Scaled(0.5)))
	if err != nil {
		t.Fatalf("XYZToDCS() failed: %v", err)
	}

	if dcsValue.Channels() != 2 || compareFloat64(dcsValue[0], 0.5) || compareFloat64(dcsValue[1], 1) {
		t.Errorf("XYZToDCS() returned wrong device color space vector. Got %v, want %v", dcsValue, []float64{0.5, 1})
	}
}
