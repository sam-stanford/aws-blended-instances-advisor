package utils

import (
	"errors"
	"fmt"
)

// LinearSearchInt searches for a given int value in a slice using a linear search.
func LinearSearchInt(slice []int, value int) (int, error) {
	if len(slice) == 0 {
		return -1, errors.New("cannot search empty slice")
	}

	answer := 0

	for thisIndex, thisValue := range slice {
		if thisValue == value {
			return thisIndex, nil
		}
		if (slice[answer] < value && thisValue >= slice[answer]) ||
			(thisValue > value && thisValue < slice[answer]) {
			answer = thisIndex
		}
	}

	return answer, nil
}

// LinearSearchFloat searches for a given float value in a slice using a linear search.
func LinearSearchFloat(slice []float64, value float64) (int, error) {
	if len(slice) == 0 {
		return -1, errors.New("cannot search empty slice")
	}

	answer := 0

	for thisIndex, thisValue := range slice {
		if thisValue == value {
			return thisIndex, nil
		}
		if (slice[answer] < value && thisValue >= slice[answer]) ||
			(thisValue > value && thisValue < slice[answer]) {
			answer = thisIndex
		}
	}

	return answer, nil
}

// LinearSearchString searches for a given string value in a slice using a linear search.
func LinearSearchString(slice []string, value string) (int, error) {
	if len(slice) == 0 {
		return -1, errors.New("cannot search empty slice")
	}

	answer := 0

	for thisIndex, thisValue := range slice {
		if thisValue == value {
			return thisIndex, nil
		}
		if (slice[answer] < value && thisValue >= slice[answer]) ||
			(thisValue > value && thisValue < slice[answer]) {
			answer = thisIndex
		}
	}

	return answer, nil
}

// BinarySearchInt searches for a given int value in a slice using binary search.
//
// The behvious of BinarySearchInt is undefined if the provided slice is not sorted.
func BinarySearchInt(sortedSlice []int, value int) (int, error) {
	if len(sortedSlice) == 0 {
		return -1, errors.New("cannot search empty slice")
	}

	left, right := 0, len(sortedSlice)-1
	for left <= right {
		if left == right {
			return left, nil
		}

		leftVal, rightVal := sortedSlice[left], sortedSlice[right]

		mid := left + (right-left)/2
		midVal := sortedSlice[mid]

		if mid == left {
			if value > leftVal {
				return right, nil
			}
			return left, nil
		}

		if leftVal > rightVal {
			return -1, fmt.Errorf(
				"provided slice is not sorted. Index %d (val=%d) greater than %d (val=%d)",
				left, leftVal, right, rightVal,
			)
		}
		if midVal > rightVal {
			return -1, fmt.Errorf(
				"provided slice is not sorted. Index %d (val=%d) greater than %d (val=%d)",
				mid, midVal, right, rightVal,
			)
		}

		if midVal == value {
			// Find first element == value
			for midVal == value {
				if mid == 0 {
					return 0, nil
				}
				mid -= 1
				midVal = sortedSlice[mid]
			}
			return mid + 1, nil
		}

		if value < midVal {
			right = mid
		} else {
			left = mid
		}
	}

	return -1, errors.New("failed to find wanted index")
}

// BinarySearchFloat searches for a given float value in a slice using binary search.
//
// The behvious of BinarySearchFloat is undefined if the provided slice is not sorted.
func BinarySearchFloat(sortedSlice []float64, value float64) (int, error) {
	if len(sortedSlice) == 0 {
		return -1, errors.New("cannot search empty slice")
	}

	left, right := 0, len(sortedSlice)-1
	for left <= right {
		if left == right {
			return left, nil
		}

		leftVal, rightVal := sortedSlice[left], sortedSlice[right]

		mid := left + (right-left)/2
		midVal := sortedSlice[mid]

		if mid == left {
			if value > leftVal {
				return right, nil
			}
			return left, nil
		}

		if leftVal > rightVal {
			return -1, fmt.Errorf(
				"provided slice is not sorted. Index %d (val=%f) greater than %d (val=%f)",
				left, leftVal, right, rightVal,
			)
		}
		if midVal > rightVal {
			return -1, fmt.Errorf(
				"provided slice is not sorted. Index %d (val=%f) greater than %d (val=%f)",
				mid, midVal, right, rightVal,
			)
		}

		if midVal == value {
			// Find first element == value
			for midVal == value {
				if mid == 0 {
					return 0, nil
				}
				mid -= 1
				midVal = sortedSlice[mid]
			}
			return mid + 1, nil
		}

		if value < midVal {
			right = mid
		} else {
			left = mid
		}
	}

	return -1, errors.New("failed to find wanted index")
}

// BinarySearchString searches for a given string value in a slice using binary search.
//
// The behvious of BinarySearchString is undefined if the provided slice is not sorted.
func BinarySearchString(sortedSlice []string, value string) (int, error) {
	if len(sortedSlice) == 0 {
		return -1, errors.New("cannot search empty slice")
	}

	left, right := 0, len(sortedSlice)-1
	for left <= right {
		if left == right {
			return left, nil
		}

		leftVal, rightVal := sortedSlice[left], sortedSlice[right]

		mid := left + (right-left)/2
		midVal := sortedSlice[mid]

		if mid == left {
			if value > leftVal {
				return right, nil
			}
			return left, nil
		}

		if leftVal > rightVal {
			return -1, fmt.Errorf(
				"provided slice is not sorted. Index %d (val=%s) greater than %d (val=%s)",
				left, leftVal, right, rightVal,
			)
		}
		if midVal > rightVal {
			return -1, fmt.Errorf(
				"provided slice is not sorted. Index %d (val=%s) greater than %d (val=%s)",
				mid, midVal, right, rightVal,
			)
		}

		if midVal == value {
			// Find first element == value
			for midVal == value {
				if mid == 0 {
					return 0, nil
				}
				mid -= 1
				midVal = sortedSlice[mid]
			}
			return mid + 1, nil
		}

		if value < midVal {
			right = mid
		} else {
			left = mid
		}
	}

	return -1, errors.New("failed to find wanted index")
}
