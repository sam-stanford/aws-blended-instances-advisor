package instances

import (
	"ec2-test/utils"
)

// Returns the index of the element in instances from startIndex (inclusive) to endIndex (exclusive) that has the smallest memory value greater than wantedMemory.
// Returns the index of the element with largest memory if no elements have a memory greater than wantedMemory.
// Returns an error if there is a problem with the given indexes.
//
// FindMinimumMemorySortedInstances should be called for improved performance on a sorted slice.
func FindMinimumMemory(instances []Instance, wantedMemory float64, startIndex int, endIndex int) (int, error) {
	err := utils.ValidateIndexes(len(instances), startIndex, endIndex)
	if err != nil {
		return -1, err
	}

	memoryValues := []float64{}
	for _, instance := range instances[startIndex:endIndex] {
		memoryValues = append(memoryValues, instance.MemoryGb)
	}

	index, err := utils.LinearSearchFloat(memoryValues, wantedMemory)
	if err != nil {
		return -1, utils.PrependToError(err, "error searching instances")
	}

	return startIndex + index, nil
}

// Returns the index of the element in sortedInstances from startIndex (inclusive) to endIndex (exclusive) that has the smallest memory value that is larger than wantedMemory.
// Undefined behaviour for a given slice of unsorted instances; however, an error will likely be returned.
//
// FindMinimumMemory should be used for unsorted slices.
func FindMinimumMemorySortedInstances(sortedInstances []Instance, wantedMemory float64, startIndex int, endIndex int) (int, error) {
	err := utils.ValidateIndexes(len(sortedInstances), startIndex, endIndex)
	if err != nil {
		return -1, err
	}

	memoryValues := []float64{}
	for _, instance := range sortedInstances[startIndex:endIndex] {
		memoryValues = append(memoryValues, instance.MemoryGb)
	}

	index, err := utils.BinarySearchFloat(memoryValues, wantedMemory)
	if err != nil {
		return -1, utils.PrependToError(err, "error searching instances")
	}

	return startIndex + index, nil
}

// TODO: Doc
func FindMinimumPrice(instances []Instance, wantedPrice float64, startIndex int, endIndex int) (int, error) {
	err := utils.ValidateIndexes(len(instances), startIndex, endIndex)
	if err != nil {
		return -1, err
	}

	prices := []float64{}
	for _, instance := range instances[startIndex:endIndex] {
		prices = append(prices, instance.PricePerHour)
	}

	index, err := utils.LinearSearchFloat(prices, wantedPrice)
	if err != nil {
		return -1, utils.PrependToError(err, "error searching instances")
	}

	return startIndex + index, nil
}

// TODO: Doc
func FindMinimumPriceSortedInstances(sortedInstances []Instance, wantedPrice float64, startIndex int, endIndex int) (int, error) {
	err := utils.ValidateIndexes(len(sortedInstances), startIndex, endIndex)
	if err != nil {
		return -1, err
	}

	prices := []float64{}
	for _, instance := range sortedInstances[startIndex:endIndex] {
		prices = append(prices, instance.PricePerHour)
	}

	index, err := utils.BinarySearchFloat(prices, wantedPrice)
	if err != nil {
		return -1, utils.PrependToError(err, "error searching instances")
	}

	return startIndex + index, nil
}

// TODO: Doc
func FindMinimumVcpu(instances []Instance, wantedVcpu int, startIndex int, endIndex int) (int, error) {
	err := utils.ValidateIndexes(len(instances), startIndex, endIndex)
	if err != nil {
		return -1, err
	}

	vcpus := []int{}
	for _, instance := range instances[startIndex:endIndex] {
		vcpus = append(vcpus, instance.Vcpus) // TODO: Rename Vcpus to vcpu cos it don't make sense
	}

	index, err := utils.LinearSearchInt(vcpus, wantedVcpu)
	if err != nil {
		return -1, utils.PrependToError(err, "error searching instances")
	}

	return startIndex + index, nil
}

// TODO: Doc
func FindMinimumVcpuSortedInstances(sortedInstances []Instance, wantedVcpu int, startIndex int, endIndex int) (int, error) {
	err := utils.ValidateIndexes(len(sortedInstances), startIndex, endIndex)
	if err != nil {
		return -1, err
	}

	vcpus := []int{}
	for _, instance := range sortedInstances[startIndex:endIndex] {
		vcpus = append(vcpus, instance.Vcpus)
	}

	index, err := utils.BinarySearchInt(vcpus, wantedVcpu)
	if err != nil {
		return -1, utils.PrependToError(err, "error searching instances")
	}

	return startIndex + index, nil
}

// TODO: Doc
func FindMinimumRevocationProbability(instances []Instance, wantedProbability float64, startIndex int, endIndex int) (int, error) {
	err := utils.ValidateIndexes(len(instances), startIndex, endIndex)
	if err != nil {
		return -1, err
	}

	probabilities := []float64{}
	for _, instance := range instances[startIndex:endIndex] {
		probabilities = append(probabilities, instance.RevocationProbability)
	}

	index, err := utils.LinearSearchFloat(probabilities, wantedProbability)
	if err != nil {
		return -1, utils.PrependToError(err, "error searching instances")
	}

	return startIndex + index, nil
}

// TODO: Doc
func FindMinimumRevocationProbabilitySortedInstances(sortedInstances []Instance, wantedProbability float64, startIndex int, endIndex int) (int, error) {
	err := utils.ValidateIndexes(len(sortedInstances), startIndex, endIndex)
	if err != nil {
		return -1, err
	}

	probabilities := []float64{}
	for _, instance := range sortedInstances[startIndex:endIndex] {
		probabilities = append(probabilities, instance.RevocationProbability)
	}

	index, err := utils.BinarySearchFloat(probabilities, wantedProbability)
	if err != nil {
		return -1, utils.PrependToError(err, "error searching instances")
	}

	return startIndex + index, nil
}
