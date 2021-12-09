// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"embed"
	"encoding/json"
	"flag"
	"io/fs"
	"log"
	"net/http"

	"github.com/Dadido3/D3iot/light/wiz"
)

//go:embed static
var staticFiles embed.FS

var flagDeviceAddress = flag.String("address", "wiz-d47cf3:38899", "The address of the device to be controlled. Example: \"--address wiz-123abc:38899\" or \"--address 192.168.1.123:38899\"")
var flagServerBind = flag.String("server-bind", ":8081", "The server binding. Example: \"--server-bind :8081\"")

var matches = map[RGBWValue]CIE1931XYZColor{
	/*{255, 0, 0, 0, 13}:    {0.339, 0.178, 0.020},
	{0, 208, 0, 12, 23}:   {0.196, 0.354, 0.080},
	{0, 0, 255, 20, 0}:    {0.157, 0.118, 0.616},
	{0, 0, 0, 255, 0}:     {0.795, 0.883, 0.771},
	{0, 0, 0, 0, 140}:     {0.549, 0.485, 0.161},
	{0, 0, 0, 50, 50}:     {0.384, 0.374, 0.228},
	{226, 155, 43, 0, 0}:  {0.267, 0.254, 0.146},
	{100, 100, 100, 0, 0}: {0.173, 0.170, 0.288},
	{60, 150, 100, 0, 0}:  {0.136, 0.188, 0.285},
	{40, 40, 40, 40, 40}:  {0.361, 0.353, 0.290},*/
}

func main() {
	// Debug optimization.
	/*primaries, ssr, err := CalculatePrimaries(matches)
	if err != nil {
		log.Printf("Failed to generate primaries: %v", err)
		return
	}
	log.Printf("Optimized primaries: %v", primaries)
	log.Printf("SSR is %f.", ssr)
	log.Print(RGBWValue{0, 0, 0, 255, 255}.CIE1931XYZColor(primaries))*/

	flag.Parse()

	if *flagDeviceAddress == "" {
		log.Printf("No device address given. Start program with the \"--address\" parameter set.")
		log.Printf("Example: profiling --address wiz-123abc:38899")
		return
	}

	light := wiz.NewLight(*flagDeviceAddress)

	staticFilesSub, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", http.FileServer(http.FS(staticFilesSub)))

	http.HandleFunc("/api/setRGBW", func(w http.ResponseWriter, r *http.Request) {
		var rgbw RGBWValue

		if err := json.NewDecoder(r.Body).Decode(&rgbw); err != nil {
			log.Printf("Couldn't decode JSON: %v", err)
			return
		}

		light.SetPilot(wiz.NewPilotWithRGBW(100, rgbw.R, rgbw.G, rgbw.B, rgbw.CW, rgbw.WW))
	})

	http.HandleFunc("/api/addDataPoint", func(w http.ResponseWriter, r *http.Request) {
		var data struct {
			RGBWValue
			CIE1931XYZColor
		}

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			log.Printf("Couldn't decode JSON: %v", err)
			return
		}

		log.Printf("Adding %v to matches list.", data)
		matches[data.RGBWValue] = data.CIE1931XYZColor

		primaries, ssr, err := CalculatePrimaries(matches)
		if err != nil {
			log.Printf("Failed to generate primaries: %v", err)
			return
		}

		log.Printf("Optimized primaries: %v", primaries)
		log.Printf("SSR is %f.", ssr)
	})

	log.Printf("Server is listening at %q. Connect to \"http://localhost:8081\", or whatever port you specified.", *flagServerBind) // TODO: Show correct port.

	log.Fatal(http.ListenAndServe(*flagServerBind, nil))
}
