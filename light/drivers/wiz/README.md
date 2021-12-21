# WiZ light devices

This package contains everything you need to query and control WiZ light bulbs.

## Features

- Supports setting and getting pilots from WiZ lights.
- Pilots support all scenes, color temperatures, raw RGBW settings, speed settings and dimming settings.
- Library automatically detects properties and abilities of a device. See `light.Product()`.
- Supports color profiles for correct color rendering.
- Implements `light.Light` which provides a simple interface to set/get colors in a fully color managed manner. See the [light package](../../) for more information.

## Usage

Install via

``` shell
go get github.com/Dadido3/D3iot
```

and import

``` go
import "github.com/Dadido3/D3iot/light/drivers/wiz"
```

### Connecting

First you have add your light bulb to your network, which for now can only be done via the WiZ app over bluetooth.
Once the device has been connected to your Wi-Fi, you can use the app to retrieve the light's IP.

To connect to your device via its IP, use

``` go
light, err := wiz.NewLight("192.168.1.123:38899")
```

where you have to replace `192.168.1.123` with the corresponding IP of your device.

Alternatively, you can use the DHCP name of the device to connect. As example use

``` go
light, err := wiz.NewLight("wiz-123abc:38899")
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
The list is not complete and may contain mistakes, if a device is on this list it doesn't mean that it was tested.

| ModuleName | Device class | `wcr` | `nowc` | Model ID | Product | Product Number |
| ---------- | ------------ | --: | ---: | -------- | ------- | -------------- |
| `ESP01_SHDW_01`      | DW    | 20 | 2 |          | WiZ A60 B22 WiZ60 DW |
| `ESP01_SHDW1_31`     | DW    |    |   |          | |
| `ESP01_SHRGB_03`     | RGBTW | 20 | 2 |          | |
| `ESP01_SHRGB1C_31`   | RGBTW |    |   |          | Philips 555623 recessed |
| `ESP01_SHRGB1C_31`   | RGBTW |    |   |          | Philips 556167 A19 Frosted Full Colour and Tunable White |
| `ESP01_SHRGB1C_31`   | RGBTW | 44 | 1 |          | Philips Wiz Full color recessed Downlights |
| `ESP01_SHRGB1C_31`   | RGBTW | 37 | 1 |          | Philips Smart Wi-Fi LED Flood Light, Full Color 65W, BR30 | MPN: 9290022657 |
| `ESP01_SHRGB1C_31`   | RGBTW | 37 | 1 |          | WIZ COLORS BR30 E26 BULB (IZ0087581) |
| `ESP01_SHTW_03`      | TW    |    |   |          | WiZ LED light bulb WZ0126071 E27 11.5 W | EAN: 5420060420566 |
| `ESP01_SHTW1C_31`    | TW    | 20 | 1 |          | Philips 555599 Tunable White 5/6 in. LED 65W recessed light |
| `ESP01_SHTW1C_31`    | TW    | 20 | 1 |          | WiZ Connected Tunable White Wi-Fi LED (A19) |
| `ESP03_SHRGB_31`     | RGBTW | 25 | 1 |          | atom AT4010/WH/WIZ/TR |
| `ESP03_SHRGB1C_01`   | RGBTW |    |   |          | Philips Color &. Tunable-White A19 |
| `ESP03_SHRGB1C_01`   | RGBTW |    |   |          | [WiZ A60 E27](https://www.wizconnected.com/en/consumer/products/8718699787059/) | EAN: 8718699787059 |
| `ESP03_SHRGB1C_01`   | RGBTW | 70 | 1 |          | [WiZ G95 E27](https://www.wizconnected.com/en/consumer/products/8718699786359/) | EAN: 8718699786359 |
| `ESP03_SHRGB1W_01`   | RGBTW |    |   |          | Philips Color &. Tunable-White A21 |
| `ESP03_SHRGB1W_01`   | RGBTW | 20 | 2 | B27285   | [WiZ WI-FI BLE 100W A67 E27 922-65 RGB](https://www.wizconnected.com/en/consumer/products/8718699786199/) | EAN: 8718699786199 |
| `ESP03_SHRGB3_01ABI` | RGBTW |    |   |          | |
| `ESP03_SHRGBP_31`    | RGBTW | 13 | 2 |          | [Trio Leuchten WiZ LED](https://products.trio-lighting.com/de/451850101-2/) | EAN: 4017807407464 |
| `ESP05_SHDW_01`      | DW    | 20 | 1 |          | |
| `ESP05_SHRGBL_21`    | RGBTW | 30 | 1 |          | WiZ WI-FI Color A19 |
| `ESP06_SHDW1_01`     | DW    |    |   |          | Filament amber ST64 E27 | EAN: 8718699787332 |
| `ESP06_SHDW9_01`     | DW    | 20 | 1 |          | [Filament amber A19 E26](https://www.usa.lighting.philips.com/consumer/p/smart-led-filament-amber-a19-e26/046677555528) | EAN: 046677555528 |
| `ESP14_SHRGB1C_01`   | RGBTW | 40 | 1 |          | |
| `ESP15_SHTW1_01I`    | TW    | 20 | 1 |          | WiZ A60 E27 WiZ60 TW F |
| `ESP17_SHTW9_01`     | TW    |    |   |          | WiZ Filament Bulb | EAN: 8718699786793 |
| `ESP56_SHTW3_01`     | TW    | 20 | 1 |          | WiZ G25 Filament bulb |
| `ESP56_SHTW3_01`     | TW    | 20 | 1 |          | WiZ filament G95 E27 |

## WiZ documentation

There is no real documentation of WiZ devices, but there is the [WiZ Pro API documentation](https://docs.pro.wizconnected.com) that shares a lot of information with the API that this library uses to directly communicate with the devices.
