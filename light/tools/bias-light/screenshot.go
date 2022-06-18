// Copyright (c) 2021-2022 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"fmt"
	"image"

	"github.com/kbinani/screenshot"
)

func takeScreenshot() (*image.RGBA, error) {
	n := screenshot.NumActiveDisplays()

	if n <= 0 {
		return nil, fmt.Errorf("couldn't determine display to capture from")
	}

	bounds := screenshot.GetDisplayBounds(0)

	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		return nil, fmt.Errorf("failed to capture screenshot: %w", err)
	}

	return img, nil
}
