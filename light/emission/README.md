# Emission

A library to handle and transform colors for light emitting things like LEDs.
This provides everything needed to transform from CIE 1931 XYZ or other color spaces into any device color space (DCS).
The transformation in the opposite direction is supported too.

This library works a bit differently than the usual color management pipelines.
There is no common media white point like in ICC profiles (Which uses D50 for the profile connection space).
If you use the CIE 1931 XYZ color of the standard illuminant D65, the transformed DCS values will make a light bulb produce light that is perceived as that standard illuminant.

Another difference is how the luminance is defined.
The Y coordinate represents an absolute luminance with its unit being lumen, instead of the usual cd/m².

## Features

- XYZ to DCS transformation.
- DCS to XYZ transformation.
- Supports lights/modules with up to 6 differently emitters of different color (Up to 3 primaries like RGB, and up to 3 higher CRI "white" LEDs).
- Optimizes automatically for high CRI and high luminance output, if the device supports that (RGBW or RGBCW).

## Usage

### Color spaces

There are two main color spaces that the lib uses.
The CIE 1931 XYZ color space, and the device color space (DCS).

``` go
xyzColor := CIE1931XYZColor{1500, 1600, 1400}
dcsColor := DCSColor{0.5,0.5,0.5,0.5,0.5}
```

- `xyzColor` describes some color with a luminance of 1600 lumen.
- `dcsColor` describes a color in the device color space.

### Module profiles

Setting up a profile for a module that can be used to transform color spaces is simple.
A module is just a thing that can output color by using a set of light emitting things of different color.

``` go
moduleProfile := &ModuleProfileGeneral{
    WhitePointColor:  CIE1931XYZColor{}.Sum(standardRGBRed, standardRGBGreen, standardRGBBlue),
    PrimaryColors:    TransformationLinDCSToXYZ{standardRGBRed, standardRGBGreen, standardRGBBlue},
    WhiteColors:      TransformationLinDCSToXYZ{CIE1931XYZColor{30, 30, 30}},
    OutputLimiter:    OutputLimiterSum{3},
    TransferFunction: TransferFunctionStandardRGB,
}
moduleProfile.MustInit() // Precalculate some internal values.
```

- `WhitePointColor` is the brightest color that module can output. It's just a CIE 1931 XYZ color.
- `PrimaryColors` contains a list of colors that span the color gamut of the module.
- `WhiteColors` contains all white or high CRI emitters of this module.
- `OutputLimiter` is a type that implements some method to limit the output channels in some way.
- `TransferFunction` is a type that implements the transformation between linear and non linear space.

### Transforming colors

`moduleProfile` can then be used to transform between color spaces.

``` go
xyzColor, err := moduleProfile.DCSToXYZ(dcsColor)
dcsColor := moduleProfile.XYZToDCS(xyzColor)
```

### Value interface

Most objects in this library that represent some sort of emission implement the `emission.Value` interface.
