// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package emission

// OutputLimiter limits the light output of a module/device by specific rules.
// Some devices may already limit their total output internally, so if you send a DCSVector{1,1,1,1,1} to a device, its real output may just be DCSVector{0.4, 0.4, 0.4, 0.4, 0.4}.
// The limiter can reproduce this behavior, which allows the lib to work with correct brightness values.
type OutputLimiter interface {
	LimitDCS(v LinDCSVector) LinDCSVector // Returns a version that is in some way modified, so that it obeys the limits.
}

// OutputLimiterSum implements OutputLimiter in a way that scales DCS values so that their sum is below a given limit.
type OutputLimiterSum struct {
	Limit float64
}

// LimitDCS implements OutputLimiter.
func (ol OutputLimiterSum) LimitDCS(v LinDCSVector) LinDCSVector {
	sum := v.ComponentSum()

	// Scale it so that the sum of all values doesn't exceed the limit.
	if sum > ol.Limit {
		return v.Scaled(ol.Limit / float64(sum))
	}

	return v
}
