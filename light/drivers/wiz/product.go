// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package wiz

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Dadido3/D3iot/light/emission"
)

// Product represents a specific WiZ model.
// It may contain abilities, limits, profiling and calibration data.
type Product struct {
	// The ModuleName that the device returns via GetDeviceInfo() or GetSystemConfig().
	moduleName string

	// Describes which LEDs a product contains.
	// This is specific to WiZ devices.
	deviceClass deviceClass

	// Describes the set of light emitting things that together generate a single color impression.
	// This is responsible for color transformations.
	// This must not be nil!
	moduleProfile emission.ModuleProfile

	// The valid color temperatures are described by the interval [MinTemp, MaxTemp].
	// This doesn't necessarily correspond with the range that the white LEDs can output.
	minTemp, maxTemp *uint
}

// ModuleName returns the module name.
// This may not be the exact module name of the device used, but a compatible one.
func (p Product) ModuleName() string {
	return p.moduleName
}

// DimmingCapability returns the min and max dimming value that the product supports.
// If the returned bool is false, the device doesn't have any dimming control.
func (p Product) DimmingCapability() (min, max uint, has bool) {
	return 10, 100, true
}

// TempCapability returns the interval of available color temperatures [min, max] that the product supports.
// If the returned bool is false, the device doesn't have any temperature control.
func (p Product) TempCapability() (min, max uint, has bool) {
	if p.minTemp != nil && p.maxTemp != nil {
		return *p.minTemp, *p.maxTemp, true
	}

	// Fallback to device class temperature capability.
	return p.deviceClass.TempCapability()
}

// RGBWCapability returns which LEDs this device contains and allows control over.
func (p Product) RGBWCapability() (r, g, b, cw, ww bool) {
	return p.deviceClass.RGBWCapability()
}

// ScenesCapability returns a list of scenes that the product supports.
func (p Product) ScenesCapability() []Scene {
	return p.deviceClass.ScenesCapability()
}

// MarshalJSON implements the JSON marshaler interface.
func (p Product) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

func (p Product) String() string {
	result := fmt.Sprintf("wiz.Product{%q, DeviceClass: %q", p.moduleName, p.deviceClass)

	if min, max, hasDimming := p.DimmingCapability(); hasDimming {
		result += fmt.Sprintf(", Dimming: [%d %%, %d %%]", min, max)
	}
	if min, max, hasTemp := p.TempCapability(); hasTemp {
		result += fmt.Sprintf(", ColorTemp: [%d K, %d K]", min, max)
	}
	if r, g, b, cw, ww := p.RGBWCapability(); r || g || b || cw || ww {
		result += fmt.Sprintf(", R:%t, G:%t, B:%t, CW:%t, WW:%t", r, g, b, cw, ww)
	}
	if scenes := p.ScenesCapability(); len(scenes) > 0 {
		result += ", Scenes: {"
		for i, scene := range scenes {
			result += fmt.Sprint(scene)
			if i < len(scenes)-1 {
				result += ", "
			}
		}
		result += "}"
	}

	return result + "}"
}

// deviceClass describes groups/classes of devices that have common abilities.
type deviceClass string

const (
	deviceClassDW    deviceClass = "DW"    // Dimmable white.
	deviceClassTW    deviceClass = "TW"    // Tweakable/Tunable white.
	deviceClassRGBTW deviceClass = "RGBTW" // RGB + cold white + warm white. DeviceClassRGBTW is called "RGB" in the moduleName.
)

// deviceClassDWScenes is a predefined list of scenes that are available to DW class devices.
var deviceClassDWScenes = []Scene{SceneWakeUp, SceneBedtime, SceneCoolWhite, SceneNightLight, SceneCandlelight, SceneGoldenWhite, ScenePulse, SceneSteampunk}

// deviceClassTWScenes is a predefined list of scenes that are available to TW class devices.
var deviceClassTWScenes = []Scene{SceneCozy, SceneWakeUp, SceneBedtime, SceneWarmWhite, SceneDaylight, SceneCoolWhite, SceneNightLight, SceneFocus, SceneRelax, SceneTVTime, SceneCandlelight, SceneGoldenWhite, ScenePulse, SceneSteampunk}

// ScenesCapability returns a list of scenes that the device class supports.
func (dc deviceClass) ScenesCapability() []Scene {
	switch dc {
	case "DW":
		return deviceClassDWScenes
	case "TW":
		return deviceClassTWScenes
	case "RGBTW":
		return ScenesList
	}

	return nil
}

// TempCapability returns the interval of available color temperatures [min, max] that the device class supports.
// If the returned bool is false, the device class allows for no temperature control.
// This capability may be overwritten by a Product, if there is
func (dc deviceClass) TempCapability() (min, max uint, has bool) {
	switch dc {
	case "DW":
		return 0, 0, false
	case "TW", "RGBTW":
		return 2700, 6500, true
	}

	return 0, 0, false
}

// RGBWCapability returns which LEDs this device class contains and allows control over.
func (dc deviceClass) RGBWCapability() (r, g, b, cw, ww bool) {
	switch dc {
	case "DW":
		return false, false, false, false, false
	case "TW":
		return false, false, false, true, true // Assuming you can control CW and WW individually.
	case "RGBTW":
		return true, true, true, true, true
	}

	return false, false, false, false, false
}

func parseDeviceClass(moduleName string) (deviceClass, error) {
	// Split moduleName into moduleFamily, details, and revision.
	splitted := strings.Split(moduleName, "_")
	if len(splitted) != 3 {
		return "", fmt.Errorf("unexpected moduleName format. Got %d sub-strings, want %d", len(splitted), 3)
	}

	details := splitted[1]

	// Assume that all devices are single headed for now.
	switch {
	case strings.HasPrefix(details, "SHDW"):
		return deviceClassDW, nil
	case strings.HasPrefix(details, "SHTW"):
		return deviceClassTW, nil
	case strings.HasPrefix(details, "SHRGB"):
		return deviceClassRGBTW, nil
	}

	return "", fmt.Errorf("%q doesn't match with any known device class", details)
}

// determineProduct returns a matching product for the given moduleName.
// This can be a similar product if there is no exact match.
func determineProduct(moduleName string) (*Product, error) {
	// Find exact match.
	for _, product := range products {
		if product.moduleName == moduleName {
			return &product, nil
		}
	}

	// Not found, match some similar device.

	deviceClass, err := parseDeviceClass(moduleName)
	if err != nil {
		return nil, err
	}

	// Find the first product with the same device class.
	// TODO: Improve how unknown devices are matched
	for _, product := range products {
		if product.deviceClass == deviceClass {
			return &product, nil
		}
	}

	// No device found.
	return nil, fmt.Errorf("couldn't find matching device for moduleName %q", moduleName)
}
