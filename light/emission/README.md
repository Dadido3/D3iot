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
