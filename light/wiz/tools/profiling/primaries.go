// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"fmt"

	"github.com/lucasb-eyer/go-colorful"
	"gonum.org/v1/gonum/optimize"
)

// RGBWPrimaries contains the primaries that a specific WiZ light device uses.
// This basically describes the color space.
type RGBWPrimaries struct {
	R, G, B, CW, WW CIE1931XYZColor
}

func NewPrimariesFromValues(values [15]float64) RGBWPrimaries {
	return RGBWPrimaries{
		R:  CIE1931XYZColor{values[0], values[1], values[2]},
		G:  CIE1931XYZColor{values[3], values[4], values[5]},
		B:  CIE1931XYZColor{values[6], values[7], values[8]},
		CW: CIE1931XYZColor{values[9], values[10], values[11]},
		WW: CIE1931XYZColor{values[12], values[13], values[14]},
	}
}

// Scale returns p with all color primaries scaled by s.
func (p RGBWPrimaries) Scale(s float64) RGBWPrimaries {
	return RGBWPrimaries{p.R.Scale(s), p.G.Scale(s), p.B.Scale(s), p.CW.Scale(s), p.WW.Scale(s)}
}

// Normalize returns p scaled such that value.CIE1931XYZColor(p).Y equals 1.
//
// If value is {0, 0, 0, 255, 255}, a light bulb that is fully turned on (Only cold and warm white, no RGB) has a CIE 1931 XYZ color of {something, 1, something}.
// In this case the illuminant is not D65, but some color between the cold and warm white LED.
func (p RGBWPrimaries) Normalize(value RGBWValue) RGBWPrimaries {
	c := value.CIE1931XYZColor(p)
	return p.Scale(1 / c.Y)
}

func (p RGBWPrimaries) Values() [15]float64 {
	return [15]float64{
		p.R.X, p.R.Y, p.R.Z,
		p.G.X, p.G.Y, p.G.Z,
		p.B.X, p.B.Y, p.B.Z,
		p.CW.X, p.CW.Y, p.CW.Z,
		p.WW.X, p.WW.Y, p.WW.Z,
	}
}

// CalculatePrimaries takes a list of matches between the device color space and CIE 1931 XYZ, and returns the color primaries of the WiZ light.
func CalculatePrimaries(matches map[RGBWValue]CIE1931XYZColor) (RGBWPrimaries, float64, error) {

	if len(matches) < 5 {
		return RGBWPrimaries{}, 0, fmt.Errorf("not enough matches to calculate primaries. Got %d, want %d", len(matches), 5)
	}

	// Function to optimize.
	optimizeFunc := func(x []float64) float64 {

		xFixed := (*[15]float64)(x)
		primaries := NewPrimariesFromValues(*xFixed)

		ssr := 0.0

		// Penalty for negative primaries.
		for _, x := range xFixed {
			if x < 0 {
				scaled := x * 10
				ssr += scaled * scaled
			}
		}

		for matchRGBW, matchXYZ := range matches {

			transformedXYZ := matchRGBW.CIE1931XYZColor(primaries)
			l1, a1, b1 := colorful.XyzToLab(transformedXYZ.X, transformedXYZ.Y, transformedXYZ.Z)
			l2, a2, b2 := colorful.XyzToLab(matchXYZ.X, matchXYZ.Y, matchXYZ.Z)

			lDiff, aDiff, bDiff := l1-l2, a1-a2, b1-b2
			ssr += lDiff*lDiff + aDiff*aDiff + bDiff*bDiff
		}

		return ssr
	}

	p := optimize.Problem{
		Func: optimizeFunc,
	}

	// Set up the initial parameters.
	init := make([]float64, 15)
	for i := range init {
		init[i] = 0.5
	}

	res, err := optimize.Minimize(p, init, nil, &optimize.NelderMead{SimplexSize: 0.5})
	if err != nil {
		return RGBWPrimaries{}, 0, err
	}
	if err = res.Status.Err(); err != nil {
		return RGBWPrimaries{}, 0, fmt.Errorf("optimization status error: %w", err)
	}

	xFixed := (*[15]float64)(res.X)

	primaries := NewPrimariesFromValues(*xFixed)
	normalizedPrimaries := primaries.Normalize(RGBWValue{0, 0, 0, 255, 255})

	return normalizedPrimaries, optimizeFunc(res.X), nil
}
