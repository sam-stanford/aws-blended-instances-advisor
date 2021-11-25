package instances

type instanceFloatSearchTest struct {
	name          string
	instances     []Instance
	searchValue   float64
	start, end    int
	expectedIndex int
}

// TODO: Name proper
// func TestGetIndexOfMinimumMemoryFromUnsortedInstances(t *testing.T) {

// }

// func TestGetIndexOfMinimumMemoryFromSortedInstances(t *testing.T) {
// 	i0 := Instance{Name: "0", MemoryGb: 4}
// 	i1 := Instance{Name: "1", MemoryGb: 8}
// 	i2 := Instance{Name: "2", MemoryGb: 16}
// 	i3 := Instance{Name: "3", MemoryGb: 62.4}

// 	sortedInstances := []Instance{i0, i1, i2, i3}

// 	tests := []instanceFloatSearchTest{
// 		{name: "0", instances: sortedInstances, start: 0, end: 4, searchValue: 3.2, expectedIndex: 0},
// 		{name: "1", instances: sortedInstances, start: 0, end: 4, searchValue: 4, expectedIndex: 0},
// 		{name: "2", instances: sortedInstances, start: 0, end: 4, searchValue: 4.1, expectedIndex: 1},
// 		{name: "3", instances: sortedInstances, start: 0, end: 4, searchValue: 7, expectedIndex: 1},

// 		{name: "4", instances: sortedInstances, start: 0, end: 4, searchValue: 8, expectedIndex: 1},
// 		{name: "5", instances: sortedInstances, start: 0, end: 4, searchValue: 8.01, expectedIndex: 2},
// 		{name: "6", instances: sortedInstances, start: 0, end: 4, searchValue: 62.5, expectedIndex: 3},
// 		{name: "7", instances: sortedInstances, start: 0, end: 4, searchValue: 63, expectedIndex: 3},
// 		{name: "8", instances: sortedInstances, start: 0, end: 4, searchValue: 80, expectedIndex: 3},

// 		{name: "9", instances: sortedInstances, start: 0, end: 2, searchValue: 3, expectedIndex: 0},
// 		{name: "10", instances: sortedInstances, start: 0, end: 2, searchValue: 4, expectedIndex: 0},
// 		{name: "11", instances: sortedInstances, start: 0, end: 2, searchValue: 5, expectedIndex: 1},
// 		{name: "12", instances: sortedInstances, start: 0, end: 2, searchValue: 15, expectedIndex: 1},

// 		{name: "13", instances: sortedInstances, start: 0, end: 3, searchValue: 3, expectedIndex: 0},
// 		{name: "14", instances: sortedInstances, start: 0, end: 3, searchValue: 4, expectedIndex: 0},
// 		{name: "15", instances: sortedInstances, start: 0, end: 3, searchValue: 5, expectedIndex: 1},
// 		{name: "16", instances: sortedInstances, start: 0, end: 3, searchValue: 15.1, expectedIndex: 2},
// 		{name: "17", instances: sortedInstances, start: 0, end: 3, searchValue: 16, expectedIndex: 2},
// 		{name: "18", instances: sortedInstances, start: 0, end: 3, searchValue: 16.2, expectedIndex: 2},
// 		{name: "19", instances: sortedInstances, start: 0, end: 3, searchValue: 128, expectedIndex: 2},

// 		{name: "20", instances: sortedInstances, start: 2, end: 4, searchValue: 1, expectedIndex: 2},
// 		{name: "21", instances: sortedInstances, start: 2, end: 4, searchValue: 16, expectedIndex: 2},
// 		{name: "22", instances: sortedInstances, start: 2, end: 4, searchValue: 16.1, expectedIndex: 3},
// 		{name: "23", instances: sortedInstances, start: 2, end: 4, searchValue: 1000, expectedIndex: 3},

// 		{name: "24", instances: sortedInstances, start: 0, end: 1, searchValue: 3, expectedIndex: 0},
// 		{name: "25", instances: sortedInstances, start: 0, end: 1, searchValue: 4, expectedIndex: 0},
// 		{name: "26", instances: sortedInstances, start: 0, end: 1, searchValue: 5, expectedIndex: 0},

// 		{name: "27", instances: sortedInstances, start: 3, end: 4, searchValue: 1, expectedIndex: 3},
// 		{name: "28", instances: sortedInstances, start: 3, end: 4, searchValue: 16, expectedIndex: 3},
// 		{name: "29", instances: sortedInstances, start: 3, end: 4, searchValue: 16.1, expectedIndex: 3},
// 		{name: "30", instances: sortedInstances, start: 3, end: 4, searchValue: 62.4, expectedIndex: 3},
// 		{name: "31", instances: sortedInstances, start: 3, end: 4, searchValue: 1000, expectedIndex: 3},

// 		{name: "32", instances: []Instance{i0}, start: 0, end: 1, searchValue: 1, expectedIndex: 0},
// 		{name: "33", instances: []Instance{i0}, start: 0, end: 1, searchValue: 1000, expectedIndex: 0},

// 		{name: "34", instances: []Instance{i0, i0, i0, i0}, start: 0, end: 4, searchValue: 1, expectedIndex: 0},
// 		{name: "35", instances: []Instance{i0, i0, i0, i0}, start: 0, end: 4, searchValue: 4, expectedIndex: 3},
// 		{name: "36", instances: []Instance{i0, i0, i0, i0}, start: 0, end: 4, searchValue: 1000, expectedIndex: 3},
// 	}

// 	for _, test := range tests {
// 		foundIndex, err := GetIndexOfMinimumMemoryFromSortedInstances(test.instances, test.searchValue, test.start, test.end)
// 		if err != nil {
// 			t.Fatalf(
// 				"Error occurred when searching for test %s: %s",
// 				test.name,
// 				err.Error(),
// 			)
// 		}
// 		if foundIndex != test.expectedIndex {
// 			t.Fatalf(
// 				"Incorrect index found for test %s. Wanted: %d, got: %d",
// 				test.name,
// 				test.expectedIndex,
// 				foundIndex,
// 			)
// 		}
// 	}
// }
