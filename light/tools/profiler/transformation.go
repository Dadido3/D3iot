// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"fmt"

	"github.com/Dadido3/D3iot/light/emission"
	"gonum.org/v1/gonum/optimize"
)

// primariesToTransformation returns a transformation from the given list of floats.
// len(x) must be a multiple of 3!
func primariesToTransformation(x []float64) emission.TransformationLinDCSToXYZ {
	transformation := make(emission.TransformationLinDCSToXYZ, len(x)/3)

	for i := range transformation {
		transformation[i] = emission.CIE1931XYZAbs{X: x[i*3+0], Y: x[i*3+1], Z: x[i*3+2]}
	}

	return transformation
}

// calculateTransformation takes a list of matches between the DCS and CIE 1931 XYZ colors, and returns a transformation that represents this transformation with the least total ΔE.
func calculateTransformation(matches []Match) (emission.TransformationLinDCSToXYZ, float64, error) {
	if len(matches) < 1 {
		return nil, 0, fmt.Errorf("no matches given to calculate primaries with")
	}
	channels := matches[0].dcs.Channels()
	if len(matches) < channels {
		return nil, 0, fmt.Errorf("not enough matches to calculate primaries. Got %d, want %d", len(matches), channels)
	}

	// The first color is the brightest color.
	whitePoint := matches[0].xyz
	whitePointLumen := *flagMaxLuminance

	// Function to optimize.
	optimizeFunc := func(x []float64) float64 {

		transformation := primariesToTransformation(x)

		ssr := 0.0

		// Penalty for negative primaries.
		/*for _, x := range x {
			if x < 0 {
				scaled := x * 10
				ssr += scaled * scaled
			}
		}*/

		// Calculate distance between matches and DCT --> XYZ transformed colors.
		for _, match := range matches {
			// Color space transformation.
			transformedXYZ, _ := transformation.Multiplied(match.dcs)
			transformedXYZRel := transformedXYZ.Relative(whitePointLumen / whitePoint.Y)

			distSqr := match.xyz.CIE1976LABDistanceSqr(transformedXYZRel, emission.StandardIlluminantD65)

			ssr += distSqr
		}

		return ssr
	}

	p := optimize.Problem{
		Func: optimizeFunc,
	}

	// Set up the initial parameters.
	init := make([]float64, channels*3)
	for i := range init {
		init[i] = 1
	}

	res, err := optimize.Minimize(p, init, nil, &optimize.NelderMead{})
	if err != nil {
		return nil, 0, err
	}
	if err = res.Status.Err(); err != nil {
		return nil, 0, fmt.Errorf("optimization status error: %w", err)
	}

	transformation := primariesToTransformation(res.X)

	return transformation, optimizeFunc(res.X), nil
}
