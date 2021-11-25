package instances

import "testing"

type instanceFloatSearchTest struct {
	instances  []Instance
	value      float64
	start, end int
	want       int
}

type instanceIntSearchTest struct {
	instances  []Instance
	value      int
	start, end int
	want       int
}

func TestFindMinimumMemory(t *testing.T) {
	i0 := Instance{Name: "0", MemoryGb: 4}
	i1 := Instance{Name: "1", MemoryGb: 8}
	i2 := Instance{Name: "2", MemoryGb: 16}
	i3 := Instance{Name: "3", MemoryGb: 62.4}

	tests := map[string]instanceFloatSearchTest{
		"equals value, singleton slice": {instances: []Instance{i0}, value: 4, start: 0, end: 1, want: 0},

		"equals value, sorted slice":                     {instances: []Instance{i0, i1, i2, i3}, value: 8, start: 0, end: 4, want: 1},
		"between values, sorted slice":                   {instances: []Instance{i0, i1, i2, i3}, value: 10, start: 0, end: 4, want: 2},
		"greater than all values, sorted subslice start": {instances: []Instance{i0, i1, i2, i3}, value: 100, start: 0, end: 2, want: 1},
		"equals value, sorted slice, duplicates":         {instances: []Instance{i0, i0, i1, i1}, value: 8, start: 0, end: 4, want: 2},

		"equals, sorted subslice start":             {instances: []Instance{i0, i1, i2, i3}, value: 4, start: 0, end: 2, want: 0},
		"less than all values, sorted subslice end": {instances: []Instance{i0, i1, i2, i3}, value: 1, start: 2, end: 4, want: 2},
		"equals value, sorted subslice middle":      {instances: []Instance{i0, i1, i2, i3}, value: 16, start: 1, end: 3, want: 2},

		"equals value, unsorted slice":                     {instances: []Instance{i1, i0, i3, i2}, value: 8, start: 0, end: 4, want: 0},
		"between values, unsorted slice":                   {instances: []Instance{i0, i1, i3, i2}, value: 10, start: 0, end: 4, want: 3},
		"greater than all values, unsorted subslice start": {instances: []Instance{i3, i2, i1, i0}, value: 100, start: 0, end: 2, want: 0},
		"equals value, unsorted subslice start":            {instances: []Instance{i1, i0, i3, i2}, value: 4, start: 0, end: 2, want: 1},
		"less than all values, unsorted subslice end":      {instances: []Instance{i1, i0, i3, i2}, value: 1, start: 2, end: 4, want: 3},
		"equals value, unsorted subslice middle":           {instances: []Instance{i1, i0, i3, i2}, value: 16, start: 1, end: 3, want: 2},
	}

	for name, test := range tests {
		foundIndex, err := FindMinimumMemory(test.instances, test.value, test.start, test.end)
		if err != nil {
			t.Fatalf(
				"Error occurred when searching for test \"%s\": %s",
				name,
				err.Error(),
			)
		}
		if foundIndex != test.want {
			t.Fatalf(
				"Incorrect index found for test \"%s\". Wanted: %d, got: %d, test: %+v",
				name,
				test.want,
				foundIndex,
				test,
			)
		}
	}

	errorThrowingTests := map[string]instanceFloatSearchTest{
		"zero size slice":             {instances: []Instance{}, value: 0, start: 0, end: 1, want: -1},
		"subslice of zero elements":   {instances: []Instance{i0, i0, i0, i1}, value: 0, start: 1, end: 1, want: -1},
		"start less than 0":           {instances: []Instance{i0, i0, i0, i1}, value: 0, start: -1, end: 3, want: -1},
		"end greater than slice size": {instances: []Instance{i0, i0, i0, i1}, value: 0, start: 1, end: 5, want: -1},
		"start after end":             {instances: []Instance{i0, i0, i0, i1}, value: 0, start: 3, end: 1, want: -1},
	}

	for name, test := range errorThrowingTests {
		_, err := FindMinimumMemory(test.instances, test.value, test.start, test.end)
		if err == nil {
			t.Fatalf(
				"Test did not return error: \"%s\"",
				name,
			)
		}
	}
}

func TestFindMinimumMemorySortedInstances(t *testing.T) {
	i0 := Instance{Name: "0", MemoryGb: 4}
	i1 := Instance{Name: "1", MemoryGb: 8}
	i2 := Instance{Name: "2", MemoryGb: 16}
	i3 := Instance{Name: "3", MemoryGb: 62.4}

	sortedInstances := []Instance{i0, i1, i2, i3}

	tests := map[string]instanceFloatSearchTest{
		"equals value, singleton slice":                  {instances: []Instance{i0}, value: 4, start: 0, end: 1, want: 0},
		"equals value, sorted slice":                     {instances: sortedInstances, value: 8, start: 0, end: 4, want: 1},
		"between values, sorted slice":                   {instances: sortedInstances, value: 10, start: 0, end: 4, want: 2},
		"greater than all values, sorted subslice start": {instances: sortedInstances, value: 100, start: 0, end: 2, want: 1},
		"equals value, sorted subslice start":            {instances: sortedInstances, value: 4, start: 0, end: 2, want: 0},
		"less than all values, sorted subslice end":      {instances: sortedInstances, value: 1, start: 2, end: 4, want: 2},
		"equals value, sorted subslice middle":           {instances: sortedInstances, value: 16, start: 1, end: 3, want: 2},
		"equals value, sorted slice, duplicates":         {instances: []Instance{i0, i0, i1, i1}, value: 8, start: 0, end: 4, want: 2},
	}

	for name, test := range tests {
		foundIndex, err := FindMinimumMemorySortedInstances(test.instances, test.value, test.start, test.end)
		if err != nil {
			t.Fatalf(
				"Error occurred when searching for test \"%s\": %s",
				name,
				err.Error(),
			)
		}
		if foundIndex != test.want {
			t.Fatalf(
				"Incorrect index found for test \"%s\". Wanted: %d, got: %d, test: %+v",
				name,
				test.want,
				foundIndex,
				test,
			)
		}
	}

	errorThrowingTests := map[string]instanceFloatSearchTest{
		"zero size slice":             {instances: []Instance{}, value: 0, start: 0, end: 1, want: -1},
		"subslice of zero elements":   {instances: []Instance{i0, i0, i0, i1}, value: 0, start: 1, end: 1, want: -1},
		"start less than 0":           {instances: []Instance{i0, i0, i0, i1}, value: 0, start: -1, end: 3, want: -1},
		"end greater than slice size": {instances: []Instance{i0, i0, i0, i1}, value: 0, start: 1, end: 5, want: -1},
		"start after end":             {instances: []Instance{i0, i0, i0, i1}, value: 0, start: 3, end: 1, want: -1},
		"unsorted instances":          {instances: []Instance{i3, i2, i1, i0}, value: 0, start: 0, end: 4, want: -1},
	}

	for name, test := range errorThrowingTests {
		_, err := FindMinimumMemorySortedInstances(test.instances, test.value, test.start, test.end)
		if err == nil {
			t.Fatalf(
				"Test did not return error: \"%s\"",
				name,
			)
		}
	}
}

func TestFindMinimumPrice(t *testing.T) {
	i0 := Instance{Name: "0", PricePerHour: 0.001}
	i1 := Instance{Name: "1", PricePerHour: 0.005}
	i2 := Instance{Name: "2", PricePerHour: 0.01}
	i3 := Instance{Name: "3", PricePerHour: 0.05}

	tests := map[string]instanceFloatSearchTest{
		"equals value, singleton slice": {instances: []Instance{i0}, value: 0.001, start: 0, end: 1, want: 0},

		"equals value, sorted slice":                     {instances: []Instance{i0, i1, i2, i3}, value: 0.005, start: 0, end: 4, want: 1},
		"between values, sorted slice":                   {instances: []Instance{i0, i1, i2, i3}, value: 0.009, start: 0, end: 4, want: 2},
		"greater than all values, sorted subslice start": {instances: []Instance{i0, i1, i2, i3}, value: 10.568, start: 0, end: 2, want: 1},
		"equals value, sorted slice, duplicates":         {instances: []Instance{i0, i0, i1, i1}, value: 0.005, start: 0, end: 4, want: 2},

		"equals, sorted subslice start":             {instances: []Instance{i0, i1, i2, i3}, value: 0.001, start: 0, end: 2, want: 0},
		"less than all values, sorted subslice end": {instances: []Instance{i0, i1, i2, i3}, value: 0, start: 2, end: 4, want: 2},
		"equals value, sorted subslice middle":      {instances: []Instance{i0, i1, i2, i3}, value: 0.01, start: 1, end: 3, want: 2},

		"equals value, unsorted slice":                     {instances: []Instance{i1, i0, i3, i2}, value: 0.005, start: 0, end: 4, want: 0},
		"between values, unsorted slice":                   {instances: []Instance{i0, i1, i3, i2}, value: 0.009, start: 0, end: 4, want: 3},
		"greater than all values, unsorted subslice start": {instances: []Instance{i3, i2, i1, i0}, value: 10.568, start: 0, end: 2, want: 0},
		"equals value, unsorted subslice start":            {instances: []Instance{i1, i0, i3, i2}, value: 0.001, start: 0, end: 2, want: 1},
		"less than all values, unsorted subslice end":      {instances: []Instance{i1, i0, i3, i2}, value: 0, start: 2, end: 4, want: 3},
		"equals value, unsorted subslice middle":           {instances: []Instance{i1, i0, i3, i2}, value: 0.01, start: 1, end: 3, want: 2},
	}

	for name, test := range tests {
		foundIndex, err := FindMinimumPrice(test.instances, test.value, test.start, test.end)
		if err != nil {
			t.Fatalf(
				"Error occurred when searching for test \"%s\": %s",
				name,
				err.Error(),
			)
		}
		if foundIndex != test.want {
			t.Fatalf(
				"Incorrect index found for test \"%s\". Wanted: %d, got: %d, test: %+v",
				name,
				test.want,
				foundIndex,
				test,
			)
		}
	}

	errorThrowingTests := map[string]instanceFloatSearchTest{
		"zero size slice":             {instances: []Instance{}, value: 0, start: 0, end: 1, want: -1},
		"subslice of zero elements":   {instances: []Instance{i0, i0, i0, i1}, value: 0, start: 1, end: 1, want: -1},
		"start less than 0":           {instances: []Instance{i0, i0, i0, i1}, value: 0, start: -1, end: 3, want: -1},
		"end greater than slice size": {instances: []Instance{i0, i0, i0, i1}, value: 0, start: 1, end: 5, want: -1},
		"start after end":             {instances: []Instance{i0, i0, i0, i1}, value: 0, start: 3, end: 1, want: -1},
	}

	for name, test := range errorThrowingTests {
		_, err := FindMinimumPrice(test.instances, test.value, test.start, test.end)
		if err == nil {
			t.Fatalf(
				"Test did not return error: \"%s\"",
				name,
			)
		}
	}
}

func TestFindMinimumPriceSortedInstances(t *testing.T) {
	i0 := Instance{Name: "0", PricePerHour: 0.001}
	i1 := Instance{Name: "1", PricePerHour: 0.005}
	i2 := Instance{Name: "2", PricePerHour: 0.01}
	i3 := Instance{Name: "3", PricePerHour: 0.05}

	sortedInstances := []Instance{i0, i1, i2, i3}

	tests := map[string]instanceFloatSearchTest{
		"equals value, singleton slice":                  {instances: []Instance{i0}, value: 0.001, start: 0, end: 1, want: 0},
		"equals value, sorted slice":                     {instances: sortedInstances, value: 0.005, start: 0, end: 4, want: 1},
		"between values, sorted slice":                   {instances: sortedInstances, value: 0.009, start: 0, end: 4, want: 2},
		"greater than all values, sorted subslice start": {instances: sortedInstances, value: 1000, start: 0, end: 2, want: 1},
		"equals value, sorted subslice start":            {instances: sortedInstances, value: 0.001, start: 0, end: 2, want: 0},
		"less than all values, sorted subslice end":      {instances: sortedInstances, value: 0.00001, start: 2, end: 4, want: 2},
		"equals value, sorted subslice middle":           {instances: sortedInstances, value: 0.01, start: 1, end: 3, want: 2},
		"equals value, sorted slice, duplicates":         {instances: []Instance{i0, i0, i1, i1}, value: 0.005, start: 0, end: 4, want: 2},
	}

	for name, test := range tests {
		foundIndex, err := FindMinimumPriceSortedInstances(test.instances, test.value, test.start, test.end)
		if err != nil {
			t.Fatalf(
				"Error occurred when searching for test \"%s\": %s",
				name,
				err.Error(),
			)
		}
		if foundIndex != test.want {
			t.Fatalf(
				"Incorrect index found for test \"%s\". Wanted: %d, got: %d, test: %+v",
				name,
				test.want,
				foundIndex,
				test,
			)
		}
	}

	errorThrowingTests := map[string]instanceFloatSearchTest{
		"zero size slice":             {instances: []Instance{}, value: 0, start: 0, end: 1, want: -1},
		"subslice of zero elements":   {instances: []Instance{i0, i0, i0, i1}, value: 0, start: 1, end: 1, want: -1},
		"start less than 0":           {instances: []Instance{i0, i0, i0, i1}, value: 0, start: -1, end: 3, want: -1},
		"end greater than slice size": {instances: []Instance{i0, i0, i0, i1}, value: 0, start: 1, end: 5, want: -1},
		"start after end":             {instances: []Instance{i0, i0, i0, i1}, value: 0, start: 3, end: 1, want: -1},
		"unsorted instances":          {instances: []Instance{i3, i2, i1, i0}, value: 0, start: 0, end: 4, want: -1},
	}

	for name, test := range errorThrowingTests {
		_, err := FindMinimumPriceSortedInstances(test.instances, test.value, test.start, test.end)
		if err == nil {
			t.Fatalf(
				"Test did not return error: \"%s\"",
				name,
			)
		}
	}
}

func TestFindMinimumRevocationProbability(t *testing.T) {
	i0 := Instance{Name: "0", RevocationProbability: 0.01}
	i1 := Instance{Name: "1", RevocationProbability: 0.1}
	i2 := Instance{Name: "2", RevocationProbability: 0.2}
	i3 := Instance{Name: "3", RevocationProbability: 0.5}

	tests := map[string]instanceFloatSearchTest{
		"equals value, singleton slice": {instances: []Instance{i0}, value: 0.01, start: 0, end: 1, want: 0},

		"equals value, sorted slice":                     {instances: []Instance{i0, i1, i2, i3}, value: 0.1, start: 0, end: 4, want: 1},
		"between values, sorted slice":                   {instances: []Instance{i0, i1, i2, i3}, value: 0.15, start: 0, end: 4, want: 2},
		"greater than all values, sorted subslice start": {instances: []Instance{i0, i1, i2, i3}, value: 1.1, start: 0, end: 2, want: 1},
		"equals value, sorted slice, duplicates":         {instances: []Instance{i0, i0, i1, i1}, value: 0.1, start: 0, end: 4, want: 2},

		"equals, sorted subslice start":             {instances: []Instance{i0, i1, i2, i3}, value: 0.01, start: 0, end: 2, want: 0},
		"less than all values, sorted subslice end": {instances: []Instance{i0, i1, i2, i3}, value: 0, start: 2, end: 4, want: 2},
		"equals value, sorted subslice middle":      {instances: []Instance{i0, i1, i2, i3}, value: 0.2, start: 1, end: 3, want: 2},

		"equals value, unsorted slice":                     {instances: []Instance{i1, i0, i3, i2}, value: 0.1, start: 0, end: 4, want: 0},
		"between values, unsorted slice":                   {instances: []Instance{i0, i1, i3, i2}, value: 0.15, start: 0, end: 4, want: 3},
		"greater than all values, unsorted subslice start": {instances: []Instance{i3, i2, i1, i0}, value: 1.1, start: 0, end: 2, want: 0},
		"equals value, unsorted subslice start":            {instances: []Instance{i1, i0, i3, i2}, value: 0.01, start: 0, end: 2, want: 1},
		"less than all values, unsorted subslice end":      {instances: []Instance{i1, i0, i3, i2}, value: 0, start: 2, end: 4, want: 3},
		"equals value, unsorted subslice middle":           {instances: []Instance{i1, i0, i3, i2}, value: 0.2, start: 1, end: 3, want: 2},
	}

	for name, test := range tests {
		foundIndex, err := FindMinimumRevocationProbability(test.instances, test.value, test.start, test.end)
		if err != nil {
			t.Fatalf(
				"Error occurred when searching for test \"%s\": %s",
				name,
				err.Error(),
			)
		}
		if foundIndex != test.want {
			t.Fatalf(
				"Incorrect index found for test \"%s\". Wanted: %d, got: %d, test: %+v",
				name,
				test.want,
				foundIndex,
				test,
			)
		}
	}

	errorThrowingTests := map[string]instanceFloatSearchTest{
		"zero size slice":             {instances: []Instance{}, value: 0, start: 0, end: 1, want: -1},
		"subslice of zero elements":   {instances: []Instance{i0, i0, i0, i1}, value: 0, start: 1, end: 1, want: -1},
		"start less than 0":           {instances: []Instance{i0, i0, i0, i1}, value: 0, start: -1, end: 3, want: -1},
		"end greater than slice size": {instances: []Instance{i0, i0, i0, i1}, value: 0, start: 1, end: 5, want: -1},
		"start after end":             {instances: []Instance{i0, i0, i0, i1}, value: 0, start: 3, end: 1, want: -1},
	}

	for name, test := range errorThrowingTests {
		_, err := FindMinimumRevocationProbability(test.instances, test.value, test.start, test.end)
		if err == nil {
			t.Fatalf(
				"Test did not return error: \"%s\"",
				name,
			)
		}
	}
}

func TestFindMinimumRevocationProbabilitySortedInstances(t *testing.T) {
	i0 := Instance{Name: "0", RevocationProbability: 0}
	i1 := Instance{Name: "1", RevocationProbability: 0.01}
	i2 := Instance{Name: "2", RevocationProbability: 0.05}
	i3 := Instance{Name: "3", RevocationProbability: 0.1}

	sortedInstances := []Instance{i0, i1, i2, i3}

	tests := map[string]instanceFloatSearchTest{
		"equals value, singleton slice":                  {instances: []Instance{i0}, value: 0, start: 0, end: 1, want: 0},
		"equals value, sorted slice":                     {instances: sortedInstances, value: 0.01, start: 0, end: 4, want: 1},
		"between values, sorted slice":                   {instances: sortedInstances, value: 0.025, start: 0, end: 4, want: 2},
		"greater than all values, sorted subslice start": {instances: sortedInstances, value: 1.1, start: 0, end: 2, want: 1},
		"equals value, sorted subslice start":            {instances: sortedInstances, value: 0.0, start: 0, end: 2, want: 0},
		"less than all values, sorted subslice end":      {instances: sortedInstances, value: 0.0, start: 2, end: 4, want: 2},
		"equals value, sorted subslice middle":           {instances: sortedInstances, value: 0.05, start: 1, end: 3, want: 2},
		"equals value, sorted slice, duplicates":         {instances: []Instance{i0, i0, i1, i1}, value: 0.01, start: 0, end: 4, want: 2},
	}

	for name, test := range tests {
		foundIndex, err := FindMinimumRevocationProbabilitySortedInstances(test.instances, test.value, test.start, test.end)
		if err != nil {
			t.Fatalf(
				"Error occurred when searching for test \"%s\": %s",
				name,
				err.Error(),
			)
		}
		if foundIndex != test.want {
			t.Fatalf(
				"Incorrect index found for test \"%s\". Wanted: %d, got: %d, test: %+v",
				name,
				test.want,
				foundIndex,
				test,
			)
		}
	}

	errorThrowingTests := map[string]instanceFloatSearchTest{
		"zero size slice":             {instances: []Instance{}, value: 0, start: 0, end: 1, want: -1},
		"subslice of zero elements":   {instances: []Instance{i0, i0, i0, i1}, value: 0, start: 1, end: 1, want: -1},
		"start less than 0":           {instances: []Instance{i0, i0, i0, i1}, value: 0, start: -1, end: 3, want: -1},
		"end greater than slice size": {instances: []Instance{i0, i0, i0, i1}, value: 0, start: 1, end: 5, want: -1},
		"start after end":             {instances: []Instance{i0, i0, i0, i1}, value: 0, start: 3, end: 1, want: -1},
		"unsorted instances":          {instances: []Instance{i3, i2, i1, i0}, value: 0, start: 0, end: 4, want: -1},
	}

	for name, test := range errorThrowingTests {
		_, err := FindMinimumRevocationProbabilitySortedInstances(test.instances, test.value, test.start, test.end)
		if err == nil {
			t.Fatalf(
				"Test did not return error: \"%s\"",
				name,
			)
		}
	}
}

func TestFindMinimumVcpu(t *testing.T) {
	i0 := Instance{Name: "0", Vcpus: 4}
	i1 := Instance{Name: "1", Vcpus: 8}
	i2 := Instance{Name: "2", Vcpus: 16}
	i3 := Instance{Name: "3", Vcpus: 32}

	tests := map[string]instanceIntSearchTest{
		"equals value, singleton slice": {instances: []Instance{i0}, value: 4, start: 0, end: 1, want: 0},

		"equals value, sorted slice":                     {instances: []Instance{i0, i1, i2, i3}, value: 8, start: 0, end: 4, want: 1},
		"between values, sorted slice":                   {instances: []Instance{i0, i1, i2, i3}, value: 10, start: 0, end: 4, want: 2},
		"greater than all values, sorted subslice start": {instances: []Instance{i0, i1, i2, i3}, value: 1000, start: 0, end: 2, want: 1},
		"equals value, sorted slice, duplicates":         {instances: []Instance{i0, i0, i1, i1}, value: 8, start: 0, end: 4, want: 2},

		"equals, sorted subslice start":             {instances: []Instance{i0, i1, i2, i3}, value: 4, start: 0, end: 2, want: 0},
		"less than all values, sorted subslice end": {instances: []Instance{i0, i1, i2, i3}, value: 1, start: 2, end: 4, want: 2},
		"equals value, sorted subslice middle":      {instances: []Instance{i0, i1, i2, i3}, value: 16, start: 1, end: 3, want: 2},

		"equals value, unsorted slice":                     {instances: []Instance{i1, i0, i3, i2}, value: 8, start: 0, end: 4, want: 0},
		"between values, unsorted slice":                   {instances: []Instance{i0, i1, i3, i2}, value: 10, start: 0, end: 4, want: 3},
		"greater than all values, unsorted subslice start": {instances: []Instance{i3, i2, i1, i0}, value: 1000, start: 0, end: 2, want: 0},
		"equals value, unsorted subslice start":            {instances: []Instance{i1, i0, i3, i2}, value: 4, start: 0, end: 2, want: 1},
		"less than all values, unsorted subslice end":      {instances: []Instance{i1, i0, i3, i2}, value: 1, start: 2, end: 4, want: 3},
		"equals value, unsorted subslice middle":           {instances: []Instance{i1, i0, i3, i2}, value: 16, start: 1, end: 3, want: 2},
	}

	for name, test := range tests {
		foundIndex, err := FindMinimumVcpu(test.instances, test.value, test.start, test.end)
		if err != nil {
			t.Fatalf(
				"Error occurred when searching for test \"%s\": %s",
				name,
				err.Error(),
			)
		}
		if foundIndex != test.want {
			t.Fatalf(
				"Incorrect index found for test \"%s\". Wanted: %d, got: %d, test: %+v",
				name,
				test.want,
				foundIndex,
				test,
			)
		}
	}

	errorThrowingTests := map[string]instanceIntSearchTest{
		"zero size slice":             {instances: []Instance{}, value: 0, start: 0, end: 1, want: -1},
		"subslice of zero elements":   {instances: []Instance{i0, i0, i0, i1}, value: 0, start: 1, end: 1, want: -1},
		"start less than 0":           {instances: []Instance{i0, i0, i0, i1}, value: 0, start: -1, end: 3, want: -1},
		"end greater than slice size": {instances: []Instance{i0, i0, i0, i1}, value: 0, start: 1, end: 5, want: -1},
		"start after end":             {instances: []Instance{i0, i0, i0, i1}, value: 0, start: 3, end: 1, want: -1},
	}

	for name, test := range errorThrowingTests {
		_, err := FindMinimumVcpu(test.instances, test.value, test.start, test.end)
		if err == nil {
			t.Fatalf(
				"Test did not return error: \"%s\"",
				name,
			)
		}
	}
}

func TestFindMinimumVcpuSortedInstances(t *testing.T) {
	i0 := Instance{Name: "0", Vcpus: 4}
	i1 := Instance{Name: "1", Vcpus: 8}
	i2 := Instance{Name: "2", Vcpus: 16}
	i3 := Instance{Name: "3", Vcpus: 32}

	sortedInstances := []Instance{i0, i1, i2, i3}

	tests := map[string]instanceIntSearchTest{
		"equals value, singleton slice":                  {instances: []Instance{i0}, value: 4, start: 0, end: 1, want: 0},
		"equals value, sorted slice":                     {instances: sortedInstances, value: 8, start: 0, end: 4, want: 1},
		"between values, sorted slice":                   {instances: sortedInstances, value: 10, start: 0, end: 4, want: 2},
		"greater than all values, sorted subslice start": {instances: sortedInstances, value: 100, start: 0, end: 2, want: 1},
		"equals value, sorted subslice start":            {instances: sortedInstances, value: 4, start: 0, end: 2, want: 0},
		"less than all values, sorted subslice end":      {instances: sortedInstances, value: 1, start: 2, end: 4, want: 2},
		"equals value, sorted subslice middle":           {instances: sortedInstances, value: 16, start: 1, end: 3, want: 2},
		"equals value, sorted slice, duplicates":         {instances: []Instance{i0, i0, i1, i1}, value: 8, start: 0, end: 4, want: 2},
	}

	for name, test := range tests {
		foundIndex, err := FindMinimumVcpuSortedInstances(test.instances, test.value, test.start, test.end)
		if err != nil {
			t.Fatalf(
				"Error occurred when searching for test \"%s\": %s",
				name,
				err.Error(),
			)
		}
		if foundIndex != test.want {
			t.Fatalf(
				"Incorrect index found for test \"%s\". Wanted: %d, got: %d, test: %+v",
				name,
				test.want,
				foundIndex,
				test,
			)
		}
	}

	errorThrowingTests := map[string]instanceIntSearchTest{
		"zero size slice":             {instances: []Instance{}, value: 0, start: 0, end: 1, want: -1},
		"subslice of zero elements":   {instances: []Instance{i0, i0, i0, i1}, value: 0, start: 1, end: 1, want: -1},
		"start less than 0":           {instances: []Instance{i0, i0, i0, i1}, value: 0, start: -1, end: 3, want: -1},
		"end greater than slice size": {instances: []Instance{i0, i0, i0, i1}, value: 0, start: 1, end: 5, want: -1},
		"start after end":             {instances: []Instance{i0, i0, i0, i1}, value: 0, start: 3, end: 1, want: -1},
		"unsorted instances":          {instances: []Instance{i3, i2, i1, i0}, value: 0, start: 0, end: 4, want: -1},
	}

	for name, test := range errorThrowingTests {
		_, err := FindMinimumVcpuSortedInstances(test.instances, test.value, test.start, test.end)
		if err == nil {
			t.Fatalf(
				"Test did not return error: \"%s\"",
				name,
			)
		}
	}
}
