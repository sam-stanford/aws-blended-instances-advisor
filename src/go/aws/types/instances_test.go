package types

import (
	"reflect"
	"testing"
)

func TestSortInstancesByPrice(t *testing.T) {
	i1 := Instance{Name: "i1", PricePerHour: 0.00001}
	i2 := Instance{Name: "i2", PricePerHour: 0.00002}
	i3 := Instance{Name: "i3", PricePerHour: 0.00003}
	i4 := Instance{Name: "i4", PricePerHour: 0.00010}

	slice1 := []Instance{i1, i2, i3, i4}
	slice2 := []Instance{i2, i3, i4, i1}
	slice3 := []Instance{i3, i4, i1, i2}
	slice4 := []Instance{i4, i3, i2, i1}
	slice5 := []Instance{i4, i3, i2, i1}
	slice6 := []Instance{i2, i2, i1, i1}

	sortedSlice := []Instance{i1, i2, i3, i4}
	wantedSlice4 := []Instance{i3, i4, i2, i1}
	wantedSlice5 := []Instance{i4, i3, i1, i2}
	wantedSlice6 := []Instance{i1, i1, i2, i2}

	SortInstancesByPrice(slice1, 0, len(slice1))
	SortInstancesByPrice(slice2, 0, len(slice2))
	SortInstancesByPrice(slice3, 0, len(slice3))
	SortInstancesByPrice(slice4, 0, 2)
	SortInstancesByPrice(slice5, 2, len(slice5))
	SortInstancesByPrice(slice6, 0, len(slice6))

	if !reflect.DeepEqual(slice1, sortedSlice) {
		t.Fatalf("Instances are not sorted correctly for slice1. Wanted: %v, got: %v", sortedSlice, slice1)
	}
	if !reflect.DeepEqual(slice2, sortedSlice) {
		t.Fatalf("Instances are not sorted correctly for slice2. Wanted: %v, got: %v", sortedSlice, slice2)
	}
	if !reflect.DeepEqual(slice3, sortedSlice) {
		t.Fatalf("Instances are not sorted correctly for slice3. Wanted: %v, got: %v", sortedSlice, slice3)
	}
	if !reflect.DeepEqual(slice4, wantedSlice4) {
		t.Fatalf("Instances are not sorted correctly for slice4. Wanted: %v, got: %v", wantedSlice4, slice4)
	}
	if !reflect.DeepEqual(slice5, wantedSlice5) {
		t.Fatalf("Instances are not sorted correctly for slice5. Wanted: %v, got: %v", wantedSlice5, slice5)
	}
	if !reflect.DeepEqual(slice6, wantedSlice6) {
		t.Fatalf("Instances are not sorted correctly for slice6. Wanted: %v, got: %v", wantedSlice6, slice6)
	}
}

func TestSortInstancesByMemory(t *testing.T) {
	i1 := Instance{Name: "i1", MemoryGb: 0.2}
	i2 := Instance{Name: "i2", MemoryGb: 4.6}
	i3 := Instance{Name: "i3", MemoryGb: 64}
	i4 := Instance{Name: "i4", MemoryGb: 128.0}

	slice1 := []Instance{i1, i2, i3, i4}
	slice2 := []Instance{i2, i3, i4, i1}
	slice3 := []Instance{i3, i4, i1, i2}
	slice4 := []Instance{i4, i3, i2, i1}
	slice5 := []Instance{i4, i3, i2, i1}
	slice6 := []Instance{i2, i2, i1, i1}

	sortedSlice := []Instance{i1, i2, i3, i4}
	wantedSlice4 := []Instance{i3, i4, i2, i1}
	wantedSlice5 := []Instance{i4, i3, i1, i2}
	wantedSlice6 := []Instance{i1, i1, i2, i2}

	SortInstancesByMemory(slice1, 0, len(slice1))
	SortInstancesByMemory(slice2, 0, len(slice2))
	SortInstancesByMemory(slice3, 0, len(slice3))
	SortInstancesByMemory(slice4, 0, 2)
	SortInstancesByMemory(slice5, 2, len(slice5))
	SortInstancesByMemory(slice6, 0, len(slice6))

	if !reflect.DeepEqual(slice1, sortedSlice) {
		t.Fatalf("Instances are not sorted correctly for slice1. Wanted: %v, got: %v", sortedSlice, slice1)
	}
	if !reflect.DeepEqual(slice2, sortedSlice) {
		t.Fatalf("Instances are not sorted correctly for slice2. Wanted: %v, got: %v", sortedSlice, slice2)
	}
	if !reflect.DeepEqual(slice3, sortedSlice) {
		t.Fatalf("Instances are not sorted correctly for slice3. Wanted: %v, got: %v", sortedSlice, slice3)
	}
	if !reflect.DeepEqual(slice4, wantedSlice4) {
		t.Fatalf("Instances are not sorted correctly for slice4. Wanted: %v, got: %v", wantedSlice4, slice4)
	}
	if !reflect.DeepEqual(slice5, wantedSlice5) {
		t.Fatalf("Instances are not sorted correctly for slice5. Wanted: %v, got: %v", wantedSlice5, slice5)
	}
	if !reflect.DeepEqual(slice6, wantedSlice6) {
		t.Fatalf("Instances are not sorted correctly for slice6. Wanted: %v, got: %v", wantedSlice6, slice6)
	}
}

func TestSortInstancesByVcpus(t *testing.T) {
	i1 := Instance{Name: "i1", Vcpus: 1}
	i2 := Instance{Name: "i2", Vcpus: 4}
	i3 := Instance{Name: "i3", Vcpus: 16}
	i4 := Instance{Name: "i4", Vcpus: 64}

	slice1 := []Instance{i1, i2, i3, i4}
	slice2 := []Instance{i2, i3, i4, i1}
	slice3 := []Instance{i3, i4, i1, i2}
	slice4 := []Instance{i4, i3, i2, i1}
	slice5 := []Instance{i4, i3, i2, i1}
	slice6 := []Instance{i2, i2, i1, i1}

	sortedSlice := []Instance{i1, i2, i3, i4}
	wantedSlice4 := []Instance{i3, i4, i2, i1}
	wantedSlice5 := []Instance{i4, i3, i1, i2}
	wantedSlice6 := []Instance{i1, i1, i2, i2}

	SortInstancesByVcpus(slice1, 0, len(slice1))
	SortInstancesByVcpus(slice2, 0, len(slice2))
	SortInstancesByVcpus(slice3, 0, len(slice3))
	SortInstancesByVcpus(slice4, 0, 2)
	SortInstancesByVcpus(slice5, 2, len(slice5))
	SortInstancesByVcpus(slice6, 0, len(slice6))

	if !reflect.DeepEqual(slice1, sortedSlice) {
		t.Fatalf("Instances are not sorted correctly for slice1. Wanted: %v, got: %v", sortedSlice, slice1)
	}
	if !reflect.DeepEqual(slice2, sortedSlice) {
		t.Fatalf("Instances are not sorted correctly for slice2. Wanted: %v, got: %v", sortedSlice, slice2)
	}
	if !reflect.DeepEqual(slice3, sortedSlice) {
		t.Fatalf("Instances are not sorted correctly for slice3. Wanted: %v, got: %v", sortedSlice, slice3)
	}
	if !reflect.DeepEqual(slice4, wantedSlice4) {
		t.Fatalf("Instances are not sorted correctly for slice4. Wanted: %v, got: %v", wantedSlice4, slice4)
	}
	if !reflect.DeepEqual(slice5, wantedSlice5) {
		t.Fatalf("Instances are not sorted correctly for slice5. Wanted: %v, got: %v", wantedSlice5, slice5)
	}
	if !reflect.DeepEqual(slice6, wantedSlice6) {
		t.Fatalf("Instances are not sorted correctly for slice6. Wanted: %v, got: %v", wantedSlice6, slice6)
	}
}
