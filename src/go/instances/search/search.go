package search

import (
	. "aws-blended-instances-advisor/instances"
	"aws-blended-instances-advisor/instances/sort"
	"aws-blended-instances-advisor/utils"
)

// Returns the index of the element in instances from startIndex (inclusive) to endIndex (exclusive) that has the smallest memory value greater than wantedMemory.
// Returns the index of the element with largest memory if no elements have a memory greater than wantedMemory.
// Returns an error if there is a problem with the given indexes.
//
// FindMemorySorted should be called for improved performance on a sorted slice.
func FindMemory(instances []*Instance, wantedMemory float64, startIndex int, endIndex int) (int, error) {
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
// FindMemory should be used for unsorted slices.
func FindMemorySorted(sortedInstances []*Instance, wantedMemory float64, startIndex int, endIndex int) (int, error) {
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

func SortAndFindMemory(instances []*Instance, wantedMemory float64, startIndex, endIndex int) (int, error) {
	sort.SortInstancesByMemory(instances, startIndex, endIndex)
	return FindMemorySorted(instances, wantedMemory, startIndex, endIndex)
}

func FindPrice(instances []*Instance, wantedPrice float64, startIndex int, endIndex int) (int, error) {
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

func FindPriceSorted(sortedInstances []*Instance, wantedPrice float64, startIndex int, endIndex int) (int, error) {
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

func SortAndFindPrice(instances []*Instance, wantedPrice float64, startIndex, endIndex int) (int, error) {
	sort.SortInstancesByPrice(instances, startIndex, endIndex)
	return FindPriceSorted(instances, wantedPrice, startIndex, endIndex)
}

func FindVcpu(instances []*Instance, wantedVcpu int, startIndex int, endIndex int) (int, error) {
	err := utils.ValidateIndexes(len(instances), startIndex, endIndex)
	if err != nil {
		return -1, err
	}

	vcpus := []int{}
	for _, instance := range instances[startIndex:endIndex] {
		vcpus = append(vcpus, instance.Vcpu)
	}

	index, err := utils.LinearSearchInt(vcpus, wantedVcpu)
	if err != nil {
		return -1, utils.PrependToError(err, "error searching instances")
	}

	return startIndex + index, nil
}

func FindVcpuSorted(sortedInstances []*Instance, wantedVcpu int, startIndex int, endIndex int) (int, error) {
	err := utils.ValidateIndexes(len(sortedInstances), startIndex, endIndex)
	if err != nil {
		return -1, err
	}

	vcpus := []int{}
	for _, instance := range sortedInstances[startIndex:endIndex] {
		vcpus = append(vcpus, instance.Vcpu)
	}

	index, err := utils.BinarySearchInt(vcpus, wantedVcpu)
	if err != nil {
		return -1, utils.PrependToError(err, "error searching instances")
	}

	return startIndex + index, nil
}

func SortAndFindVcpu(instances []*Instance, wantedVcpu int, startIndex, endIndex int) (int, error) {
	sort.SortInstancesByVcpu(instances, startIndex, endIndex)
	return FindVcpuSorted(instances, wantedVcpu, startIndex, endIndex)
}

func FindRevocationProbability(instances []*Instance, wantedProbability float64, startIndex int, endIndex int) (int, error) {
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

func FindRevocationProbabilitySorted(sortedInstances []*Instance, wantedProbability float64, startIndex int, endIndex int) (int, error) {
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

func SortAndFindRevocationProbability(instances []*Instance, wantedProbability float64, startIndex, endIndex int) (int, error) {
	sort.SortInstancesByRevocationProbability(instances, startIndex, endIndex)
	return FindRevocationProbabilitySorted(instances, wantedProbability, startIndex, endIndex)
}
