// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package wiz

// products contains known WiZ light products.
// This can't be a map, as otherwise matching would get nondeterministic.
var products = []Product{
	{
		moduleName:  "ESP03_SHRGB1W_01",
		deviceClass: deviceClassRGBTW,
		maxRGBWSum:  newPtrUInt(512),
		minTemp:     newPtrUInt(2200), maxTemp: newPtrUInt(6500),
	},
}
