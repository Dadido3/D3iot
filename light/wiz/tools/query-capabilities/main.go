// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Dadido3/D3iot/light/wiz"
)

// Result contains all important data of the query.
type Result struct {
	ModuleName     string
	MatchedProduct *wiz.Product // The closest matching product.

	CurrentPilot        wiz.Pilot        // Current pilot.
	CurrentFavs         wiz.Favs         // Current Favorite settings.
	CurrentSystemConfig wiz.SystemConfig // Current system configuration.
	CurrentUserConfig   wiz.UserConfig   // Current user configuration.

	SupportedScenes  []wiz.Scene // List of supported scenes.
	MinTemp, MaxTemp *uint       // Temperature range [minTemp, maxTemp] in K.

	DebugEntries []string // Debug output of the library, split by newline characters.
}

var flagDeviceAddress = flag.String("address", "wiz-d47cf3:38899", "The address of the device to be queried. Example: \"--address wiz-123abc:38899\" or \"--address 192.168.1.123:38899\"")

func main() {
	flag.Parse()

	if *flagDeviceAddress == "" {
		log.Printf("No device address given. Start program with the \"--address\" parameter set.")
		log.Printf("Example: query-capabilities --address wiz-123abc:38899")
		return
	}

	var res Result

	light, err := wiz.NewLight(*flagDeviceAddress)
	if err != nil {
		log.Panicf("wiz.NewLight() failed: %v", err)
	}

	// Write debug output into buffer.
	debugBuffer := new(bytes.Buffer)
	light.DebugWriter = debugBuffer

	// Get device information.
	if devInfo, err := light.GetDeviceInfo(); err != nil {
		log.Panicf("light.GetDeviceInfo() failed: %v", err)
	} else {
		res.ModuleName = devInfo.ModuleName
	}

	// Get favorite settings.
	if favs, err := light.GetFavs(); err != nil {
		log.Panicf("light.GetFavs() failed: %v", err)
	} else {
		res.CurrentFavs = favs
	}

	// Get system configuration settings.
	if systemConfig, err := light.GetSystemConfig(); err != nil {
		log.Panicf("light.GetSystemConfig() failed: %v", err)
	} else {
		res.CurrentSystemConfig = systemConfig
	}

	// Get user configuration settings.
	if userConfig, err := light.GetUserConfig(); err != nil {
		log.Panicf("light.GetUserConfig() failed: %v", err)
	} else {
		res.CurrentUserConfig = userConfig
	}

	// Get matched product and its capabilities.
	res.MatchedProduct = light.Product()

	// Get current pilot.
	if pilot, err := light.GetPilot(); err != nil {
		log.Panicf("light.GetPilot() failed: %v", err)
	} else {
		res.CurrentPilot = pilot
	}

	// Create list of available scenes.
	for i := uint64(0); i <= 64; i++ {
		// Get or create scene from ID.
		scene := wiz.Scene{}
		scene.UnmarshalJSON([]byte(strconv.FormatUint(i, 10)))

		if err := light.SetPilot(wiz.NewPilotWithScene(scene, 50, 100)); err != nil {
			// Check if the code returned by the bulb signals that the scene is not supported.
			var e *wiz.ErrQueryFailed
			if errors.As(err, &e) && e.QueryErrorCode() == wiz.QueryErrorCodeInvalidParams {
				// Don't add scene to supported list.
				continue
			} else {
				// Something else caused an error.
				log.Panicf("light.SetPilot() with given scene %s failed: %v", scene, err)
			}
		}

		// Add scene to supported list.
		res.SupportedScenes = append(res.SupportedScenes, scene)
	}

	// Determine temperature range.
	for temp := uint(1000); temp <= 7000; temp += 100 {
		if err := light.SetPilot(wiz.NewPilotWithTemp(50, temp)); err != nil {
			// Check if the code returned by the bulb signals that the parameters are not supported.
			var e *wiz.ErrQueryFailed
			if errors.As(err, &e) && e.QueryErrorCode() == wiz.QueryErrorCodeInvalidParams {
				// Not even able to set the value.
				continue
			} else {
				// Something else caused an error.
				log.Panicf("light.SetPilot() with given temperature %d failed: %v", temp, err)
			}
		}

		// Read pilot and use that to determine minTemp and maxTemp.
		// Some temperatures may be setable, but the device will clip them into its allowed range.
		if pilot, err := light.GetPilot(); err != nil {
			log.Panicf("light.GetPilot() failed: %v", err)
		} else {
			if res.MinTemp == nil || (pilot.Temp != nil && *res.MinTemp > *pilot.Temp) {
				res.MinTemp = pilot.Temp
			}
			if res.MaxTemp == nil || (pilot.Temp != nil && *res.MaxTemp < *pilot.Temp) {
				res.MaxTemp = pilot.Temp
			}
		}
	}

	// Reset light to previous pilot.
	if err := light.SetPilot(res.CurrentPilot); err != nil {
		log.Printf("light.SetPilot() failed: %v", err)
	}

	// Write debug buffer into result.
	res.DebugEntries = strings.Split(debugBuffer.String(), "\n")

	// Write result.
	os.Mkdir("queried", 0755)
	filename := filepath.Join("queried", res.ModuleName+".json")
	if err := writeResult(filename, res); err != nil {
		log.Panicf("Failed to write file %q: %v", filename, err)
	}
}

func writeResult(filename string, res Result) error {
	resultFile, err := os.Create(filename)
	if err != nil {
		return err
	}

	enc := json.NewEncoder(resultFile)
	enc.SetIndent("", "\t")

	if err := enc.Encode(res); err != nil {
		return err
	}

	return resultFile.Close()
}
