# Light package

This contains the light interface that supported light devices implement.

## Usage

For now there is no way to discover devices, the only way to connect to devices is by adding them manually.

``` go
light, err := wiz.NewLight("192.168.1.123:38899")
```

Check the packages for the devices you want to connect with for more details:

- [WiZ](drivers/wiz/)

### Modules

A light device contains modules.
These are single units that contain one or multiple different types of light emitting things (LEDs, ...).
Each module is responsible for a single color impression, and each module can be controlled independently.

Most devices just contain one module (E.g. RGB light bulbs, ...).
You can use `light.Modules()` to get the number of modules.

To get information about a module's color space and abilities, use `light.ColorProfiles()`.

``` go
// Get the color profiles of all modules.
// There is always at least one module!
profile := light.ColorProfiles()[0]

whitePoint := profile.WhitePointColor()
halfWhitePoint := whitePoint.Scaled(0.5)

primaryColors := profile.ChannelPoints()
```

- `whitePoint` is the brightest color the module can output.
- `halfWhitePoint` has the same color, but half the luminance.
- `primaryColors` contains a list of the colors that span color space of the device.

### Set colors

First you need to define colors.
The following shows some example light emitters.

``` go
// XYZ values of the standard illuminant D65 with an absolute luminance of 500 lumen.
xyzColor := emission.CIE1931XYZRel{X: 0.95047, Y: 1, Z: 1.08883}.Absolute(500) 

// xyY values of the standard illuminant D65 with an absolute luminance of 500 lumen.
xyYColor := emission.CIE1931xyYAbs{X: 0.31271, Y: 0.32902, LuminanceY: 500}

// Value in linear device color space. The resulting color depends on the device.
// The example is for modules with 1 channels.
singleChannel := emission.LinDCSVector{0.5}

// Value in device color space. The resulting color depends on the device.
// The example is for modules with 3 channels (Like RGB).
dcsRGBColor := emission.DCSVector{0.1, 0.2, 1.0}

// Black body radiator with 5000 K color temperature and a luminance of 500 lumen.
blackbodyColor := emission.BlackBodyFixed{Temperature: 5000, Luminance: 500}

// CIE standard illuminant A with an absolute luminance of 400 lumen.
incandescent := emission.StandardIlluminantA.Absolute(400)
```

These colors can then be passed to any light by using the `SetColors()` method.

``` go
err := light.SetColor(xyYColor)
```

If you have a device with multiple modules, like addressable RGB strips, you can pass multiple colors to it.
Even if they are in different color spaces.

``` go
err := light.SetColor(xyYColor, XYZColor, dcsRGBColor, blackbodyColor)
err := light.SetColor(colors...)
```

### Get colors

Querying a light device for its current set of colors is simple too.

``` go
var xyYColor emission.CIE1931xyYAbs
err := light.GetColor(&xyYColor)
```

The result will be written into `xyYColor`.
You can write into any type that implements the `emission.Value` interface.
If you write into `emission.DCSVector`, you will get a vector in the device color space.
That is the raw RGBW values or whatever defines the color space of that device.
