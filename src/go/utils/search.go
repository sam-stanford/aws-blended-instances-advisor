package utils

import (
	"errors"
	"fmt"
)

// TODO: Docs & testing

func LinearSearchInt(slice []int, value int) (int, error) {
	if len(slice) == 0 {
		return -1, errors.New("cannot search empty slice")
	}

	answer := 0

	for thisIndex, thisValue := range slice {
		if (slice[answer] < value && thisValue > slice[answer]) ||
			(thisValue >= value && thisValue < slice[answer]) {
			answer = thisIndex
		}
	}

	return answer, nil
}

func LinearSearchFloat(slice []float64, value float64) (int, error) {
	if len(slice) == 0 {
		return -1, errors.New("cannot search empty slice")
	}

	answer := 0

	for thisIndex, thisValue := range slice {
		if (slice[answer] < value && thisValue > slice[answer]) ||
			(thisValue >= value && thisValue < slice[answer]) {
			answer = thisIndex
		}
	}

	return answer, nil
}

func LinearSearchString(slice []string, value string) (int, error) {
	if len(slice) == 0 {
		return -1, errors.New("cannot search empty slice")
	}

	answer := 0

	for thisIndex, thisValue := range slice {
		if (slice[answer] < value && thisValue > slice[answer]) ||
			(thisValue >= value && thisValue < slice[answer]) {
			answer = thisIndex
		}
	}

	return answer, nil
}

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
			for mid >= 0 && midVal == value {
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
			for mid >= 0 && midVal == value {
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
			for mid >= 0 && midVal == value {
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
