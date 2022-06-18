// Copyright (c) 2021-2022 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"flag"
	"image"
	"log"
	"time"

	"github.com/Dadido3/D3iot/light"
	"github.com/Dadido3/D3iot/light/drivers/wiz"
	"github.com/Dadido3/D3iot/light/emission"
	"golang.org/x/image/draw"
)

var flagDeviceWiZ = flag.String("device-wiz", "wiz-d47cf3:38899", "The address of the WiZ device to be profiled. Example: \"--device-wiz wiz-123abc:38899\" or \"--device-wiz 192.168.1.123:38899\".")
var flagMaxLuminance = flag.Float64("max-luminance", 0, "The maximum luminance that will be output for a fully white screen in lumens.")

func main() {
	flag.Parse()

	// Connect to light device.
	var light light.Light
	switch {
	case *flagDeviceWiZ != "":
		var err error
		if light, err = wiz.NewLight(*flagDeviceWiZ); err != nil {
			log.Printf("wiz.NewLight(%q) failed: %v", *flagDeviceWiZ, err)
			return
		}

	default:
		log.Printf("No device address given. Start program with any \"--device-...\" parameter set.")
		log.Printf("Example: bias-light --device-wiz wiz-123abc:38899")
		return
	}

	for {
		srcImg, err := takeScreenshot()
		if err != nil {
			log.Printf("Failed to take screenshot: %s", err)
			continue
		}

		// Scale image to 1x1 pixels.
		dstImg := image.NewNRGBA(image.Rect(0, 0, 1, 1))
		draw.CatmullRom.Scale(dstImg, dstImg.Bounds(), srcImg, srcImg.Bounds(), draw.Over, &draw.Options{})

		r, g, b, _ := dstImg.At(0, 0).RGBA()
		emissionValue := emission.StandardRGB{R: float64(r) / 65535.0, G: float64(g) / 65535.0, B: float64(b) / 65535.0}
		if *flagMaxLuminance > 0 {
			// Convert into absolute emission value and set the first module's color.
			light.SetColors(emissionValue.CIE1931XYZRel().Absolute(*flagMaxLuminance)) // TODO: Set all available modules.
		} else {
			// Use relative emission value and set the first module's color.
			light.SetColors(emissionValue) // TODO: Set all available modules.
		}

		time.Sleep(20 * time.Millisecond)
	}
}