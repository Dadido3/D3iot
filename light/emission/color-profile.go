// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package emission

// ColorProfile describes a set of light emitting things (a module) that together create a single color impression.
//
// This contains everything that is necessary to convert between the CIE 1931 XYZ color space and the device color space.
//
// Most devices contain just one module with some number of channels (E.g. RGBW light bulbs).
// But there can be multiple modules per device, e.g. multi headed lamps, addressable LED strips.
type ColorProfile interface {
	// Channels returns the dimensionality of the device color space.
	Channels() int

	// WhitePoint returns the white point as a CIE 1931 XYZ color.
	// This is also the brightest color a module can output.
	WhitePoint() CIE1931XYZAbs

	// ChannelPoints returns a list of channel colors.
	// Depending on the module type, this could be the colors for:
	//
	//	- Single white emitter.
	//	- RGB emitters.
	//	- RGB + white emitters.
	//	- RGB + cold white + warm white emitters.
	ChannelPoints() []CIE1931XYZAbs // TODO: Find better name, you don't get channel points for watching your light bulb glow

	// XYZToDCS takes a color and returns a vector in the device color space that reproduces the given color as closely as possible.
	//
	// Short: XYZ --> device color space.
	XYZToDCS(color CIE1931XYZAbs) DCSVector

	// DCSToXYZ takes a vector from the device color space and returns the color it represents.
	//
	// Short: Device color space --> XYZ.
	DCSToXYZ(v DCSVector) (CIE1931XYZAbs, error)

	TransferFunction() TransferFunction
}
