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
type TransformationLinDCSToXYZ []CIE1931XYZColor

// Multiplied returns the multiplication of t with a vector in the linear device color space.
// The result is a CIE 1931 XYZ color.
func (t TransformationLinDCSToXYZ) Multiplied(values []float64) (CIE1931XYZColor, error) {
	if len(t) != len(values) {
		return CIE1931XYZColor{}, fmt.Errorf("number of primaries %d doesn't match with the dimensionality %d of the DCS vector", len(t), len(values))
	}

	result := CIE1931XYZColor{}

	for i, primary := range t {
		result = result.Add(primary.Scale(values[i]))
	}

	return result, nil
}

// Inverted returns the inverted transformation matrix as list of column vectors.
//
// 1x3 and 2x3 matrices are handled in a special way.
func (t TransformationLinDCSToXYZ) Inverted() (TransformationXYZToLinDCS, error) {
	switch m, n := len(t), 3; m {
	case 1:
		// Add two new arbitrary vectors that are perpendicular to each other.
		tExt := TransformationLinDCSToXYZ{
			t[0],
			t[0].CrossProd(CIE1931XYZColor{1, 0, 0}),
			t[0].CrossProd(t[0].CrossProd(CIE1931XYZColor{1, 0, 0})),
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

// TransformationXYZToLinDCS represents a XYZ -> linear DCS transformation matrix by its row vectors.
//
//	t[0].X t[0].Y t[0].Z
//	t[1].X t[1].Y t[1].Z
//	t[2].X t[2].Y t[2].Z
//	...    ...    ...
type TransformationXYZToLinDCS []CIE1931XYZColor

// Multiplied returns the multiplication of t with a color in the XYZ color space.
// The result is a vector in the linear device color space.
func (t TransformationXYZToLinDCS) Multiplied(color CIE1931XYZColor) []float64 {
	result := make([]float64, len(t))

	for i, primary := range t {
		result[i] += primary.X*color.X + primary.Y*color.Y + primary.Z*color.Z
	}

	return result
}
