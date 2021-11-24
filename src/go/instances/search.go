package instances

import (
	"ec2-test/utils"
	"errors"
)

// Returns the index of the element in instances from startIndex (inclusive) to endIndex (exclusive) that has the smallest memory value greater than wantedMemory.
// Returns the index of the element with largest memory if no elements have a memory greater than wantedMemory.
// Returns an error if there is a problem with the given indexes.
//
// GetIndexOfMinimumMemoryFromSortedInstances should be called for improved performance on a sorted slice.
func GetIndexOfMinimumMemoryFromInstances(instances []Instance, wantedMemory float64, startIndex int, endIndex int) (int, error) {
	err := utils.ValidateIndexes(len(instances), startIndex, endIndex)
	if err != nil {
		return -1, err
	}

	closestIndex := startIndex
	closestMemory := instances[startIndex].MemoryGb

	for index, instance := range instances[startIndex:endIndex] {
		if (closestMemory < wantedMemory && instance.MemoryGb > closestMemory) ||
			(closestMemory >= wantedMemory && instance.MemoryGb < closestMemory) {
			closestIndex = index
			closestMemory = instance.MemoryGb
		}
	}

	return closestIndex, nil
}

// Returns the index of the element in sortedInstances from startIndex (inclusive) to endIndex (exclusive) that has the smallest memory value that is larger than wantedMemory.
// Undefined behaviour for a given slice of unsorted instances; however, an error will likely be returned.
//
// GetIndexOfMinimumMemoryFromInstances should be used for unsorted slices.
func GetIndexOfMinimumMemoryFromSortedInstances(sortedInstances []Instance, wantedMemory float64, startIndex int, endIndex int) (int, error) {
	// TODO: Validate indexes
	left, right := startIndex, endIndex-1
	for left <= right {
		if left == right {
			return left, nil
		}

		leftVal, rightVal := sortedInstances[left].MemoryGb, sortedInstances[right].MemoryGb

		mid := left + (right-left)/2
		midVal := sortedInstances[mid].MemoryGb

		if mid == left {
			if wantedMemory > leftVal {
				return right, nil
			}
			return left, nil
		}

		if leftVal > rightVal || midVal > rightVal {
			return -1, errors.New("provided instances are not sorted with respect to memory")
		}

		if midVal == wantedMemory {
			return mid, nil
		}

		if wantedMemory < midVal {
			right = mid
		} else {
			left = mid
		}
	}

	return -1, errors.New("failed to find wanted index")
}
