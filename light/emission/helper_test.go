package emission

import "math"

var float64Precision = 1e-9

func compareFloat64(a, b float64) bool {
	return math.Abs(a-b) > float64Precision
}
