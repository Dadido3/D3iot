package main

import (
	"D3domotics/light/wiz"
	"log"
	"time"
)

func main() {
	light := wiz.NewLight("192.168.1.174:38899")

	// Interpolate between these pilots.
	pilots := []wiz.Pilot{
		//wiz.Pilot{}.WithRGBW(0, 0, 0, 0, 45, 100),    // Warm white with a good CRI.
		wiz.Pilot{}.WithRGBW(0, 0, 0, 45, 0, 100),    // Cold white with a good CRI.
		wiz.Pilot{}.WithRGBW(80, 100, 40, 0, 0, 100), // Same cold white simlated with RGB colors. Bad CRI.

		// "Fire" sequence.
		/*wiz.Pilot{}.WithRGBW(0, 0, 0, 0, 40, 100),
		wiz.Pilot{}.WithRGBW(20, 0, 0, 0, 40, 100),
		wiz.Pilot{}.WithRGBW(30, 0, 0, 10, 30, 100),
		wiz.Pilot{}.WithRGBW(10, 0, 0, 0, 50, 100),*/
	}

	/*if err := light.Pulse(1000, 100*time.Millisecond); err != nil {
		log.Printf("light.Pulse() failed: %v", err)
	}*/

	/*if result, err := light.GetUserConfig(); err != nil {
		log.Printf("query failed: %v", err)
	} else {
		log.Printf("%#v", result)
	}*/

	/*if err := light.SetPilot(wiz.Pilot{}.WithRGBW(50, 0, 0, 0, 00, 100)); err != nil {
		log.Printf("light.SetPilot() failed: %v", err)
	}*/

	/*if err := light.SetPilot(wiz.Pilot{}.WithScene(wiz.SceneNoScene, 6000, 50, 50)); err != nil {
		log.Printf("light.SetPilot() failed: %v", err)
	}*/

	// Init first pilot for blending/mixing.
	p1 := pilots[0]

	steps := 10

	for {
		for _, p2 := range pilots {
			for i := 0; i < steps; i++ {

				factor1, factor2 := 1-float64(i)/float64(steps), float64(i)/float64(steps)
				if p1.HasRGBW() && p2.HasRGBW() {

					r := uint8(float64(*p1.R)*factor1 + float64(*p2.R)*factor2)
					g := uint8(float64(*p1.G)*factor1 + float64(*p2.G)*factor2)
					b := uint8(float64(*p1.B)*factor1 + float64(*p2.B)*factor2)
					cw := uint8(float64(*p1.CW)*factor1 + float64(*p2.CW)*factor2)
					ww := uint8(float64(*p1.WW)*factor1 + float64(*p2.WW)*factor2)
					dimming := uint(float64(p1.Dimming)*factor1 + float64(p2.Dimming)*factor2)

					mixedPilot := wiz.Pilot{}.WithRGBW(r, g, b, cw, ww, dimming)

					if err := light.SetPilot(mixedPilot); err != nil {
						log.Printf("light.SetPilot() failed: %v", err)
					}
				}
				time.Sleep(100 * time.Millisecond)
			}
			p1 = p2
		}
	}
}
