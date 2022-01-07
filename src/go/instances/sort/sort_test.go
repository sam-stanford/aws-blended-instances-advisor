package sort

import (
	awsTypes "aws-blended-instances-advisor/aws/types"
	. "aws-blended-instances-advisor/instances"
	"reflect"
	"testing"
)

type instanceSortTest struct {
	instances  []*Instance
	start, end int
	expected   []*Instance
}

func TestSortInstancesByPrice(t *testing.T) {
	i0 := &Instance{Name: "0", PricePerHour: 0.00001}
	i1 := &Instance{Name: "1", PricePerHour: 0.00002}
	i2 := &Instance{Name: "2", PricePerHour: 0.00003}
	i3 := &Instance{Name: "3", PricePerHour: 0.00010}

	sortedSlice := []*Instance{i0, i1, i2, i3}

	tests := []instanceSortTest{
		{instances: []*Instance{i0, i1, i2, i3}, start: 0, end: 4, expected: sortedSlice},
		{instances: []*Instance{i1, i2, i3, i0}, start: 0, end: 4, expected: sortedSlice},
		{instances: []*Instance{i2, i3, i0, i1}, start: 0, end: 4, expected: sortedSlice},
		{instances: []*Instance{i3, i2, i1, i0}, start: 0, end: 2, expected: []*Instance{i2, i3, i1, i0}},
		{instances: []*Instance{i3, i2, i1, i0}, start: 2, end: 4, expected: []*Instance{i3, i2, i0, i1}},
		{instances: []*Instance{i1, i1, i0, i0}, start: 0, end: 4, expected: []*Instance{i0, i0, i1, i1}},
	}

	for index, test := range tests {
		SortInstancesByPrice(test.instances, test.start, test.end)
		if !reflect.DeepEqual(test.instances, test.expected) {
			t.Fatalf(
				"Instances are not sorted correctly for test %d. Wanted: %+v, got: %+v",
				index,
				test.instances,
				test.expected,
			)
		}
	}
}

func TestSortInstancesByMemory(t *testing.T) {
	i0 := &Instance{Name: "0", MemoryGb: 4}
	i1 := &Instance{Name: "1", MemoryGb: 8}
	i2 := &Instance{Name: "2", MemoryGb: 16}
	i3 := &Instance{Name: "3", MemoryGb: 62.4}

	sortedSlice := []*Instance{i0, i1, i2, i3}

	tests := []instanceSortTest{
		{instances: []*Instance{i0, i1, i2, i3}, start: 0, end: 4, expected: sortedSlice},
		{instances: []*Instance{i1, i2, i3, i0}, start: 0, end: 4, expected: sortedSlice},
		{instances: []*Instance{i2, i3, i0, i1}, start: 0, end: 4, expected: sortedSlice},
		{instances: []*Instance{i3, i2, i1, i0}, start: 0, end: 2, expected: []*Instance{i2, i3, i1, i0}},
		{instances: []*Instance{i3, i2, i1, i0}, start: 2, end: 4, expected: []*Instance{i3, i2, i0, i1}},
		{instances: []*Instance{i1, i1, i0, i0}, start: 0, end: 4, expected: []*Instance{i0, i0, i1, i1}},
	}

	for index, test := range tests {
		SortInstancesByMemory(test.instances, test.start, test.end)
		if !reflect.DeepEqual(test.instances, test.expected) {
			t.Fatalf(
				"Instances are not sorted correctly for test %d. Wanted: %+v, got: %+v",
				index,
				test.instances,
				test.expected,
			)
		}
	}
}

func TestSortInstancesByVcpu(t *testing.T) {
	i0 := &Instance{Name: "0", Vcpu: 1}
	i1 := &Instance{Name: "1", Vcpu: 4}
	i2 := &Instance{Name: "2", Vcpu: 16}
	i3 := &Instance{Name: "3", Vcpu: 128}

	sortedSlice := []*Instance{i0, i1, i2, i3}

	tests := []instanceSortTest{
		{instances: []*Instance{i0, i1, i2, i3}, start: 0, end: 4, expected: sortedSlice},
		{instances: []*Instance{i1, i2, i3, i0}, start: 0, end: 4, expected: sortedSlice},
		{instances: []*Instance{i2, i3, i0, i1}, start: 0, end: 4, expected: sortedSlice},
		{instances: []*Instance{i3, i2, i1, i0}, start: 0, end: 2, expected: []*Instance{i2, i3, i1, i0}},
		{instances: []*Instance{i3, i2, i1, i0}, start: 2, end: 4, expected: []*Instance{i3, i2, i0, i1}},
		{instances: []*Instance{i1, i1, i0, i0}, start: 0, end: 4, expected: []*Instance{i0, i0, i1, i1}},
	}

	for index, test := range tests {
		SortInstancesByVcpu(test.instances, test.start, test.end)
		if !reflect.DeepEqual(test.instances, test.expected) {
			t.Fatalf(
				"Instances are not sorted correctly for test %d. Wanted: %+v, got: %+v",
				index,
				test.instances,
				test.expected,
			)
		}
	}
}

func TestSortInstancesByRevocationProbability(t *testing.T) {
	i0 := &Instance{Name: "0", RevocationProbability: 0}
	i1 := &Instance{Name: "1", RevocationProbability: 0.05}
	i2 := &Instance{Name: "2", RevocationProbability: 0.10}
	i3 := &Instance{Name: "3", RevocationProbability: 0.215}

	sortedSlice := []*Instance{i0, i1, i2, i3}

	tests := []instanceSortTest{
		{instances: []*Instance{i0, i1, i2, i3}, start: 0, end: 4, expected: sortedSlice},
		{instances: []*Instance{i1, i2, i3, i0}, start: 0, end: 4, expected: sortedSlice},
		{instances: []*Instance{i2, i3, i0, i1}, start: 0, end: 4, expected: sortedSlice},
		{instances: []*Instance{i3, i2, i1, i0}, start: 0, end: 2, expected: []*Instance{i2, i3, i1, i0}},
		{instances: []*Instance{i3, i2, i1, i0}, start: 2, end: 4, expected: []*Instance{i3, i2, i0, i1}},
		{instances: []*Instance{i1, i1, i0, i0}, start: 0, end: 4, expected: []*Instance{i0, i0, i1, i1}},
	}

	for index, test := range tests {
		SortInstancesByRevocationProbability(test.instances, test.start, test.end)
		if !reflect.DeepEqual(test.instances, test.expected) {
			t.Fatalf(
				"Instances are not sorted correctly for test %d. Wanted: %+v, got: %+v",
				index,
				test.instances,
				test.expected,
			)
		}
	}
}

func TestSortInstancesByRegion(t *testing.T) {
	i0 := &Instance{Name: "0", Region: awsTypes.ApNorthEast1}
	i1 := &Instance{Name: "1", Region: awsTypes.EuNorth1}
	i2 := &Instance{Name: "2", Region: awsTypes.UsEast1}
	i3 := &Instance{Name: "3", Region: awsTypes.UsWest2}

	sortedSlice := []*Instance{i0, i1, i2, i3}

	tests := []instanceSortTest{
		{instances: []*Instance{i0, i1, i2, i3}, start: 0, end: 4, expected: sortedSlice},
		{instances: []*Instance{i1, i2, i3, i0}, start: 0, end: 4, expected: sortedSlice},
		{instances: []*Instance{i2, i3, i0, i1}, start: 0, end: 4, expected: sortedSlice},
		{instances: []*Instance{i3, i2, i1, i0}, start: 0, end: 2, expected: []*Instance{i2, i3, i1, i0}},
		{instances: []*Instance{i3, i2, i1, i0}, start: 2, end: 4, expected: []*Instance{i3, i2, i0, i1}},
		{instances: []*Instance{i1, i1, i0, i0}, start: 0, end: 4, expected: []*Instance{i0, i0, i1, i1}},
	}

	for index, test := range tests {
		SortInstancesByRegion(test.instances, test.start, test.end)
		if !reflect.DeepEqual(test.instances, test.expected) {
			t.Fatalf(
				"Instances are not sorted correctly for test %d. Wanted: %+v, got: %+v",
				index,
				test.instances,
				test.expected,
			)
		}
	}
}
