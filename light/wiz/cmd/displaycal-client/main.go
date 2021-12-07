// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Dadido3/D3iot/light/wiz"
)

var flagDeviceAddress = flag.String("address", "", "The address of the device to be controlled. Example: \"--address wiz-123abc:38899\" or \"--address 192.168.1.123:38899\"")
var flagDisplayCALAddress = flag.String("http-server", "http://localhost:8080/", "The address of the DisplayCAL web-server. Example: \"--http-server http://localhost:8080/\"")
var flagMaxR = flag.Uint("max-r", 255, "The maximum value of the red channel. The input will be scaled to fit into [0, max].")
var flagMaxG = flag.Uint("max-g", 255, "The maximum value of the green channel. The input will be scaled to fit into [0, max].")
var flagMaxB = flag.Uint("max-b", 255, "The maximum value of the blue channel. The input will be scaled to fit into [0, max].")

func main() {
	flag.Parse()

	if *flagDeviceAddress == "" {
		log.Printf("No device address given. Start program with the \"--address\" parameter set.")
		log.Printf("Example: displaycal-client --address wiz-123abc:38899")
		return
	}

	light := wiz.NewLight(*flagDeviceAddress)

	// Setup HTTP request.
	req, err := http.NewRequest("GET", *flagDisplayCALAddress, nil)
	if err != nil {
		log.Panicf("http.NewRequest() failed: %v", err)
	}
	req.URL.Path = "/ajax/messages"
	log.Printf("Using %q to query messages.", req.URL)

	// Request colors in a loop.
	for {
		response, err := queryMessage(req)
		if err != nil {
			log.Panicf("queryMessage() failed: %v", err)
		}

		var r, g, b uint
		if _, err := fmt.Sscanf(response, "#%02x%02x%02x", &r, &g, &b); err != nil {
			log.Panicf("failed to parse %q as color: %v", response, err)
		}

		rScaled, gScaled, bScaled := uint8(r*(*flagMaxR)/255), uint8(g*(*flagMaxG)/255), uint8(b*(*flagMaxB)/255)

		if err := light.SetPilot(wiz.NewPilotWithRGB(100, rScaled, gScaled, bScaled)); err != nil {
			log.Panicf("light.SetPilot() with RGB %d, %d, %d failed: %v", rScaled, gScaled, bScaled, err)
		}
	}
}

func queryMessage(request *http.Request) (response string, err error) {
	client := new(http.Client)

	resp, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("client.Do() failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ioutil.ReadAll() failed: %w", err)
	}

	return string(respBody), nil
}
