// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package emission

import (
	"fmt"
	"testing"
)

func TestBlackBodyFixed_CIE1931XYZAbs(t *testing.T) {
	for temp := 1700; temp <= 25000; temp += 100 {

		t.Run(fmt.Sprintf("%d K", temp), func(t *testing.T) {

			modelResult := BlackBodyFixed{Temperature: float64(temp), Luminance: 1}.CIE1931XYZAbs()
			integratedResult := BlackBodyArea{Temperature: float64(temp), Area: 1}.CIE1931XYZAbs()

			// Convert into relative CIE 1931 XYZ.
			modelRelative := modelResult.Relative(1)
			integratedRelative := integratedResult.Relative(integratedResult.Y)

			dist := integratedRelative.CIE1976LABDistance(modelRelative, StandardIlluminantD65.CIE1931XYZRel())
			if dist > 1 {
				t.Errorf("Blackbody model result %v differs from integrated result %v. ΔE* = %f", modelRelative.CIE1931xyYRel(), integratedRelative.CIE1931xyYRel(), dist)
			}
		})
	}
}
