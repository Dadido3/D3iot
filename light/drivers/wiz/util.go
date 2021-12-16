// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

// TODO: Replace functions in here with go native alternatives once they are available.

package wiz

import "math"

// newPtrUInt returns a pointer to the given uint.
func newPtrUInt(n uint) *uint {
	return &n
}

// clamp01 returns v clamped into the interval [0,1].
func clamp01(v float64) float64 {
	return math.Min(math.Max(v, 0), 1)
}

// normFloatToInt converts the float v in the range [0, 1] to an integer in the range [0, maxInt].
//
// Nearest neighbor rounding is used.
// The float v is clamped to [0, 1].
func normFloatToInt(v float64, maxInt int) int {
	return int(clamp01(v)*float64(maxInt) + 0.5)
}

// normFloatToUint8 converts the float v in the range [0, 1] to an uint8 in the range [0, 255].
//
// Nearest neighbor rounding is used.
// The float v is clamped to [0, 1].
func normFloatToUint8(v float64) uint8 {
	return uint8(normFloatToInt(v, 255))
}
