// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"log"
	"time"

	"github.com/Dadido3/D3iot/light/emission"
	"github.com/Dadido3/D3iot/light/wiz"
)

func main() {
	light, err := wiz.NewLight("wiz-d47cf3:38899")
	if err != nil {
		log.Printf("wiz.NewLight() failed: %v", err)
		return
	}

	/*if err := light.Pulse(50, 100*time.Millisecond); err != nil {
		log.Printf("light.Pulse() failed: %v", err)
	}*/

	/*if result, err := light.GetSystemConfig(); err != nil {
		log.Printf("light.GetSystemConfig() failed: %v", err)
	} else {
		log.Printf("%#v", result)
	}*/

	/*if info, err := light.GetDeviceInfo(); err != nil {
		log.Printf("light.GetDeviceInfo() failed: %v", err)
	} else {
		log.Printf("%#v", info)
	}*/

	/*if info, err := light.GetUserConfig(); err != nil {
		log.Printf("light.GetUserConfig() failed: %v", err)
	} else {
		log.Printf("%#v", info)
	}*/

	/*if favs, err := light.GetFavs(); err != nil {
		log.Printf("light.GetFavs() failed: %v", err)
	} else {
		log.Printf("Favs: %v", favs)
	}*/

	/*if product, err := light.Product(); err != nil {
		log.Printf("light.Product() failed: %v", err)
	} else {
		log.Printf("Product: %v", product)
	}*/

	/*//scene := wiz.SceneBedtime
	//temp := uint(4200)
	r := uint8(255)
	dimming := uint(10)
	pilot := wiz.Pilot{
		State:   true,
		Dimming: &dimming,
		R:       &r,
		//Scene:   &scene,
	}*/
	//pilot := wiz.NewPilotWithRGBW(100, 0, 0, 0, 50, 0)
	//pilot := wiz.NewPilot(false)
	//pilot := wiz.NewPilotWithScene(wiz.SceneCozy, 20, 0)
	/*if err := light.SetPilot(pilot); err != nil {
		log.Printf("light.SetPilot() failed: %v", err)
	}
	log.Printf("Set pilot to %v", pilot)*/

	/*for {
		if result, err := light.GetPilot(); err != nil {
			log.Printf("query failed: %v", err)
		} else {
			log.Printf("%s", result)
		}
		time.Sleep(1 * time.Second)
	}*/

	/*if err := light.SetPilot(wiz.Pilot{}.WithRGBW(0, 0, 0, 50, 0)); err != nil {
		log.Printf("light.SetPilot() failed: %v", err)
	}*/

	/*if err := light.SetPilot(wiz.Pilot{}.WithScene(wiz.SceneClub, 6000, 50, 50)); err != nil {
		log.Printf("light.SetPilot() failed: %v", err)
	}*/

	//moduleProfile := light.ModuleProfiles()[0]

	//emissionValue := emission.CIE1931XYZAbs{X: 0.95047, Y: 1, Z: 1.08883}.Scaled(1521 * 0.1)
	//emissionValue := moduleProfile.WhitePoint().Scaled(0.13)
	//emissionValue := emission.StandardIlluminantA.Scaled(200)
	emissionValue := emission.BlackBodyFixed{Temperature: 2800, Luminance: 200}

	if err := light.SetColors(emissionValue); err != nil {
		log.Printf("light.SetColors() failed: %v", err)
	}

	var res emission.CIE1931XYZAbs
	if err := light.GetColors(&res); err != nil {
		log.Printf("light.GetColors() failed: %v", err)
	} else {
		log.Printf("Returned colors: %v", res)
	}

	if pilot, err := light.GetPilot(); err != nil {
		log.Printf("light.GetPilot() failed: %v", err)
	} else {
		log.Printf("Returned pilot: %v", pilot)
	}

	if true {
		return
	}

	// Interpolate between these pilots.
	pilots := []wiz.Pilot{
		//wiz.NewPilot(true).WithDimming(100).WithRGBW(0, 0, 0, 0, 45),    // Warm white with a good CRI.
		wiz.NewPilot(true).WithDimming(100).WithRGBW(0, 0, 0, 0, 170), // Cold white with a good CRI.
		//wiz.NewPilot(true).WithDimming(100).WithRGBW(90, 120, 50, 0, 0), // Same cold white simlated with RGB colors. Bad CRI.

		// "Fire" sequence.
		//wiz.NewPilot(true).WithDimming(100).WithRGBW(0, 0, 0, 0, 40),
		//wiz.NewPilot(true).WithDimming(100).WithRGBW(20, 0, 0, 0, 40),
		//wiz.NewPilot(true).WithDimming(100).WithRGBW(30, 0, 0, 10, 30),
		//wiz.NewPilot(true).WithDimming(100).WithRGBW(10, 0, 0, 0, 50),

		//wiz.NewPilot(true).WithDimming(100).WithRGBW(255, 0, 0, 0, 0),
		//wiz.NewPilot(true).WithDimming(100).WithRGBW(0, 255, 0, 0, 0),
		//wiz.NewPilot(true).WithDimming(100).WithRGBW(0, 0, 255, 0, 0),
	}

	// Init first pilot for blending/mixing.
	p1 := pilots[0]

	steps := 1

	for {
		for _, p2 := range pilots {
			for i := 0; i < steps; i++ {

				factor1, factor2 := 1-float64(i)/float64(steps), float64(i)/float64(steps)
				if p1.HasRGBW() && p2.HasRGBW() && p1.HasDimming() && p2.HasDimming() {

					r := uint8(float64(*p1.R)*factor1 + float64(*p2.R)*factor2)
					g := uint8(float64(*p1.G)*factor1 + float64(*p2.G)*factor2)
					b := uint8(float64(*p1.B)*factor1 + float64(*p2.B)*factor2)
					cw := uint8(float64(*p1.CW)*factor1 + float64(*p2.CW)*factor2)
					ww := uint8(float64(*p1.WW)*factor1 + float64(*p2.WW)*factor2)
					dimming := uint(float64(*p1.Dimming)*factor1 + float64(*p2.Dimming)*factor2)

					mixedPilot := wiz.NewPilotWithRGBW(dimming, r, g, b, cw, ww)

					if err := light.SetPilot(mixedPilot); err != nil {
						log.Printf("light.SetPilot() failed: %v", err)
					}
				}
				time.Sleep(1000 * time.Millisecond)
			}
			p1 = p2
		}
	}
}
