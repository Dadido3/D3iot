# Emission

A library to handle and transform colors for light emitting things like LEDs.
It provides everything you need to convert between different color spaces.
This includes device color spaces with an arbitrary dimensionality.

The goal of this library is to implement the color space standards and their transformations as correctly as possible.
Though, there are some adjustments that make it more usable for light emitting devices, and less usable as a color management system for image data:

- There is no common media white point like in ICC profiles (Which use D50 for the profile connection space):

- There are no rendering intents, every transformation is handled as absolute colorimetric.
If you transform the CIE standard illuminant A into the device color space of a lamp, the perception of the produced light will resemble that standard illuminant as closely as possible.
No white point adjustment is done.

- Correct chromaticity is preferred over correct luminance.
When a color has a higher luminance than the device allows, the color is clipped in a way that doesn't change its chromaticity.

- It's not optimized for speed and doesn't use any hardware acceleration (GPU, ...).
This shouldn't matter when controlling lights, even if you update many thousands of lights hundreds of times per second.
But it's probably too slow to convert images pixel by pixel in a reasonable time.

- The absolute luminance is defined in lumen, instead of the usual cd/m².

Another goal was to leverage the go type system in a way that prevents users from transforming color spaces in wrong ways.

## Features

- Supports devices/modules with up to 6 emitters of different color (Up to 3 primaries like RGB, and up to 3 higher CRI "white" LEDs).
- Optimizes automatically for high CRI and high luminance output, if the device supports that (RGBW or RGBCW).
- Comes with support for the following color spaces:
  - CIE 1931 XYZ with absolute luminance
  - CIE 1931 XYZ with relative luminance
  - CIE 1931 xyY with absolute luminance
  - CIE 1931 xyY with relative luminance
  - CIE 1976 L\*a\*b\*
  - sRGB (IEC 61966-2-1:1999)
  - Linear and non-linear device color spaces with arbitrary dimensionality
- Transfer functions:
  - sRGB
  - Gamma curves
  - Custom transfer functions can implement `emission.TransferFunction`
- Light emitters:
  - CIE standard illuminants
  - Black body radiator

## Usage

### Color spaces

There are several color spaces that the library supports.
For a better overview we will group them into

- Human perception color spaces
- Device color spaces
- RGB color spaces

#### Human perception color spaces

These color spaces describe the human perception in some way.

``` go
xyzAbsolute := emission.CIE1931XYZAbs{X: 1500, Y: 1600, Z: 1400} // CIE 1931 XYZ color with an absolute luminance of 1600 lumen.
xyzRelative := emission.CIE1931XYZRel{X: 0.7, Y: 0.8, Z: 0.6}    // CIE 1931 XYZ color with a relative luminance of 0.8.
xyYAbsolute := emission.CIE1931xyYAbs{X: 0.33, Y: 0.33, LuminanceY: 1000} // CIE 1931 xyY color with an absolute luminance of 1000 lumen.
xyYRelative := emission.CIE1931xyYRel{X: 0.33, Y: 0.33, LuminanceY: 0.8}  // CIE 1931 xyY color with a relative luminance of 0.8.

// CIE 1976 L*a*b* color with a white point of D65.
cieLAB := emission.CIE1976LAB{
    L: 75,
    A: 29.1,
    B: 79.9,
    WhitePoint: emission.StandardIlluminantD65,
}
```

#### Device color spaces

A device color space (DCS) is just a vector.
The meaning of each component depends on the light device.

``` go
dcsVectorRGBCW := emission.DCSVector{0.5, 0.5, 0.5, 0.5, 0.5}
dcsVectorRGB := emission.DCSVector{0, 0, 0}
dcsVectorSingle := emission.DCSVector{0.2}
```

With the help of a transfer function you can linearize or de-linearize them.

``` go
linDCS := dcsVectorRGBCW.ClampedAndLinearized(emission.TransferFunctionStandardRGB)
dcsVector := linDCS.ClampedAndDeLinearized(emission.TransferFunctionStandardRGB)
```

Linear and non linear DCS vectors are of different type, so you can't just use the wrong one accidentally.

#### RGB color spaces

Any type that implements the `emission.RGB` represents a RGB color space.
For now the transformation is limited between RGB color spaces and `emission.CIE1931XYZRel`.

``` go
rgbColor := emission.StandardRGB{R: 0.5, G: 0.6, B: 0.7} // sRGB color.

xyzColor := rgbColor.CIE1931XYZRel().Absolute(300) // CIE 1931 XYZ color of that sRGB value with an absolute luminance of 300 lumen.

// Convert xyzColor into sRGB
var rgbColor2 emission.StandardRGB
xyzColor.TransformRGB(&rgbColor2) // Writes the color into rgbColor2.
```

### Color profiles

Setting up a color profile for a module is simple.
The library comes with a general implementation of its ColorProfile interface that supports most lamps:

``` go
colorProfile := &ColorProfileGeneral{
    WhitePointColor:  CIE1931XYZAbs{}.Sum(standardRGBRed, standardRGBGreen, standardRGBBlue),
    PrimaryColors:    TransformationLinDCSToXYZ{standardRGBRed, standardRGBGreen, standardRGBBlue},
    WhiteColors:      TransformationLinDCSToXYZ{CIE1931XYZAbs{30, 30, 30}},
    OutputLimiter:    OutputLimiterSum{3},
    TransferFunction: TransferFunctionStandardRGB,
}
colorProfile.MustInit() // Precalculate some internal values.
```

- `WhitePointColor` is the brightest color that module can output. It's just an absolute CIE 1931 XYZ color, with its luminance being in lumen.
- `PrimaryColors` contains a list of colors that span the color gamut of the module. Also in absolute luminance.
- `WhiteColors` contains all white or high CRI emitters of this module. Also in absolute luminance.
- `OutputLimiter` is a type that implements some method to limit the output channels in some way.
- `TransferFunction` is a type that implements the transformation between linear and non linear space.

### Value interface

Most objects in this library that represent some sort of emission implement the `emission.Value` interface.

### High CRI and high luminance optimization

By default, the library attempts to maximize the performance of high CRI white emitters, which should be desirable in most applications.
To disable this behavior, you can wrap any light object with `emission.NoWhiteOptimization`, which will prevent the use of any white emitters:

``` go
highCRILight := emission.StandardIlluminantA.Absolute(200)
lowCRILight := emission.NoWhiteOptimization{emission.StandardIlluminantA.Absolute(200)}
```

Using the `WiZ WI-FI BLE 100W A67 E27 922-65 RGB` as example, this will result in the following spectra:

![Alt text](images/WiZ%20WI-FI%20BLE%20100W%20A67%20E27%20922-65%20RGB%20CRI-comparison.png)

|                | Plot color | CRI (Ra) | Planckian temperature | WiZ pilot
| -------------- | ---------- | -------: | --------------------: | ---------
| `highCRILight` | black      |     96.3 |     2974K (DE2K -0.7) | `{State: true, Dimming: 100 %, R: 0, G: 0, B: 0, CW: 7, WW: 59}`
| `lowCRILight`  | red        |     16.6 |      3070K (DE2K 7.4) | `{State: true, Dimming: 100 %, R: 251, G: 151, B: 21, CW: 0, WW: 0}`

In direct comparison on a white surface, both light sources appear fairly identical.
However, due to "spiky" spectrum of `lowCRILight`, objects illuminated by it may show exaggerated or muted colors.
