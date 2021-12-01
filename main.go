package main

import (
	"D3domotics/light/wiz"
	"fmt"
	"log"
	"time"
)

func main() {
	light := wiz.NewLight("192.168.1.174:38899")

	// Interpolate between these pilots.
	pilots := []wiz.Pilot{
		//{State: true, R: 0, G: 0, B: 0, C: 0, W: 45, Dimming: 100}, // Warm white with a good CRI.
		//{State: true, R: 0, G: 0, B: 0, C: 45, W: 0, Dimming: 30}, // Cold white with a good CRI.
		//{State: true, R: 80, G: 100, B: 40, C: 0, W: 0, Dimming: 30}, // Same cold white simlated with RGB colors. Bad CRI.

		// "Fire" sequence.
		{State: true, R: 0, G: 0, B: 0, CW: 0, WW: 40, Dimming: 100},
		{State: true, R: 20, G: 0, B: 0, CW: 0, WW: 40, Dimming: 100},
		{State: true, R: 30, G: 0, B: 0, CW: 10, WW: 30, Dimming: 100},
		{State: true, R: 10, G: 0, B: 0, CW: 0, WW: 50, Dimming: 100},
	}

	// Init old pilot for mixing.
	p1 := wiz.Pilot{State: true, R: 0, G: 0, B: 0, CW: 0, WW: 0, Dimming: 100}

	steps := 10

	/*if err := light.Pulse(100, 100*time.Millisecond); err != nil {
		log.Printf("light.Pulse() failed: %v", err)
	}*/

	if result, err := light.GetUserConfig(); err != nil {
		log.Printf("query failed: %v", err)
	} else {
		log.Printf("%#v", result)
	}

	if true {
		return
	}

	for {
		for _, p2 := range pilots {
			for i := 0; i < steps; i++ {

				a, b := 1-float64(i)/float64(steps), float64(i)/float64(steps)
				mixedPilot := wiz.Pilot{

					State:   true,
					R:       uint8(float64(p1.R)*a + float64(p2.R)*b),
					G:       uint8(float64(p1.G)*a + float64(p2.G)*b),
					B:       uint8(float64(p1.B)*a + float64(p2.B)*b),
					CW:      uint8(float64(p1.CW)*a + float64(p2.CW)*b),
					WW:      uint8(float64(p1.WW)*a + float64(p2.WW)*b),
					Dimming: uint(float64(p1.Dimming)*a + float64(p2.Dimming)*b),
				}

				fmt.Println(mixedPilot)

				if err := light.SetPilot(mixedPilot); err != nil {
					log.Printf("light.SetPilot() failed: %v", err)
				}
				time.Sleep(100 * time.Millisecond)
			}
			p1 = p2
		}
	}
}
