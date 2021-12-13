// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package emission

// ModuleDescriptor describes a group of LEDs that together generate a single color impression.
//
// This contains everything that is necessary to convert from the CIE 1931 XYZ color space into the device color space, and vice versa.
//
// Most devices contain just one module with some number of channels (E.g. RGBW light bulbs).
// But there can be multiple modules per device, e.g. multi headed lamps, addressable LED strips.
type ModuleDescriptor interface {
	// Channels returns the dimensionality of the device color space.
	Channels() int

	// WhitePoint returns the white point as a CIE 1931 XYZ color.
	// This is also the brightest color a module can output.
	WhitePoint() CIE1931XYZColor

	// ChannelPoints returns a list of channel colors.
	// Depending on the module type, this could be the colors for:
	//
	//	- Single white emitter.
	//	- RGB emitters.
	//	- RGB + white emitters.
	//	- RGB + cold white + warm white emitters.
	ChannelPoints() []CIE1931XYZColor

	// XYZToDCS takes a color and returns a vector in the device color space that reproduces the given color as closely as possible.
	//
	// Short: XYZ --> device color space.
	XYZToDCS(color CIE1931XYZColor) (DCSColor, error)

	// DCSToXYZ takes a vector from the device color space and returns the color it represents.
	//
	// Short: Device color space --> XYZ.
	DCSToXYZ(v DCSColor) (CIE1931XYZColor, error)
}
