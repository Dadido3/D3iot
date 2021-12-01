# WiZ light devices

This package contains everything you need to query and control WiZ light bulbs.

## Usage

Install via

``` shell
go get github.com/Dadido3/D3iot/light/wiz
```

and import

``` go
import "github.com/Dadido3/D3iot/light/wiz"
```

### Connecting

First you have add your light bulb to your network, which for now can only be done via the WiZ app over bluetooth.
Once the device has been connected to your Wi-Fi, you can use the app to retrieve the light's IP.

Finally, to connect to your device, use

``` go
light := wiz.NewLight("192.168.1.123:38899")
```

where you have to replace `192.168.1.123` with the corresponding IP of your device.

### Read device information

``` go
info, err := light.GetDeviceInfo()
fmt.Printf("%#v", info)
```

### Set the pilot

The pilot contains all the settings that are related to the device's light output.
This is the "raw" way to control the device, no color transformation is done.
You can either use predefined scenes, or RGBW values.

``` go
pilot := wiz.Pilot{}.WithScene(wiz.SceneClub, 6000, 50, 50)
err := light.SetPilot(pilot)
```

``` go
pilot := wiz.Pilot{}.WithRGBW(50, 0, 0, 0, 0, 100) // R, G, B, cold white, warm white, dimming value
err := light.SetPilot(pilot)
```

### Get the pilot

To query the currently active pilot and retrieve the RGBW values, use

``` go
pilot, err := light.GetPilot()

if pilot.HasRGBW() {
    r, g, b, cw, ww := *pilot.R, *pilot.G, *pilot.B, *pilot.CW, *pilot.WW
}
```

### Pulse

If you have multiple lamps and need to identify a specific device, you can make the lamp change its light output for a given amount of time.

``` go
light.Pulse(50, 100*time.Millisecond)
```
