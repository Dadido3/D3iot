// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package wiz

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

// DevInfo contains the bulb's device information.
// It's basically a stripped down version of SystemConfig.
type DevInfo struct {
	Mac        string `json:"mac"`
	DevMac     string `json:"devMac"`
	ModuleName string `json:"moduleName"`
}

// Favs contains the bulb's favorite settings.
// There can be 4 favorites.
type Favs struct {
	Favs [4][7]int `json:"favs"` // The first entry of each favorite entry is the scene ID.
	Opts [4][]int  `json:"opts"` // No idea.
}

type SystemConfig struct {
	Mac         string `json:"mac"`
	HomeID      uint   `json:"homeId"`
	RoomID      uint   `json:"roomId"`
	Rgn         string `json:"rgn"` // Region, e.g. "eu".
	ModuleName  string `json:"moduleName"`
	FWVersion   string `json:"fwVersion"`
	GroupID     uint   `json:"groupId"`
	TypeID      uint   `json:"typeId"`
	HomeLock    bool   `json:"homeLock"`
	PairingLock bool   `json:"pairingLock"`
	DrvConf     []int  `json:"drvConf"`
	Ping        uint   `json:"ping"`
}

type UserConfig struct {
	FadeIn     uint `json:"fadeIn"`     // Fade-in time in milliseconds.
	FadeOut    uint `json:"fadeOut"`    // Fade-out time in milliseconds.
	DFTDim     uint `json:"dftDim"`     // Not sure. Default dimming value in percent?
	OpMode     int  `json:"opMode"`     // No idea.
	PO         bool `json:"po"`         // No idea.
	MinDimming uint `json:"minDimming"` // Minimal dimming value in percent.
	TapSensor  int  `json:"tapSensor"`  // Not sure. Amount of tap sensors?
}

// method represents a query method.
type method string

const (
	// Temporarily affects the bulb.

	methodPulse  method = "pulse"
	methodReboot method = "reboot"

	// These methods change parameters.

	methodRegistration    method = "registration"
	methodReset           method = "reset"
	methodSetDevInfo      method = "setDevInfo"
	methodSetFavs         method = "setFavs"
	methodSetPilot        method = "setPilot"
	methodSetSchd         method = "setSchd"
	methodSetSchdPset     method = "setSchdPset"
	methodSetState        method = "setState"
	methodSetSystemConfig method = "setSystemConfig"
	methodSetUserConfig   method = "setUserConfig"
	//methodSetWifiConfig   method = "setWifiConfig"

	// Methods that retrieve parameters.

	methodGetDevInfo      method = "getDevInfo"
	methodGetFavs         method = "getFavs"
	methodGetPilot        method = "getPilot"
	methodGetSystemConfig method = "getSystemConfig"
	methodGetUserConfig   method = "getUserConfig"
	//methodGetWifiConfig   method = "getWifiConfig"

	// Sync stuff.

	/*methodSyncAlarm          method = "syncAlarm"
	methodSyncBroadcastPilot method = "syncBroadcastPilot"
	methodSyncConfig         method = "syncConfig"
	methodSyncSchdPset       method = "syncSchdPset"
	methodSyncSystemConfig   method = "syncSystemConfig"
	methodSyncUserConfig     method = "syncUserConfig"*/
)

// query is the data structure that holds any query.
type query struct {
	Method method      `json:"method"`
	Env    string      `json:"env,omitempty"`    // No idea. Is "pro" by default. Production, Professional, Prolapse?
	ID     uint        `json:"id,omitempty"`     // No idea. Unique ID for the light bulb, or for this query?
	Params interface{} `json:"params,omitempty"` // The parameters to transmit.
}

// response represents a general response with status code, error message and a result field.
type response struct {
	Method method      `json:"method"`
	Env    string      `json:"env,omitempty"`
	Result interface{} `json:"result,omitempty"`

	Error *struct {
		Code    int64  `json:"code"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// Check returns an error in case the response contains any error message or status code.
func (r response) Check(m method) error {
	if r.Method == "" {
		return fmt.Errorf("empty response method")
	}
	if r.Method != m {
		return fmt.Errorf("response is for different method. Got %q, want %q", r.Method, m)
	}
	if r.Error != nil {
		return fmt.Errorf("light bulb returned error %d: %v", r.Error.Code, r.Error.Message)
	}
	/*if !r.Result.Success {
		return fmt.Errorf("light bulb signalled that the operation failed")
	}*/
	return nil
}

// Pulse lets the light bulb do a single pulse of the given delta for the given duration.
// This can be used to identify a specific bulb.
func (l *Light) Pulse(delta int, duration time.Duration) error {
	q := query{
		Method: methodPulse,
		Env:    "pro",
		Params: struct {
			Delta    int `json:"delta"`
			Duration int `json:"duration"`
		}{
			Delta:    delta,
			Duration: int(duration.Milliseconds()),
		},
	}

	var r response
	l.jsonQuery(q, &r)

	// Check if the response contains any error code.
	return r.Check(q.Method)
}

// Reboot reboots the bulb.
// This will not reset any parameters.
func (l *Light) Reboot() error {
	q := query{
		Method: methodReboot,
		Env:    "pro",
	}

	var r response
	l.jsonQuery(q, &r)

	// Check if the response contains any error code.
	return r.Check(q.Method)
}

// SetPilot sends the given pilot to the light bulb.
func (l *Light) SetPilot(p Pilot) error {
	q := query{
		Method: methodSetPilot,
		Env:    "pro",
		Params: p,
	}

	var r response
	l.jsonQuery(q, &r)

	return r.Check(q.Method)
}

// GetDeviceInfo queries the bulb for its device info.
func (l *Light) GetDeviceInfo() (DevInfo, error) {
	q := query{
		Method: methodGetDevInfo,
		Env:    "pro",
	}

	result := DevInfo{}

	var r response
	r.Result = &result
	l.jsonQuery(q, &r)

	return result, r.Check(q.Method) // This may return data in case of an error.
}

// GetFavs queries the bulb for its favorites/presets.
func (l *Light) GetFavs() (Favs, error) {
	q := query{
		Method: methodGetFavs,
		Env:    "pro",
	}

	result := Favs{}

	var r response
	r.Result = &result
	l.jsonQuery(q, &r)

	return result, r.Check(q.Method) // This may return data in case of an error.
}

// GetPilot queries the bulb for its current pilot data.
func (l *Light) GetPilot() (Pilot, error) {
	q := query{
		Method: methodGetPilot,
		Env:    "pro",
	}

	result := Pilot{}

	var r response
	r.Result = &result
	l.jsonQuery(q, &r)

	return result, r.Check(q.Method) // This may return data in case of an error.
}

// GetSystemConfig queries the bulb for its system configuration.
func (l *Light) GetSystemConfig() (SystemConfig, error) {
	q := query{
		Method: methodGetSystemConfig,
		Env:    "pro",
	}

	result := SystemConfig{}

	var r response
	r.Result = &result
	l.jsonQuery(q, &r)

	return result, r.Check(q.Method) // This may return data in case of an error.
}

// GetUserConfig queries the bulb for its user configuration.
func (l *Light) GetUserConfig() (UserConfig, error) {
	q := query{
		Method: methodGetUserConfig,
		Env:    "pro",
	}

	result := UserConfig{}

	var r response
	r.Result = &result
	l.jsonQuery(q, &r)

	return result, r.Check(q.Method) // This may return data in case of an error.
}

// GetWifiConfig queries the bulb for its user configuration.
/*func (l *Light) GetWifiConfig() (WifiConfig, error) {
	q := query{
		Method: methodGetWifiConfig,
		Env:    "pro",
	}

	result := WifiConfig{}

	var r response
	r.Result = &result
	l.jsonQuery(q, &r)

	return result, r.Check(q.Method) // This may return data in case of an error.
}*/

// jsonQuery sends the given query structure as JSON, and unmarshals the JSON response into the given structure r.
func (l *Light) jsonQuery(q query, r interface{}) error {
	data, err := json.Marshal(q)
	if err != nil {
		return err
	}

	//log.Printf("%q query: %q", q.Method, string(data))

	responseData, err := l.rawSend(data)
	if err != nil {
		return err
	}

	//log.Printf("%q response: %q", q.Method, string(responseData))

	if err := json.Unmarshal(responseData, &r); err != nil {
		return err
	}

	return nil
}

// rawSend sends the given data to the light bulb via UDP.
// The answer given by the bulb will be returned as byte slice.
//
// This assumes that there is only a single connection between the local and remote address.
// If there is more communication going on, the response might be something unexpected.
func (l *Light) rawSend(data []byte) ([]byte, error) {
	conn, err := net.Dial("udp", l.address)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
	if _, err := conn.Write(data); err != nil {
		return nil, err
	}

	buf := make([]byte, 1024)

	conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}

	return buf[:n], nil
}
