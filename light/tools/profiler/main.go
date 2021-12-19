// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"github.com/Dadido3/D3iot/light"
	"github.com/Dadido3/D3iot/light/drivers/wiz"
	"github.com/Dadido3/D3iot/light/emission"
)

//go:embed static
var staticFiles embed.FS

var flagDeviceWiZ = flag.String("device-wiz", "wiz-d47cf3:38899", "The address of the WiZ device to be profiled. Example: \"--device-wiz wiz-123abc:38899\" or \"--device-wiz 192.168.1.123:38899\".")
var flagServerPort = flag.Int("server-port", 8081, "The server port. Example: \"--server-port 8081\".")
var flagModuleNumber = flag.Int("module", 0, "Number of the module that we want to profile.")
var flagMaxLuminance = flag.Float64("max-luminance", 1521, "The maximum luminance value in lumen that the light can output.")

func main() {
	flag.Parse()

	// Connect to light device.
	var light light.Light
	switch {
	case *flagDeviceWiZ != "":
		var err error
		if light, err = wiz.NewLight(*flagDeviceWiZ); err != nil {
			log.Printf("wiz.NewLight(%q) failed: %v", *flagDeviceWiZ, err)
			return
		}

	default:
		log.Printf("No device address given. Start program with any \"--device-...\" parameter set.")
		log.Printf("Example: profiler --device-wiz wiz-123abc:38899")
		return
	}

	staticFilesSub, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}

	// Check module number for validity.
	if *flagModuleNumber < 0 {
		log.Printf("Module number %d is negative.", *flagModuleNumber)
		return
	}
	if *flagModuleNumber >= light.Modules() {
		log.Printf("Module %d is outside the available modules of the light device (%d).", *flagModuleNumber, light.Modules())
		return
	}

	// Color profile of the module we want to profile.
	colorProfile := light.ColorProfiles()[*flagModuleNumber]

	// List of matches.
	matches := []Match{
		/*{dcs: emission.LinDCSVector{0, 0, 0, 0.5, 0.5}, xyz: emission.CIE1976LAB{81, 8, 34, emission.StandardIlluminantD65}.CIE1931XYZRel()},
		{dcs: emission.LinDCSVector{0, 0, 0, 0.5, 0}, xyz: emission.CIE1976LAB{61.8, -9, 13, emission.StandardIlluminantD65}.CIE1931XYZRel()},
		{dcs: emission.LinDCSVector{0, 0, 0, 0, 0.5}, xyz: emission.CIE1976LAB{60.6, 19, 46, emission.StandardIlluminantD65}.CIE1931XYZRel()},
		{dcs: emission.LinDCSVector{0.8, 0, 0, 0, 0.1}, xyz: emission.CIE1976LAB{41.1, 51, 48, emission.StandardIlluminantD65}.CIE1931XYZRel()},
		{dcs: emission.LinDCSVector{0.8, 0, 0, 0.1, 0}, xyz: emission.CIE1976LAB{40.6, 42, 26, emission.StandardIlluminantD65}.CIE1931XYZRel()},
		{dcs: emission.LinDCSVector{0, 0.6, 0, 0.2, 0}, xyz: emission.CIE1976LAB{52.1, -42, 24, emission.StandardIlluminantD65}.CIE1931XYZRel()},
		{dcs: emission.LinDCSVector{0, 0.6, 0, 0, 0.2}, xyz: emission.CIE1976LAB{52, -23, 49, emission.StandardIlluminantD65}.CIE1931XYZRel()},
		{dcs: emission.LinDCSVector{0, 0, 0.6, 0, 0.2}, xyz: emission.CIE1976LAB{43.5, 26, -24, emission.StandardIlluminantD65}.CIE1931XYZRel()},
		{dcs: emission.LinDCSVector{0, 0, 0.6, 0.2, 0}, xyz: emission.CIE1976LAB{43.3, 10, -33, emission.StandardIlluminantD65}.CIE1931XYZRel()},
		{dcs: emission.LinDCSVector{0.2, 0.2, 0.2, 0.2, 0.2}, xyz: emission.CIE1976LAB{60.1, 8, 16, emission.StandardIlluminantD65}.CIE1931XYZRel()},
		{dcs: emission.LinDCSVector{0.2, 0.2, 0.2, 0, 0}, xyz: emission.CIE1976LAB{27.9, 6, -13, emission.StandardIlluminantD65}.CIE1931XYZRel()},*/
	}

	http.Handle("/", http.FileServer(http.FS(staticFilesSub)))

	http.HandleFunc("/api/getChannels", func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(colorProfile.Channels()); err != nil {
			err = fmt.Errorf("failed to encode JSON: %w", err)
			log.Printf("HTTP error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/api/LAB2sRGB", func(w http.ResponseWriter, r *http.Request) {
		labColor := emission.CIE1976LAB{
			WhitePoint: emission.StandardIlluminantD65,
		}

		if err := json.NewDecoder(r.Body).Decode(&labColor); err != nil {
			err = fmt.Errorf("failed to decode JSON: %w", err)
			log.Printf("HTTP error: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Convert to sRGB.
		var sRGB emission.StandardRGB
		labColor.CIE1931XYZRel().TransformRGB(&sRGB)

		if err := json.NewEncoder(w).Encode(sRGB); err != nil {
			err = fmt.Errorf("failed to encode JSON: %w", err)
			log.Printf("HTTP error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/api/setDCSVector", func(w http.ResponseWriter, r *http.Request) {
		var dcsVector emission.LinDCSVector

		if err := json.NewDecoder(r.Body).Decode(&dcsVector); err != nil {
			err = fmt.Errorf("failed to decode JSON: %w", err)
			log.Printf("HTTP error: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := light.SetColors(dcsVector); err != nil {
			err = fmt.Errorf("failed to set light color to %v: %w", dcsVector, err)
			log.Printf("HTTP error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/api/DCS2LAB", func(w http.ResponseWriter, r *http.Request) {
		var dcsVector emission.LinDCSVector

		if err := json.NewDecoder(r.Body).Decode(&dcsVector); err != nil {
			err = fmt.Errorf("failed to decode JSON: %w", err)
			log.Printf("HTTP error: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		primaries, ssr, err := calculateTransformation(matches)
		if err != nil {
			err = fmt.Errorf("failed to calculate primaries: %w", err)
			log.Printf("HTTP error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Optimized primaries: %v", primaries)
		log.Printf("SSR is %f.", ssr)

		xyzColor, err := primaries.Multiplied(dcsVector)
		if err != nil {
			err = fmt.Errorf("failed to calculate color transformation: %w", err)
			log.Printf("HTTP error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// The first color is the brightest color.
		whitePoint := matches[0].xyz
		whitePointLumen := *flagMaxLuminance
		labColor := xyzColor.Relative(whitePointLumen / whitePoint.Y).CIE1976LAB(emission.StandardIlluminantD65)

		if err := json.NewEncoder(w).Encode(labColor); err != nil {
			err = fmt.Errorf("failed to encode JSON: %w", err)
			log.Printf("HTTP error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/api/addDataPoint", func(w http.ResponseWriter, r *http.Request) {
		var data struct {
			emission.LinDCSVector
			emission.CIE1976LAB
		}
		data.WhitePoint = emission.StandardIlluminantD65

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			err = fmt.Errorf("failed to decode JSON: %w", err)
			log.Printf("HTTP error: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		matches = append(matches, Match{
			dcs: data.LinDCSVector,
			xyz: data.CIE1976LAB.CIE1931XYZRel(),
		})
		log.Printf("Added match pair: DCS vector %v with {L: %v, a: %v, b: %v}", data.LinDCSVector, data.L, data.A, data.B)

		primaries, ssr, err := calculateTransformation(matches)
		if err != nil {
			err = fmt.Errorf("failed to calculate primaries: %w", err)
			log.Printf("HTTP error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Optimized primaries: %v", primaries)
		log.Printf("SSR is %f.", ssr)
	})

	log.Printf("Server is listening at port %d. Connect to \"http://localhost:%d\".", *flagServerPort, *flagServerPort)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *flagServerPort), nil))
}
