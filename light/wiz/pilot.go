// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package wiz

import (
	"encoding/json"
	"fmt"
)

// Pilot contains the parameters regarding light output.
//
// A valid pilot can either contain a scene, a color value defined by temperature or a color defined by RGB(W) values.
// There must be only one at a time.
// An exception is that the device sends a scene ID of 0 back when no scene is defined.
type Pilot struct {
	//Cnx  string `json:"cnx,omitempty"`
	Mac  string `json:"mac,omitempty"`  // Mac address.
	RSSI int    `json:"rssi,omitempty"` // Signal strength.
	Src  string `json:"src,omitempty"`  // No idea.

	State bool `json:"state"` // On off state.

	Dimming    *uint  `json:"dimming,omitempty"`    // Dimming value in percent. There is a limit as to how low the dimming value can go.
	Speed      *uint  `json:"speed,omitempty"`      // Color changing speed in percent. This only seems to influence the scene playback speed.
	Temp       *uint  `json:"temp,omitempty"`       // Color temperature in Kelvin.
	SchdPsetID *uint  `json:"schdPsetId,omitempty"` // Not sure. Scheduled preset?
	Scene      *Scene `json:"sceneId,omitempty"`    // The scene ID.

	R  *uint8 `json:"r,omitempty"` // Red luminance in range 0-255.
	G  *uint8 `json:"g,omitempty"` // Green luminance range 0-255.
	B  *uint8 `json:"b,omitempty"` // Blue luminance range 0-255.
	CW *uint8 `json:"c,omitempty"` // Cold white luminance range 0-255.
	WW *uint8 `json:"w,omitempty"` // Warm white luminance range 0-255.
}

// NewPilot returns a pilot with the given light state.
func NewPilot(state bool) Pilot {
	return Pilot{
		State: state,
	}
}

// NewPilotWithRGB returns a pilot with the given RGB values.
// No color transformation is done, all values are passed directly to the device.
func NewPilotWithRGB(dimming uint, r, g, b uint8) Pilot {
	// Edge case, if all values are 0, set the state to false.
	// Otherwise the lamp will ignore the pilot.
	var state bool
	if r > 0 || g > 0 || b > 0 {
		state = true
	} else {
		state = false
	}

	return Pilot{
		State:   state,
		Dimming: &dimming,
		R:       &r,
		G:       &g,
		B:       &b,
	}
}

// NewPilotWithRGBW returns a pilot with the given RGBW values.
// No color transformation is done, all values are passed directly to the device.
func NewPilotWithRGBW(dimming uint, r, g, b, cw, ww uint8) Pilot {
	// Edge case, if all values are 0, set the state to false.
	// Otherwise the lamp will ignore the pilot.
	var state bool
	if r > 0 || g > 0 || b > 0 || cw > 0 || ww > 0 {
		state = true
	} else {
		state = false
	}

	return Pilot{
		State:   state,
		Dimming: &dimming,
		R:       &r,
		G:       &g,
		B:       &b,
		CW:      &cw,
		WW:      &ww,
	}
}

// NewPilotWithWhite returns a pilot with the given cold and warm white values.
// No color transformation is done, all values are passed directly to the device.
func NewPilotWithWhite(dimming uint, cw, ww uint8) Pilot {
	// Edge case, if all values are 0, set the state to false.
	// Otherwise the lamp will ignore the pilot.
	var state bool
	if cw > 0 || ww > 0 {
		state = true
	} else {
		state = false
	}

	return Pilot{
		State:   state,
		Dimming: &dimming,
		CW:      &cw,
		WW:      &ww,
	}
}

// NewPilotWithTemp returns a pilot with the given color temperature values.
// No color transformation is done, all values are passed directly to the device.
func NewPilotWithTemp(dimming, temperature uint) Pilot {
	return Pilot{
		State:   true,
		Dimming: &dimming,
		Temp:    &temperature,
	}
}

// NewPilotWithScene returns a pilot with the given scene.
//
// Some scenes need a dimming and/or speed value, see
//
//	s.NeedsDimming()
//	s.NeedsSpeed()
//
// In case a scene doesn't need such value, it is ignored by this function.
func NewPilotWithScene(s Scene, dimming, speed uint) Pilot {
	p := Pilot{
		State: true,
		Scene: &s,
	}

	if s.NeedsDimming() {
		p.Dimming = &dimming
	}
	if s.NeedsSpeed() {
		p.Speed = &speed
	}

	return p
}

// WithDimming returns a copy of the pilot set to the given dimming value.
func (p Pilot) WithDimming(dimming uint) Pilot {
	p.Dimming = &dimming
	return p
}

// WithLightOff returns a copy of the pilot with the light turned off.
func (p Pilot) WithLightOff() Pilot {
	p.State = false
	return p
}

// WithLightOn returns a copy of the pilot with the light turned on.
func (p Pilot) WithLightOn() Pilot {
	p.State = true
	return p
}

// WithLightToggled returns a copy of the pilot with the light toggled.
func (p Pilot) WithLightToggled() Pilot {
	p.State = !p.State
	return p
}

// WithRGB returns a copy of the pilot with the given color values set.
// This will change the on off state.
// This will reset any other competing value like scene ID or temperature.
func (p Pilot) WithRGB(r, g, b uint8) Pilot {
	p.Scene, p.Temp, p.Speed = nil, nil, nil
	p.R, p.G, p.B, p.CW, p.WW = &r, &g, &b, nil, nil
	if r > 0 || g > 0 || b > 0 {
		p.State = true
	} else {
		p.State = false
	}
	return p
}

// WithRGBW returns a copy of the pilot with the given color values set.
// This will change the on off state.
// This will reset any other competing value like scene ID or temperature.
func (p Pilot) WithRGBW(r, g, b, cw, ww uint8) Pilot {
	p.Scene, p.Temp, p.Speed = nil, nil, nil
	p.R, p.G, p.B, p.CW, p.WW = &r, &g, &b, &cw, &ww
	if r > 0 || g > 0 || b > 0 || cw > 0 || ww > 0 {
		p.State = true
	} else {
		p.State = false
	}
	return p
}

// WithWhite returns a copy of the pilot with the given color values set.
// This will change the on off state.
// This will reset any other competing value like scene ID or temperature.
func (p Pilot) WithWhite(cw, ww uint8) Pilot {
	p.Scene, p.Temp, p.Speed = nil, nil, nil
	p.R, p.G, p.B, p.CW, p.WW = nil, nil, nil, &cw, &ww
	if cw > 0 || ww > 0 {
		p.State = true
	} else {
		p.State = false
	}
	return p
}

// WithTemp returns a copy of the pilot with the given color temperature set.
// This will not change the on/off state or dimming value of the pilot.
// This will reset any other competing value like scene ID or RGBW values.
func (p Pilot) WithTemp(temp uint) Pilot {
	p.Scene, p.Temp, p.Speed = nil, &temp, nil
	p.R, p.G, p.B, p.CW, p.WW = nil, nil, nil, nil, nil
	return p
}

// WithScene returns a copy of the pilot with the given scene set.
// This will not change the on/off state of the pilot.
// This will reset any other competing value like RGB values.
func (p Pilot) WithScene(s Scene, speed uint) Pilot {
	p.R, p.G, p.B, p.CW, p.WW, p.Temp = nil, nil, nil, nil, nil, nil
	p.Scene = &s

	if s.NeedsSpeed() {
		p.Speed = &speed
	} else {
		p.Speed = nil
	}

	if !s.NeedsDimming() {
		p.Dimming = nil
	}

	return p
}

// HasRGB returns true, if the pilot contains RGB values, including all Zero values.
// This ensures that you can dereference the R, G, and B fields.
func (p Pilot) HasRGB() bool {
	if p.R != nil && p.G != nil && p.B != nil {
		return true
	}
	return false
}

// HasRGBW returns true, if the pilot contains RGBW values, including all Zero values.
// This ensures that you can dereference the R, G, B, CW, and WW fields.
func (p Pilot) HasRGBW() bool {
	if p.R != nil && p.G != nil && p.B != nil && p.CW != nil && p.WW != nil {
		return true
	}
	return false
}

// HasWhite returns true, if the pilot contains cold and warm white values, including all Zero values.
// This ensures that you can dereference the CW and WW fields.
func (p Pilot) HasWhite() bool {
	if p.CW != nil && p.WW != nil {
		return true
	}
	return false
}

// HasScene returns true, if the pilot contains a scene value.
// This ensures that you can dereference the Scene field.
func (p Pilot) HasScene() bool {
	return p.Scene != nil
}

// HasTemp returns true, if the pilot contains a color temperature value.
// This ensures that you can dereference the Temp field.
func (p Pilot) HasTemp() bool {
	return p.Temp != nil
}

// HasSpeed returns true, if the pilot contains a speed value.
// This ensures that you can dereference the Speed field.
func (p Pilot) HasSpeed() bool {
	return p.Speed != nil
}

// HasDimming returns true, if the pilot contains a dimming value.
// This ensures that you can dereference the Dimming field.
func (p Pilot) HasDimming() bool {
	return p.Dimming != nil
}

// UnmarshalJSON implements the JSON unmarshaler interface.
func (p *Pilot) UnmarshalJSON(data []byte) error {
	type tempType Pilot
	var temp tempType
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Special case: If the scene ID is 0, remove the scene.
	// The device sends that scene ID, instead of leaving the scene field empty.
	if temp.Scene != nil && temp.Scene.id == 0 {
		temp.Scene = nil
	}

	*p = Pilot(temp)

	return nil
}

func (p Pilot) String() string {
	result := fmt.Sprintf("{State: %v", p.State)

	if p.Dimming != nil {
		result += fmt.Sprintf(", Dimming: %d %%", *p.Dimming)
	}
	if p.Speed != nil {
		result += fmt.Sprintf(", Speed: %d %%", *p.Speed)
	}
	if p.Temp != nil {
		result += fmt.Sprintf(", Temp: %d K", *p.Temp)
	}
	if p.SchdPsetID != nil {
		result += fmt.Sprintf(", SchdPsetID: %v", *p.SchdPsetID)
	}
	if p.Scene != nil {
		result += fmt.Sprintf(", Scene: %v", *p.Scene)
	}

	if p.R != nil {
		result += fmt.Sprintf(", R: %d", *p.R)
	}
	if p.G != nil {
		result += fmt.Sprintf(", G: %d", *p.G)
	}
	if p.B != nil {
		result += fmt.Sprintf(", B: %d", *p.B)
	}
	if p.CW != nil {
		result += fmt.Sprintf(", CW: %d", *p.CW)
	}
	if p.WW != nil {
		result += fmt.Sprintf(", WW: %d", *p.WW)
	}

	result += "}"

	return result
}
