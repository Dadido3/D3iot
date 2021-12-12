// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package emission

import "math"

// clamp01 returns v clamped into the interval [0,1].
func clamp01(v float64) float64 {
	return math.Min(math.Max(v, 0), 1)
}
