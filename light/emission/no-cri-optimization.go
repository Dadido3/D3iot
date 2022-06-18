// Copyright (c) 2022 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package emission

// NoWhiteOptimization can be used to wrap any emission value to disable any white emitter optimization later in the color management pipeline.
//
// Wrapping such a value will cause the emission value to be only constructed out of primary color emitters (e.g. red, green, blue), and no white emitters.
// This will not change the emitted color, but the maximum brightness may be reduced and the CRI is not as good as it could be.
//
// One use-case could be to eliminate any timing discrepancy between white and primary color emitters for fast changing colors, as these two classes of emitters may be filtered (by a low-pass) in a different way.
//
// This does not work on DCS vectors, even though they implement the emission.Value interface.
type NoWhiteOptimization struct {
	EmissionValue Value
}

var _ Value = &NoWhiteOptimization{}

// IntoDCS implements the Value interface.
func (n NoWhiteOptimization) IntoDCS(cp ColorProfile) DCSVector {
	nCP := cp.NoWhiteOptimizationColorProfile()

	return n.EmissionValue.IntoDCS(nCP)
}
