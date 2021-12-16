// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"log"
	"math"
	"time"

	"github.com/Dadido3/D3iot/light/drivers/wiz"
	"github.com/Dadido3/D3iot/light/emission"
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

	//colorProfile := light.ColorProfiles()[0]

	//emissionValue := emission.CIE1931XYZRel{X: 0.95047, Y: 1, Z: 1.08883}
	//emissionValue := colorProfile.WhitePoint().Scaled(0.1)
	//emissionValue := emission.StandardIlluminantA.Absolute(200)
	emissionValue := emission.BlackBodyArea{Temperature: 2000, Area: 0.15}

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

	//return

	frequency := 0.02 // In 1/s

	for {
		seconds := float64(time.Now().UnixNano()) * 1e-9
		sineWave := math.Sin(frequency * 2 * math.Pi * seconds)
		temp := 2300 + 1000*sineWave
		emissionValue := emission.BlackBodyArea{Temperature: temp, Area: 2}

		if err := light.SetColors(emissionValue); err != nil {
			log.Printf("light.SetColors() failed: %v", err)
		}

		time.Sleep(10 * time.Millisecond)
	}
}
