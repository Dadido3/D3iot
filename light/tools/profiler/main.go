// Copyright (c) 2021-2022 David Vogel
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
		/*{dcs: emission.LinDCSVector{0, 0, 0, 1, 1}, xyz: emission.CIE1976LAB{L: 419.048502, A: 41.67841, B: 80.909031, WhitePoint: emission.StandardIlluminantD50}.CIE1931XYZRel()},
		{dcs: emission.LinDCSVector{0, 0, 0, 0, 1}, xyz: emission.CIE1976LAB{L: 329.770063, A: 86.108957, B: 181.467988, WhitePoint: emission.StandardIlluminantD50}.CIE1931XYZRel()},
		{dcs: emission.LinDCSVector{0, 0, 0, 1, 0}, xyz: emission.CIE1976LAB{L: 327.582717, A: -25.147464, B: -17.477994, WhitePoint: emission.StandardIlluminantD50}.CIE1931XYZRel()},
		{dcs: emission.LinDCSVector{0, 0, 1, 0, 0}, xyz: emission.CIE1976LAB{L: 122.436975, A: 203.136262, B: -361.854121, WhitePoint: emission.StandardIlluminantD50}.CIE1931XYZRel()},
		{dcs: emission.LinDCSVector{0, 1, 0, 0, 0}, xyz: emission.CIE1976LAB{L: 200.53485, A: -346.7991, B: 158.705295, WhitePoint: emission.StandardIlluminantD50}.CIE1931XYZRel()},
		{dcs: emission.LinDCSVector{1, 0, 0, 0, 0}, xyz: emission.CIE1976LAB{L: 144.04654, A: 233.651276, B: 244.636863, WhitePoint: emission.StandardIlluminantD50}.CIE1931XYZRel()},
		{dcs: emission.LinDCSVector{0.5, 0.5, 0.5, 0, 0}, xyz: emission.CIE1976LAB{L: 187.845395, A: 21.465868, B: -125.669222, WhitePoint: emission.StandardIlluminantD50}.CIE1931XYZRel()},
		{dcs: emission.LinDCSVector{1, 0, 0, 0.5, 0}, xyz: emission.CIE1976LAB{L: 274.968486, A: 75.372955, B: 12.030361, WhitePoint: emission.StandardIlluminantD50}.CIE1931XYZRel()},
		{dcs: emission.LinDCSVector{0, 1, 0, 0.5, 0}, xyz: emission.CIE1976LAB{L: 297.018857, A: -138.013113, B: 36.911707, WhitePoint: emission.StandardIlluminantD50}.CIE1931XYZRel()},
		{dcs: emission.LinDCSVector{0, 0, 1, 0.5, 0}, xyz: emission.CIE1976LAB{L: 269.303661, A: 47.263185, B: -200.430022, WhitePoint: emission.StandardIlluminantD50}.CIE1931XYZRel()},
		{dcs: emission.LinDCSVector{1, 0, 0, 0, 0.5}, xyz: emission.CIE1976LAB{L: 276.6268, A: 142.680486, B: 173.870282, WhitePoint: emission.StandardIlluminantD50}.CIE1931XYZRel()},*/
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
			WhitePoint: emission.StandardIlluminantD50,
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
		labColor := xyzColor.Relative(whitePointLumen / whitePoint.Y).CIE1976LAB(emission.StandardIlluminantD50)

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
		data.WhitePoint = emission.StandardIlluminantD50

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
