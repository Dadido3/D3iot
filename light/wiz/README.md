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

To connect to your device via its IP, use

``` go
light := wiz.NewLight("192.168.1.123:38899")
```

where you have to replace `192.168.1.123` with the corresponding IP of your device.

Alternatively, you can use the DHCP name of the device to connect. As example use

``` go
light := wiz.NewLight("wiz-123abc:38899")
```

where you have to replace `123abc` with the 6 last digits of the device's MAC-Address.

### Read device information

``` go
info, err := light.GetDeviceInfo()
fmt.Printf("%#v", info)
```

### Set the pilot

The pilot contains all the settings that are related to the device's light output.
This is the "raw" way to control the device, no color transformation is done.
You can either use predefined scenes, a color temperature, or RGBW values.

``` go
pilot := wiz.NewPilotWithScene(wiz.SceneClub, 50, 100) // 50 % dimming value, 100 % speed value.
err := light.SetPilot(pilot)
```

``` go
pilot := wiz.NewPilotWithTemp(4000, 80) // Color temperature of 4000 K value, 80 % dimming value.
err := light.SetPilot(pilot)
```

``` go
pilot := wiz.NewPilotWithRGBW(100, 50, 0, 0, 0, 0) // 100 % dimming value and the R, G, B, cold white, warm white values.
err := light.SetPilot(pilot)
```

### Get the pilot

To query the currently active pilot and retrieve the RGBW values, use

``` go
pilot, err := light.GetPilot()
fmt.Printf("%v", pilot)

if pilot.HasRGBW() {
    r, g, b, cw, ww := *pilot.R, *pilot.G, *pilot.B, *pilot.CW, *pilot.WW
}
```

### Pulse

If you have multiple lamps and need to identify a specific device, you can make the lamp change its light output for a given amount of time.

``` go
light.Pulse(50, 100*time.Millisecond)
```
