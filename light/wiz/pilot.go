// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package wiz

// Pilot contains the parameters regarding light output.
type Pilot struct {
	//Cnx        string  `json:"cnx,omitempty"`
	Mac        string `json:"mac,omitempty"`        // Mac address.
	RSSI       int    `json:"rssi,omitempty"`       // Signal strength.
	Src        string `json:"src,omitempty"`        // No idea.
	Speed      *uint  `json:"speed,omitempty"`      // Color changing speed in percent. This only seems to influence the scene playback speed.
	Temp       *uint  `json:"temp,omitempty"`       // Color temperature in Kelvin.
	SchdPsetID *uint  `json:"schdPsetId,omitempty"` // Not sure. Scheduled preset?
	SceneID    *Scene `json:"sceneId,omitempty"`    // The scene ID.

	State   bool   `json:"state"`       // On off state.
	R       *uint8 `json:"r,omitempty"` // Red luminance in range 0-255.
	G       *uint8 `json:"g,omitempty"` // Green luminance range 0-255.
	B       *uint8 `json:"b,omitempty"` // Blue luminance range 0-255.
	CW      *uint8 `json:"c,omitempty"` // Cold white luminance range 0-255.
	WW      *uint8 `json:"w,omitempty"` // Warm white luminance range 0-255.
	Dimming uint   `json:"dimming"`     // Dimming value in percent.
}

// WithLightOff returns a copy of the pilot with the light turned off.
func (p Pilot) WithLightOff() Pilot {
	p.State = false
	return p
}

// WithRGBW returns a copy of the pilot with the given color values set.
// This will reset any other competing value like scene ID or temperature.
func (p Pilot) WithRGBW(r, g, b, cw, ww uint8, dimming uint) Pilot {
	p.State, p.Dimming = true, dimming
	p.SceneID, p.Temp, p.Speed = nil, nil, nil
	p.R, p.G, p.B, p.CW, p.WW = &r, &g, &b, &cw, &ww
	return p
}

// WithScene returns a copy of the pilot with the given scene set.
// This will reset any other competing value like RGB values.
func (p Pilot) WithScene(s Scene, temp uint, speed uint, dimming uint) Pilot {
	p.State, p.Dimming = true, dimming
	p.R, p.G, p.B, p.CW, p.WW = nil, nil, nil, nil, nil
	p.SceneID, p.Temp, p.Speed = &s, &temp, &speed
	return p
}

// HasRGBW returns true, if the pilot contains RGBW values, including all Zero values.
func (p Pilot) HasRGBW() bool {
	if p.R != nil && p.G != nil && p.B != nil && p.CW != nil && p.WW != nil {
		return true
	}
	return false
}

// HasScene returns true, if the pilot contains a scene value, including SceneNoScene.
func (p Pilot) HasScene() bool {
	if p.SceneID != nil && p.Temp != nil && p.Speed != nil {
		return true
	}
	return false
}
