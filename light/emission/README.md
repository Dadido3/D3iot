# Emission

A library to handle and transform colors for light emitting things like LEDs.
This provides everything needed to transform from CIE 1931 XYZ or other color spaces into any device color space (DCS).
The transformation in the opposite direction is supported too.

This library works a bit differently than the usual color management pipelines.
There is no common media white point like in ICC profiles (Which use D50 for the profile connection space).

If you use the CIE 1931 XYZ color of the standard illuminant D65, the transformed DCS values will make a light bulb produce light that is perceived as that standard illuminant.

Another difference is how the luminance is defined.
The Y coordinate represents an absolute luminance with its unit being lumen, instead of the usual cd/m².

## Features

- XYZ to DCS transformation.
- DCS to XYZ transformation.
- Supports devices/modules with up to 6 emitters of different color (Up to 3 primaries like RGB, and up to 3 higher CRI "white" LEDs).
- Optimizes automatically for high CRI and high luminance output, if the device supports that (RGBW or RGBCW).
- Comes with support for the following color spaces:
  - CIE 1931 XYZ in absolute luminance
  - CIE 1931 XYZ in relative luminance
  - CIE 1931 xyY in absolute luminance
  - CIE 1931 xyY in relative luminance
  - CIE 1976 L\*a\*b\*
  - Linear and non-linear device colors spaces with arbitrary dimensionality
- Comes with the following transfer functions:
  - sRGB
  - Custom gamma curve
- Light emitters:
  - CIE standard illuminants
  - Black body radiator

## Usage

### Color spaces

There are two main color spaces that the lib uses.
The CIE 1931 XYZ color space, and the device color space (DCS).

``` go
xyzColor := CIE1931XYZAbs{1500, 1600, 1400}
dcsVector := DCSVector{0.5,0.5,0.5,0.5,0.5}
```

- `xyzColor` describes some color with a luminance of 1600 lumen.
- `dcsVector` describes a vector/color in the device color space.

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

### Transforming colors

`colorProfile` can then be used to transform between color spaces.

``` go
xyzColor, err := colorProfile.DCSToXYZ(dcsVector)
dcsVector := colorProfile.XYZToDCS(xyzColor)
```

### Value interface

Most objects in this library that represent some sort of emission implement the `emission.Value` interface.
