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

## Devices

There are the following device classes:

1. `RGBTW` - have Red, Green, Blue, Cool White and Warm White LEDs. These type of bulbs can support all the light modes provided by WiZ System.
2. `TW` - have Cool White and Warm White LEDs. Such devices support most static light modes + CCT control.
3. `DW` - have only Dimmable white LEDs. Such devices support only dimming, and some light modes.

The following is a list of known devices by their `ModuleName`.
The list is not complete and may contain mistakes.

| ModuleName | Device class | Model ID | Example products |
| --- | --- | --- | --- |
| `ESP01_SHDW_01`      | `DW`    |          | `WiZ A60 B22 WiZ60 DW` |
| `ESP01_SHDW1_31`     | `DW`    |          | |
| `ESP01_SHRGB_03`     | `RGBTW` |          | |
| `ESP01_SHRGB1C_31`   | `RGBTW` | `23007`  | `Philips 555623 recessed`, `Philips 556167 A19 Frosted Full Colour and Tunable White` |
| `ESP01_SHTW1C_31`    | `TW`    |          | `Philips 555599 Tunable White 5/6 in. LED 65W recessed light`, `WiZ Connected Tunable White Wi-Fi LED (A19)` |
| `ESP03_SHRGB1C_01`   | `RGBTW` | `B23065` | `Philips Color &. Tunable-White A19`, `WiZ A60 E27 EAN 8718699787059`, `WiZ G95 E27 EAN 8718699786359` |
| `ESP03_SHRGB1W_01`   | `RGBTW` | `B27285` | `WiZ WI-FI BLE 100W A67 E27 922-65 RGB EAN 8718699786199`, `Philips Color &. Tunable-White A21` |
| `ESP03_SHRGB3_01ABI` | `RGBTW` |          | |
| `ESP03_SHRGBP_31`    | `RGBTW` |          | `Trio Leuchten WiZ LED` |
| `ESP05_SHDW_01`      | `DW`    |          | |
| `ESP05_SHRGBL_21`    | `RGBTW` |          | `WiZ WI-FI Color A19` |
| `ESP06_SHDW1_01`     | `DW`    |          | |
| `ESP06_SHDW9_01`     | `DW`    |          | `Filament amber A19 E26` |
| `ESP14_SHRGB1C_01`   | `RGBTW` |          | |
| `ESP15_SHTW1_01I`    | `TW`    |          | |
| `ESP17_SHTW9_01`     | `TW`    |          | `WiZ Filament Bulb EAN 8718699786793` |
| `ESP56_SHTW3_01`     | `TW`    |          | `WiZ G25 Filament bulb`, `WiZ G95 E27 720lm Filament Bulb` |

## WiZ documentation

There is no real documentation of WiZ devices, but there is the [WiZ Pro API documentation](https://docs.pro.wizconnected.com) that shares a lot of information with the API that this library uses to directly communicate with the devices.
