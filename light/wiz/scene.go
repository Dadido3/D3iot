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
	SceneOcean        = Scene{txDimming: true, txSpeed: true, rxDimming: true, rxSpeed: true, rxTemp: false, id: 1, name: "Ocean"}           //
	SceneRomance      = Scene{txDimming: true, txSpeed: true, rxDimming: true, rxSpeed: true, rxTemp: false, id: 2, name: "Romance"}         //
	SceneSunset       = Scene{txDimming: true, txSpeed: true, rxDimming: true, rxSpeed: true, rxTemp: false, id: 3, name: "Sunset"}          //
	SceneParty        = Scene{txDimming: true, txSpeed: true, rxDimming: true, rxSpeed: true, rxTemp: false, id: 4, name: "Party"}           //
	SceneFireplace    = Scene{txDimming: true, txSpeed: true, rxDimming: true, rxSpeed: true, rxTemp: false, id: 5, name: "Fireplace"}       //
	SceneCozy         = Scene{txDimming: true, txSpeed: true, rxDimming: true, rxSpeed: false, rxTemp: false, id: 6, name: "Cozy"}           //
	SceneForest       = Scene{txDimming: true, txSpeed: true, rxDimming: true, rxSpeed: true, rxTemp: false, id: 7, name: "Forest"}          //
	ScenePastelColors = Scene{txDimming: true, txSpeed: true, rxDimming: true, rxSpeed: true, rxTemp: false, id: 8, name: "PastelColors"}    //
	SceneWakeUp       = Scene{txDimming: true, txSpeed: false, rxDimming: false, rxSpeed: false, rxTemp: false, id: 9, name: "WakeUp"}       // Takes 30 min to complete. // TODO: Check if the dimming value has any influence on it.
	SceneBedtime      = Scene{txDimming: true, txSpeed: false, rxDimming: true, rxSpeed: false, rxTemp: false, id: 10, name: "Bedtime"}      // Takes 30 min to complete.
	SceneWarmWhite    = Scene{txDimming: true, txSpeed: false, rxDimming: true, rxSpeed: false, rxTemp: true, id: 11, name: "WarmWhite"}     //
	SceneDaylight     = Scene{txDimming: true, txSpeed: false, rxDimming: true, rxSpeed: false, rxTemp: true, id: 12, name: "Daylight"}      //
	SceneCoolWhite    = Scene{txDimming: true, txSpeed: false, rxDimming: true, rxSpeed: false, rxTemp: true, id: 13, name: "CoolWhite"}     //
	SceneNightLight   = Scene{txDimming: false, txSpeed: false, rxDimming: false, rxSpeed: false, rxTemp: false, id: 14, name: "NightLight"} // Just a dim light, there is no control over anything.
	SceneFocus        = Scene{txDimming: true, txSpeed: false, rxDimming: true, rxSpeed: false, rxTemp: false, id: 15, name: "Focus"}        //
	SceneRelax        = Scene{txDimming: true, txSpeed: false, rxDimming: true, rxSpeed: false, rxTemp: false, id: 16, name: "Relax"}        //
	SceneTrueColors   = Scene{txDimming: true, txSpeed: false, rxDimming: true, rxSpeed: false, rxTemp: false, id: 17, name: "TrueColors"}   //
	SceneTVTime       = Scene{txDimming: true, txSpeed: false, rxDimming: true, rxSpeed: false, rxTemp: false, id: 18, name: "TVTime"}       //
	ScenePlantgrowth  = Scene{txDimming: true, txSpeed: false, rxDimming: true, rxSpeed: false, rxTemp: false, id: 19, name: "Plantgrowth"}  //
	SceneSpring       = Scene{txDimming: true, txSpeed: true, rxDimming: true, rxSpeed: true, rxTemp: false, id: 20, name: "Spring"}         //
	SceneSummer       = Scene{txDimming: true, txSpeed: true, rxDimming: true, rxSpeed: true, rxTemp: false, id: 21, name: "Summer"}         //
	SceneFall         = Scene{txDimming: true, txSpeed: true, rxDimming: true, rxSpeed: true, rxTemp: false, id: 22, name: "Fall"}           //
	SceneDeepdive     = Scene{txDimming: true, txSpeed: true, rxDimming: true, rxSpeed: true, rxTemp: false, id: 23, name: "Deepdive"}       //
	SceneJungle       = Scene{txDimming: true, txSpeed: true, rxDimming: true, rxSpeed: true, rxTemp: false, id: 24, name: "Jungle"}         //
	SceneMojito       = Scene{txDimming: true, txSpeed: true, rxDimming: true, rxSpeed: true, rxTemp: false, id: 25, name: "Mojito"}         //
	SceneClub         = Scene{txDimming: true, txSpeed: true, rxDimming: true, rxSpeed: true, rxTemp: false, id: 26, name: "Club"}           //
	SceneChristmas    = Scene{txDimming: true, txSpeed: true, rxDimming: true, rxSpeed: true, rxTemp: false, id: 27, name: "Christmas"}      //
	SceneHalloween    = Scene{txDimming: true, txSpeed: true, rxDimming: true, rxSpeed: true, rxTemp: false, id: 28, name: "Halloween"}      //
	SceneCandlelight  = Scene{txDimming: true, txSpeed: true, rxDimming: true, rxSpeed: false, rxTemp: false, id: 29, name: "Candlelight"}   //
	SceneGoldenWhite  = Scene{txDimming: true, txSpeed: true, rxDimming: true, rxSpeed: true, rxTemp: false, id: 30, name: "GoldenWhite"}    //
	ScenePulse        = Scene{txDimming: true, txSpeed: true, rxDimming: true, rxSpeed: true, rxTemp: false, id: 31, name: "Pulse"}          //
	SceneSteampunk    = Scene{txDimming: true, txSpeed: true, rxDimming: true, rxSpeed: true, rxTemp: false, id: 32, name: "Steampunk"}      //
)

// ScenesList contains a list of all available scenes.
// The ID doesn't necessarily correspond with the position of the item.
var ScenesList = []Scene{
	SceneOcean,
	SceneRomance,
	SceneSunset,
	SceneParty,
	SceneFireplace,
	SceneCozy,
	SceneForest,
	ScenePastelColors,
	SceneWakeUp,
	SceneBedtime,
	SceneWarmWhite,
	SceneDaylight,
	SceneCoolWhite,
	SceneNightLight,
	SceneFocus,
	SceneRelax,
	SceneTrueColors,
	SceneTVTime,
	ScenePlantgrowth,
	SceneSpring,
	SceneSummer,
	SceneFall,
	SceneDeepdive,
	SceneJungle,
	SceneMojito,
	SceneClub,
	SceneChristmas,
	SceneHalloween,
	SceneCandlelight,
	SceneGoldenWhite,
	ScenePulse,
	SceneSteampunk,
}

// ScenesMap maps scene IDs to Scene objects.
var ScenesMap = scenesMapFromList(ScenesList)

// scenesMapFromList returns a map of scenes from the given list of scenes.
func scenesMapFromList(scenesList []Scene) map[uint]Scene {
	res := map[uint]Scene{}

	for _, scene := range scenesList {
		res[scene.id] = scene
	}

	return res
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
	if scene, ok := ScenesMap[id]; ok {
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

// Name returns the name of the scene.
func (s Scene) Name() string {
	return s.name
}

// ID returns the ID of the scene.
func (s Scene) ID() uint {
	return s.id
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

// HasTemp returns true if the device sends a color temperature along with the scene in the pilot.
func (s Scene) HasTemp() bool {
	return s.rxTemp
}
