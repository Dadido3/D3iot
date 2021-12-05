// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package wiz

import (
	"encoding/json"
	"fmt"
)

type Scene struct {
	id   uint   // ID of the scene.
	name string // English name of the scene.

	txDimming, txSpeed         bool // Parameters that need to be set when sending a pilot.
	rxDimming, rxSpeed, rxTemp bool // Parameters that are set when receiving a pilot from the device.
}

// TODO: Write test to check tx and rx parameters with a real device

var (
	SceneOcean        = Scene{txDimming: true, txSpeed: true, rxDimming: true, rxSpeed: true, rxTemp: false, id: 1, name: "Ocean"}           // Parameters: Dimming, Speed.
	SceneRomance      = Scene{txDimming: true, txSpeed: true, rxDimming: true, rxSpeed: true, rxTemp: false, id: 2, name: "Romance"}         // Parameters: Dimming, Speed.
	SceneSunset       = Scene{txDimming: true, txSpeed: true, rxDimming: true, rxSpeed: true, rxTemp: false, id: 3, name: "Sunset"}          // Parameters: Dimming, Speed.
	SceneParty        = Scene{txDimming: true, txSpeed: true, rxDimming: true, rxSpeed: true, rxTemp: false, id: 4, name: "Party"}           // Parameters: Dimming, Speed.
	SceneFireplace    = Scene{txDimming: true, txSpeed: true, rxDimming: true, rxSpeed: true, rxTemp: false, id: 5, name: "Fireplace"}       // Parameters: Dimming, Speed.
	SceneCozy         = Scene{txDimming: true, txSpeed: false, rxDimming: true, rxSpeed: false, rxTemp: false, id: 6, name: "Cozy"}          // Parameters: Dimming.
	SceneForest       = Scene{txDimming: true, txSpeed: true, rxDimming: true, rxSpeed: true, rxTemp: false, id: 7, name: "Forest"}          // Parameters: Dimming, Speed.
	ScenePastelColors = Scene{txDimming: true, txSpeed: true, rxDimming: true, rxSpeed: true, rxTemp: false, id: 8, name: "PastelColors"}    // Parameters: Dimming, Speed.
	SceneWakeUp       = Scene{txDimming: false, txSpeed: false, rxDimming: false, rxSpeed: false, rxTemp: false, id: 9, name: "WakeUp"}      // Takes 30 min to complete. Dimming is always set to 0. No parameters.
	SceneBedtime      = Scene{txDimming: true, txSpeed: false, rxDimming: true, rxSpeed: false, rxTemp: false, id: 10, name: "Bedtime"}      // Takes 30 min to complete. Parameters: Dimming.
	SceneWarmWhite    = Scene{txDimming: true, txSpeed: false, rxDimming: true, rxSpeed: false, rxTemp: true, id: 11, name: "WarmWhite"}     // Parameters: Dimming. Temperature doesn't need to be sent, but it is returned by the light.
	SceneDaylight     = Scene{txDimming: true, txSpeed: false, rxDimming: true, rxSpeed: false, rxTemp: true, id: 12, name: "Daylight"}      // Parameters: Dimming. Temperature doesn't need to be sent, but it is returned by the light.
	SceneCoolWhite    = Scene{txDimming: true, txSpeed: false, rxDimming: true, rxSpeed: false, rxTemp: true, id: 13, name: "CoolWhite"}     // Parameters: Dimming. Temperature doesn't need to be sent, but it is returned by the light.
	SceneNightLight   = Scene{txDimming: false, txSpeed: false, rxDimming: false, rxSpeed: false, rxTemp: false, id: 14, name: "NightLight"} // Dimming is always set to 0. Parameters: Nothing.
	SceneFocus        = Scene{txDimming: true, txSpeed: false, rxDimming: true, rxSpeed: false, rxTemp: false, id: 15, name: "Focus"}        // Parameters: Dimming.
	SceneRelax        = Scene{txDimming: true, txSpeed: false, rxDimming: true, rxSpeed: false, rxTemp: false, id: 16, name: "Relax"}        // Parameters: Dimming.
	SceneTrueColors   = Scene{txDimming: true, txSpeed: false, rxDimming: true, rxSpeed: false, rxTemp: false, id: 17, name: "TrueColors"}   // Parameters: Dimming.
	SceneTVTime       = Scene{txDimming: true, txSpeed: false, rxDimming: true, rxSpeed: false, rxTemp: false, id: 18, name: "TVTime"}       // Parameters: Dimming.
	ScenePlantgrowth  = Scene{txDimming: true, txSpeed: false, rxDimming: true, rxSpeed: false, rxTemp: false, id: 19, name: "Plantgrowth"}  // Parameters: Dimming.
	SceneSpring       = Scene{txDimming: true, txSpeed: true, rxDimming: true, rxSpeed: true, rxTemp: false, id: 20, name: "Spring"}         // Parameters: Dimming, Speed.
	SceneSummer       = Scene{txDimming: true, txSpeed: true, rxDimming: true, rxSpeed: true, rxTemp: false, id: 21, name: "Summer"}         // Parameters: Dimming, Speed.
	SceneFall         = Scene{txDimming: true, txSpeed: true, rxDimming: true, rxSpeed: true, rxTemp: false, id: 22, name: "Fall"}           // Parameters: Dimming, Speed.
	SceneDeepdive     = Scene{txDimming: true, txSpeed: true, rxDimming: true, rxSpeed: true, rxTemp: false, id: 23, name: "Deepdive"}       // Parameters: Dimming, Speed.
	SceneJungle       = Scene{txDimming: true, txSpeed: true, rxDimming: true, rxSpeed: true, rxTemp: false, id: 24, name: "Jungle"}         // Parameters: Dimming, Speed.
	SceneMojito       = Scene{txDimming: true, txSpeed: true, rxDimming: true, rxSpeed: true, rxTemp: false, id: 25, name: "Mojito"}         // Parameters: Dimming, Speed.
	SceneClub         = Scene{txDimming: true, txSpeed: true, rxDimming: true, rxSpeed: true, rxTemp: false, id: 26, name: "Club"}           // Parameters: Dimming, Speed.
	SceneChristmas    = Scene{txDimming: true, txSpeed: true, rxDimming: true, rxSpeed: true, rxTemp: false, id: 27, name: "Christmas"}      // Parameters: Dimming, Speed.
	SceneHalloween    = Scene{txDimming: true, txSpeed: true, rxDimming: true, rxSpeed: true, rxTemp: false, id: 28, name: "Halloween"}      // Parameters: Dimming, Speed.
	SceneCandlelight  = Scene{txDimming: true, txSpeed: false, rxDimming: true, rxSpeed: false, rxTemp: false, id: 29, name: "Candlelight"}  // Parameters: Dimming.
	SceneGoldenWhite  = Scene{txDimming: true, txSpeed: true, rxDimming: true, rxSpeed: true, rxTemp: false, id: 30, name: "GoldenWhite"}    // Parameters: Dimming, Speed.
	ScenePulse        = Scene{txDimming: true, txSpeed: true, rxDimming: true, rxSpeed: true, rxTemp: false, id: 31, name: "Pulse"}          // Parameters: Dimming, Speed.
	SceneSteampunk    = Scene{txDimming: true, txSpeed: true, rxDimming: true, rxSpeed: true, rxTemp: false, id: 32, name: "Steampunk"}      // Parameters: Dimming, Speed.
	SceneRhythm       = Scene{txDimming: false, txSpeed: false, rxDimming: false, rxSpeed: false, rxTemp: false, id: 1000, name: "Rhythm"}   // Not tested yet.
)

var scenes = map[uint]Scene{
	1:    SceneOcean,
	2:    SceneRomance,
	3:    SceneSunset,
	4:    SceneParty,
	5:    SceneFireplace,
	6:    SceneCozy,
	7:    SceneForest,
	8:    ScenePastelColors,
	9:    SceneWakeUp,
	10:   SceneBedtime,
	11:   SceneWarmWhite,
	12:   SceneDaylight,
	13:   SceneCoolWhite,
	14:   SceneNightLight,
	15:   SceneFocus,
	16:   SceneRelax,
	17:   SceneTrueColors,
	18:   SceneTVTime,
	19:   ScenePlantgrowth,
	20:   SceneSpring,
	21:   SceneSummer,
	22:   SceneFall,
	23:   SceneDeepdive,
	24:   SceneJungle,
	25:   SceneMojito,
	26:   SceneClub,
	27:   SceneChristmas,
	28:   SceneHalloween,
	29:   SceneCandlelight,
	30:   SceneGoldenWhite,
	31:   ScenePulse,
	32:   SceneSteampunk,
	1000: SceneRhythm,
}

// MarshalJSON implements the JSON marshaler interface.
func (s Scene) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.id)
}

// UnmarshalJSON implements the JSON unmarshaler interface.
func (s *Scene) UnmarshalJSON(data []byte) error {
	var id uint
	if err := json.Unmarshal(data, &id); err != nil {
		return err
	}

	// Write data of known scene, if there is one.
	if scene, ok := scenes[id]; ok {
		*s = scene
		return nil
	}

	// Scene not found, generate new general scene.
	*s = Scene{
		id:   id,
		name: "Unknown",
	}

	return nil
}

func (s Scene) String() string {
	return fmt.Sprintf("%d %q", s.id, s.name)
}

// NeedsDimming returns true if you need to supply a dimming value along with the scene in the pilot.
func (s Scene) NeedsDimming() bool {
	return s.txDimming
}

// NeedsSpeed returns true if you need to supply a speed value along with the scene in the pilot.
func (s Scene) NeedsSpeed() bool {
	return s.txSpeed
}

// HasDimming returns true if the device sends a dimming value along with the scene in the pilot.
func (s Scene) HasDimming() bool {
	return s.rxDimming
}

// HasSpeed returns true if the device sends a speed value along with the scene in the pilot.
func (s Scene) HasSpeed() bool {
	return s.rxSpeed
}

// HasTemperature returns true if the device sends a color temperature along with the scene in the pilot.
func (s Scene) HasTemperature() bool {
	return s.rxTemp
}
