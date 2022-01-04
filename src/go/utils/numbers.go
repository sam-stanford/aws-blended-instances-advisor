package utils

import "math"

// FloatsEqual checks if two floats are equal within a given threshold.
func FloatsEqual(a, b float64) bool {
	return math.Abs(a-b) < 0.001
}

// MinOfInts returns the smaller of two given ints.
func MinOfInts(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// MaxOfInts returns the larger of two given ints.
func MaxOfInts(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// MinOfFloats returns the smaller of two given floats.
func MinOfFloats(a, b float64) float64 {
	return math.Min(a, b)
}

// MaxOfFloats returns the larger of two given floats.
func MaxOfFloats(a, b float64) float64 {
	return math.Max(a, b)
}
