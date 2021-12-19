// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package emission

import "fmt"

// TransformationLinDCSToXYZ represents a linear DCS -> XYZ transformation matrix by its column vectors (color primaries).
//
//	t[0].X t[1].X t[2].X ...
//	t[0].Y t[1].Y t[2].Y ...
//	t[0].Z t[1].Z t[2].Z ...
type TransformationLinDCSToXYZ []CIE1931XYZAbs

// DCSChannels returns the dimensionality of the device color space.
// This is equal to m in the m x n transformation matrix.
func (t TransformationLinDCSToXYZ) DCSChannels() int {
	return len(t)
}

// Scaled returns t scaled by the scalar s.
func (t TransformationLinDCSToXYZ) Scaled(s float64) TransformationLinDCSToXYZ {
	result := make(TransformationLinDCSToXYZ, 0, t.DCSChannels())
	for _, color := range t {
		result = append(result, color.Scaled(s))
	}

	return result
}

// Multiplied returns the multiplication of t with an in the linear device color space.
// The result is a CIE 1931 XYZ color.
func (t TransformationLinDCSToXYZ) Multiplied(v LinDCSVector) (CIE1931XYZAbs, error) {
	if t.DCSChannels() != v.Channels() {
		return CIE1931XYZAbs{}, fmt.Errorf("number of primaries %d doesn't match with the dimensionality %d of the DCS vector", t.DCSChannels(), v.Channels())
	}

	result := CIE1931XYZAbs{}

	for i, primary := range t {
		result = result.Sum(primary.Scaled(v[i]))
	}

	return result, nil
}

// Inverted returns the inverted transformation matrix as list of column vectors.
//
// 1x3 and 2x3 matrices are handled in a special way.
// 0x3 matrices will just return an empty inverse transformation.
func (t TransformationLinDCSToXYZ) Inverted() (TransformationXYZToLinDCS, error) {
	switch m, n := t.DCSChannels(), 3; m {
	case 0:
		return nil, nil
	case 1:
		// Add two new arbitrary vectors that each are perpendicular to the others.
		tExt := TransformationLinDCSToXYZ{
			t[0],
			t[0].CrossProd(CIE1931XYZAbs{1, 0, 0}),
			t[0].CrossProd(t[0].CrossProd(CIE1931XYZAbs{1, 0, 0})),
		}
		inv, err := tExt.Inverted()
		if err != nil {
			return nil, err
		}
		return inv[:1], nil

	case 2:
		// Add one new arbitrary vector that is perpendicular the other two.
		tExt := TransformationLinDCSToXYZ{
			t[0],
			t[1],
			t[0].CrossProd(t[1]),
		}
		inv, err := tExt.Inverted()
		if err != nil {
			return nil, err
		}
		return inv[:2], nil

	case 3:
		// Calculate inverse of 3x3 matrix.
		det := t[0].X*(t[1].Y*t[2].Z-t[2].Y*t[1].Z) -
			t[0].Y*(t[1].X*t[2].Z-t[1].Z*t[2].X) +
			t[0].Z*(t[1].X*t[2].Y-t[1].Y*t[2].X)

		if det == 0 {
			return nil, fmt.Errorf("determinant is zero")
		}

		invDet := 1 / det

		return TransformationXYZToLinDCS{
			{
				(t[1].Y*t[2].Z - t[2].Y*t[1].Z) * invDet,
				(t[1].Z*t[2].X - t[1].X*t[2].Z) * invDet,
				(t[1].X*t[2].Y - t[2].X*t[1].Y) * invDet,
			},
			{
				(t[0].Z*t[2].Y - t[0].Y*t[2].Z) * invDet,
				(t[0].X*t[2].Z - t[0].Z*t[2].X) * invDet,
				(t[2].X*t[0].Y - t[0].X*t[2].Y) * invDet,
			},
			{
				(t[0].Y*t[1].Z - t[0].Z*t[1].Y) * invDet,
				(t[1].X*t[0].Z - t[0].X*t[1].Z) * invDet,
				(t[0].X*t[1].Y - t[1].X*t[0].Y) * invDet,
			},
		}, nil

	default:
		return nil, fmt.Errorf("unsupported transformation matrix dimension %d x %d", m, n)
	}

}

// MustInverted returns the inverted transformation matrix as list of column vectors.
//
// This is an alias of t.Inverted() that does panic on any error.
func (t TransformationLinDCSToXYZ) MustInverted() TransformationXYZToLinDCS {
	inv, err := t.Inverted()
	if err != nil {
		panic(err)
	}
	return inv
}

// TransformationXYZToLinDCS represents a XYZ -> linear DCS transformation matrix by its row vectors.
//
//	t[0].X t[0].Y t[0].Z
//	t[1].X t[1].Y t[1].Z
//	t[2].X t[2].Y t[2].Z
//	...    ...    ...
type TransformationXYZToLinDCS []CIE1931XYZAbs

// DCSChannels returns the dimensionality of the device color space.
// This is equal to n in a m x n matrix.
func (t TransformationXYZToLinDCS) DCSChannels() int {
	return len(t)
}

// Multiplied returns the multiplication of t with a color in the XYZ color space.
// The result is an unclamped vector in the linear device color space.
func (t TransformationXYZToLinDCS) Multiplied(color CIE1931XYZAbs) LinDCSVector {
	result := make(LinDCSVector, t.DCSChannels())

	for i, primary := range t {
		result[i] += primary.X*color.X + primary.Y*color.Y + primary.Z*color.Z
	}

	return result
}
