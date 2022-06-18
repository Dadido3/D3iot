// Copyright (c) 2021-2022 David Vogel
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
	ChannelPoints() []CIE1931XYZAbs

	// XYZToDCS takes a color and returns a vector in the device color space that reproduces the given color as closely as possible.
	//
	// Short: XYZ --> device color space.
	XYZToDCS(color CIE1931XYZAbs) DCSVector

	// DCSToXYZ takes a vector from the device color space and returns the color it represents.
	//
	// Short: Device color space --> XYZ.
	DCSToXYZ(v DCSVector) (CIE1931XYZAbs, error)

	TransferFunction() TransferFunction

	// NoWhiteOptimizationColorProfile returns a copy of this color profile with all high CRI white emitters disabled.
	//
	// The output of such profile can not optimize for high CRI and high luminance.
	// This will cause that the emitted color is only constructed by the primary colors (e.g. red, green, blue).
	// This will not change the emitted color, but the maximum brightness may be reduced and the CRI is not as good as it could be.
	// One use-case could be to eliminate any timing discrepancy between high CRI whites and primary color emitters, as these two classes of emitters may be filtered (by a low-pass) in a different way.
	//
	// Some color profiles do not support this and will just return themselves.
	NoWhiteOptimizationColorProfile() ColorProfile
}
