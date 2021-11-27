package utils

import "testing"

type intSearchTest struct {
	slice []int
	value int
	want  int
}

type floatSearchTest struct {
	slice []float64
	value float64
	want  int
}

type stringSearchTest struct {
	slice []string
	value string
	want  int
}

func TestLinearSearchInt(t *testing.T) {
	singletonSlice := []int{10}
	twoElementSlice := []int{20, 10}                        // Sorted: 10, 20
	standardSlice := []int{9, 8, -3, 4, 6}                  // Sorted: -3, 4, 6, 8, 9
	duplicatesSlice := []int{9, 8, 4, 8, 2, 3, -2, 2, 2, 5} // Sorted: -2, 2, 2, 2, 3, 4, 5, 8, 8, 9
	allDuplicatesSlice := []int{1, 1, 1, 1, 1}

	tests := map[string]intSearchTest{
		"less than value, singleton slice":    {slice: singletonSlice, value: 0, want: 0},
		"equals value, singleton slice":       {slice: singletonSlice, value: 10, want: 0},
		"greater than value, singleton slice": {slice: singletonSlice, value: 20, want: 0},

		"less than all values, slice length 2, no duplicates":    {slice: twoElementSlice, value: 0, want: 1},
		"equals first value, slice length 2, no duplicates":      {slice: twoElementSlice, value: 20, want: 0},
		"equals last value, slice length 2, no duplicates":       {slice: twoElementSlice, value: 10, want: 1},
		"between values, slice length 2, no duplicates":          {slice: twoElementSlice, value: 15, want: 0},
		"greater than all values, slice length 2, no duplicates": {slice: twoElementSlice, value: 100, want: 0},

		"less than all values, slice length 5, no duplicates":    {slice: standardSlice, value: -100, want: 2},
		"equals first value, slice length 5, no duplicates":      {slice: standardSlice, value: 9, want: 0},
		"equals last value, slice length 5, no duplicates":       {slice: standardSlice, value: 6, want: 4},
		"equals middle value, slice length 5, no duplicates":     {slice: standardSlice, value: -3, want: 2},
		"between values, slice length 5, no duplicates":          {slice: standardSlice, value: 7, want: 1},
		"greater than all values, slice length 5, no duplicates": {slice: standardSlice, value: 100, want: 0},

		"less than all values, slice length 10, duplicates":    {slice: duplicatesSlice, value: -100, want: 6},
		"equals first value, slice length 10, duplicates":      {slice: duplicatesSlice, value: 9, want: 0},
		"equals last value, slice length 10, duplicates":       {slice: duplicatesSlice, value: 5, want: 9},
		"equals duplicate value, slice length 10, duplicates":  {slice: duplicatesSlice, value: 2, want: 4},
		"between values, slice length 10, duplicates":          {slice: duplicatesSlice, value: 0, want: 4},
		"greater than all values, slice length 10, duplicates": {slice: duplicatesSlice, value: 1000, want: 0},

		"equals value, slice length 5, all duplcates":            {slice: allDuplicatesSlice, value: 1, want: 0},
		"less than all values, slice length 5, all duplcates":    {slice: allDuplicatesSlice, value: 0, want: 0},
		"greater than all values, slice length 5, all duplcates": {slice: allDuplicatesSlice, value: 2, want: 4},
	}

	for name, test := range tests {
		got, err := LinearSearchInt(test.slice, test.value)
		if err != nil {
			t.Fatalf(
				"Error occurred when searching for test \"%s\": %s",
				name,
				err.Error(),
			)
		}
		if got != test.want {
			t.Fatalf(
				"Incorrect index found for test \"%s\". Wanted: %d, got: %d, test: %+v",
				name,
				test.want,
				got,
				test,
			)
		}
	}

	errorThrowingTests := map[string]intSearchTest{
		"empty slice given": {slice: []int{}, value: 0, want: -1},
	}

	for name, test := range errorThrowingTests {
		_, err := LinearSearchInt(test.slice, test.value)
		if err == nil {
			t.Fatalf(
				"Test did not return error: \"%s\"",
				name,
			)
		}
	}
}

func TestLinearSearchFloat(t *testing.T) {
	singletonSlice := []float64{10}
	twoElementSlice := []float64{20.02, 9.99}                                             // Sorted: 9.99, 20.02
	standardSlice := []float64{9, 8.01, -3.2, 4.99, 6}                                    // Sorted: -3.2, 4.99, 6, 8.01, 9
	duplicatesSlice := []float64{9.00, 8.01, 4.62, 8.01, 2.1, 3.9, -2.22, 2.1, 2.1, 5.01} // Sorted: -2.22, 2.1, 2.1, 2.1, 3.9, 4.62, 5.01, 8.01, 8.01, 9.00
	allDuplicatesSlice := []float64{1.1, 1.1, 1.1, 1.1, 1.1}

	tests := map[string]floatSearchTest{
		"less than value, singleton slice":    {slice: singletonSlice, value: 9.01, want: 0},
		"equals value, singleton slice":       {slice: singletonSlice, value: 10, want: 0},
		"greater than value, singleton slice": {slice: singletonSlice, value: 11.1, want: 0},

		"less than all values, slice length 2, no duplicates":    {slice: twoElementSlice, value: 0, want: 1},
		"equals first value, slice length 2, no duplicates":      {slice: twoElementSlice, value: 20.02, want: 0},
		"equals last value, slice length 2, no duplicates":       {slice: twoElementSlice, value: 9.99, want: 1},
		"between values, slice length 2, no duplicates":          {slice: twoElementSlice, value: 15, want: 0},
		"greater than all values, slice length 2, no duplicates": {slice: twoElementSlice, value: 100, want: 0},

		"equals first value, slice length 4":     {slice: []float64{8, 4, 62.4, 16}, value: 8, want: 0},
		"equals duplicate value, slice length 4": {slice: []float64{0, 0, 0.01, 0.01}, value: 0.01, want: 2},

		"less than all values, slice length 5, no duplicates":    {slice: standardSlice, value: -100, want: 2},
		"equals first value, slice length 5, no duplicates":      {slice: standardSlice, value: 9, want: 0},
		"equals last value, slice length 5, no duplicates":       {slice: standardSlice, value: 6, want: 4},
		"equals middle value, slice length 5, no duplicates":     {slice: standardSlice, value: -3.2, want: 2},
		"between values, slice length 5, no duplicates":          {slice: standardSlice, value: 7.1, want: 1},
		"greater than all values, slice length 5, no duplicates": {slice: standardSlice, value: 100, want: 0},

		"less than all values, slice length 10, duplicates":    {slice: duplicatesSlice, value: -100, want: 6},
		"equals first value, slice length 10, duplicates":      {slice: duplicatesSlice, value: 9, want: 0},
		"equals last value, slice length 10, duplicates":       {slice: duplicatesSlice, value: 5.01, want: 9},
		"equals duplicate value, slice length 10, duplicates":  {slice: duplicatesSlice, value: 2.1, want: 4},
		"between values, slice length 10, duplicates":          {slice: duplicatesSlice, value: 0, want: 4},
		"greater than all values, slice length 10, duplicates": {slice: duplicatesSlice, value: 1000, want: 0},

		"equals value, slice length 5, all duplcates":            {slice: allDuplicatesSlice, value: 1.1, want: 0},
		"less than all values, slice length 5, all duplcates":    {slice: allDuplicatesSlice, value: 0, want: 0},
		"greater than all values, slice length 5, all duplcates": {slice: allDuplicatesSlice, value: 2, want: 4},
	}

	for name, test := range tests {
		got, err := LinearSearchFloat(test.slice, test.value)
		if err != nil {
			t.Fatalf(
				"Error occurred when searching for test \"%s\": %s",
				name,
				err.Error(),
			)
		}
		if got != test.want {
			t.Fatalf(
				"Incorrect index found for test \"%s\". Wanted: %d, got: %d, test: %+v",
				name,
				test.want,
				got,
				test,
			)
		}
	}

	errorThrowingTests := map[string]floatSearchTest{
		"empty slice given": {slice: []float64{}, value: 0, want: -1},
	}

	for name, test := range errorThrowingTests {
		_, err := LinearSearchFloat(test.slice, test.value)
		if err == nil {
			t.Fatalf(
				"Test did not return error: \"%s\"",
				name,
			)
		}
	}
}

func TestLinearSearchString(t *testing.T) {
	singletonSlice := []string{"hello"}
	twoElementSlice := []string{"zebra", "apple"}                                                     // Sorted: "apple", "zebra"
	standardSlice := []string{"potato", "octopus", "apple", "chocolate", "lemon"}                     // Sorted: "apple", "chocolate", "lemon", "octopus", "potato"
	duplicatesSlice := []string{"zebra", "banana", "dolphin", "octopus", "banana", "banana", "apple"} // Sorted: "apple", "banana", "banana", "banana", "dolphin", "octopus", "zebra"
	allDuplicatesSlice := []string{"b", "b", "b", "b", "b"}

	tests := map[string]stringSearchTest{
		"less than value, singleton slice":    {slice: singletonSlice, value: "apple", want: 0},
		"equals value, singleton slice":       {slice: singletonSlice, value: "hello", want: 0},
		"greater than value, singleton slice": {slice: singletonSlice, value: "zebra", want: 0},

		"less than all values, slice length 2, no duplicates":    {slice: twoElementSlice, value: "a", want: 1},
		"equals first value, slice length 2, no duplicates":      {slice: twoElementSlice, value: "zebra", want: 0},
		"equals last value, slice length 2, no duplicates":       {slice: twoElementSlice, value: "apple", want: 1},
		"between values, slice length 2, no duplicates":          {slice: twoElementSlice, value: "potato", want: 0},
		"greater than all values, slice length 2, no duplicates": {slice: twoElementSlice, value: "zzzz", want: 0},

		"less than all values, slice length 5, no duplicates":    {slice: standardSlice, value: "a", want: 2},
		"equals first value, slice length 5, no duplicates":      {slice: standardSlice, value: "potato", want: 0},
		"equals last value, slice length 5, no duplicates":       {slice: standardSlice, value: "lemon", want: 4},
		"equals middle value, slice length 5, no duplicates":     {slice: standardSlice, value: "apple", want: 2},
		"between values, slice length 5, no duplicates":          {slice: standardSlice, value: "night", want: 1},
		"greater than all values, slice length 5, no duplicates": {slice: standardSlice, value: "zzzz", want: 0},

		"less than all values, slice length 7, duplicates":    {slice: duplicatesSlice, value: "a", want: 6},
		"equals first value, slice length 7, duplicates":      {slice: duplicatesSlice, value: "zebra", want: 0},
		"equals last value, slice length 7, duplicates":       {slice: duplicatesSlice, value: "apple", want: 6},
		"equals duplicate value, slice length 7, duplicates":  {slice: duplicatesSlice, value: "banana", want: 1},
		"between values, slice length 7, duplicates":          {slice: duplicatesSlice, value: "azure", want: 1},
		"greater than all values, slice length 7, duplicates": {slice: duplicatesSlice, value: "zzzz", want: 0},

		"equals value, slice length 5, all duplcates":            {slice: allDuplicatesSlice, value: "b", want: 0},
		"less than all values, slice length 5, all duplcates":    {slice: allDuplicatesSlice, value: "a", want: 0},
		"greater than all values, slice length 5, all duplcates": {slice: allDuplicatesSlice, value: "c", want: 4},
	}

	for name, test := range tests {
		got, err := LinearSearchString(test.slice, test.value)
		if err != nil {
			t.Fatalf(
				"Error occurred when searching for test \"%s\": %s",
				name,
				err.Error(),
			)
		}
		if got != test.want {
			t.Fatalf(
				"Incorrect index found for test \"%s\". Wanted: %d, got: %d, test: %+v",
				name,
				test.want,
				got,
				test,
			)
		}
	}

	errorThrowingTests := map[string]stringSearchTest{
		"empty slice given": {slice: []string{}, value: "hello", want: -1},
	}

	for name, test := range errorThrowingTests {
		_, err := LinearSearchString(test.slice, test.value)
		if err == nil {
			t.Fatalf(
				"Test did not return error: \"%s\"",
				name,
			)
		}
	}
}

func TestBinarySearchInt(t *testing.T) {
	singletonSlice := []int{10}
	twoElementSlice := []int{10, 20}
	evenLengthSlice := []int{1, 5, 10, 15}
	oddLengthSlice := []int{1, 5, 10, 15, 20}
	largeSlice := []int{-10, -5, -1, 2, 4, 9, 15, 72, 93, 250, 1002, 99999}
	duplicatesSlice := []int{1, 2, 2, 5, 5, 5, 5, 8, 9, 9}
	allDuplicatesSlice := []int{1, 1, 1, 1, 1}

	tests := map[string]intSearchTest{
		"less than value, singleton slice":    {slice: singletonSlice, value: 0, want: 0},
		"equals value, singleton slice":       {slice: singletonSlice, value: 10, want: 0},
		"greater than value, singleton slice": {slice: singletonSlice, value: 20, want: 0},

		"less than all values, even length (2)":         {slice: twoElementSlice, value: 0, want: 0},
		"equals value (index 0), even length (2)":       {slice: twoElementSlice, value: 10, want: 0},
		"between values (index 0 & 1), even length (2)": {slice: twoElementSlice, value: 15, want: 1},
		"equals value (index 1), even length (2)":       {slice: twoElementSlice, value: 20, want: 1},
		"greater than all values, even length (2)":      {slice: twoElementSlice, value: 100, want: 1},

		"less than all values, even length (4)":            {slice: evenLengthSlice, value: 0, want: 0},
		"equals value (index 0), even length (4)":          {slice: evenLengthSlice, value: 1, want: 0},
		"between values (indexes 0 & 1), even length (4)":  {slice: evenLengthSlice, value: 2, want: 1},
		"equals value (index 1), even length (4)":          {slice: evenLengthSlice, value: 5, want: 1},
		"between values (indexes 1 & 2), even length (4)":  {slice: evenLengthSlice, value: 8, want: 2},
		"equals value (index 2), even length (4)":          {slice: evenLengthSlice, value: 10, want: 2},
		"between values (indexes 2 & 3), even length (4) ": {slice: evenLengthSlice, value: 14, want: 3},
		"equals value (index 3), even length (4)":          {slice: evenLengthSlice, value: 15, want: 3},
		"greater than all values, even length (4)":         {slice: evenLengthSlice, value: 100, want: 3},

		"less than all values, odd length (5)":            {slice: oddLengthSlice, value: 0, want: 0},
		"equals value (index 0), odd length (5)":          {slice: oddLengthSlice, value: 1, want: 0},
		"between values (indexes 0 & 1), odd length (5)":  {slice: oddLengthSlice, value: 2, want: 1},
		"equals value (index 1), odd length (5)":          {slice: oddLengthSlice, value: 5, want: 1},
		"between values (indexes 1 & 2), odd length (5)":  {slice: oddLengthSlice, value: 8, want: 2},
		"equals value (index 2), odd length (5)":          {slice: oddLengthSlice, value: 10, want: 2},
		"between values (indexes 2 & 3), odd length (5) ": {slice: oddLengthSlice, value: 14, want: 3},
		"equals value (index 3), odd length (5)":          {slice: oddLengthSlice, value: 15, want: 3},
		"between values (indexes 3 & 4), odd length (5) ": {slice: oddLengthSlice, value: 17, want: 4},
		"equals value (index 4), odd length (5)":          {slice: oddLengthSlice, value: 20, want: 4},
		"greater than all values, odd length (5)":         {slice: oddLengthSlice, value: 100, want: 4},

		"less than all values, large slice":           {slice: largeSlice, value: -1000, want: 0},
		"between values (indexes 1 & 2), large slice": {slice: largeSlice, value: -3, want: 2},
		"between values (indexes 5 & 6), large slice": {slice: largeSlice, value: 11, want: 6},
		"greater than all values, large slice":        {slice: largeSlice, value: 9999999, want: 11},

		"less than all values, duplicates slice":           {slice: duplicatesSlice, value: -1000, want: 0},
		"equals value (indexes 1 & 2), duplicates slice":   {slice: duplicatesSlice, value: 2, want: 1},
		"equals value (indexes 3 to 6), duplicates slice":  {slice: duplicatesSlice, value: 5, want: 3},
		"between values (indexes 2 & 3), duplicates slice": {slice: duplicatesSlice, value: 3, want: 3},
		"greater than all values, duplicates slice":        {slice: duplicatesSlice, value: 9999999, want: 9},

		"equals value, slice length 5, all duplcates":            {slice: allDuplicatesSlice, value: 1, want: 0},
		"less than all values, slice length 5, all duplcates":    {slice: allDuplicatesSlice, value: 0, want: 0},
		"greater than all values, slice length 5, all duplcates": {slice: allDuplicatesSlice, value: 2, want: 4},
	}

	for name, test := range tests {
		got, err := BinarySearchInt(test.slice, test.value)
		if err != nil {
			t.Fatalf(
				"Error occurred when searching for test \"%s\": %s",
				name,
				err.Error(),
			)
		}
		if got != test.want {
			t.Fatalf(
				"Incorrect index found for test \"%s\". Wanted: %d, got: %d, test: %+v",
				name,
				test.want,
				got,
				test,
			)
		}
	}

	errorThrowingTests := map[string]intSearchTest{
		"empty slice given": {slice: []int{}, value: 0, want: -1},
	}

	for name, test := range errorThrowingTests {
		_, err := BinarySearchInt(test.slice, test.value)
		if err == nil {
			t.Fatalf(
				"Test did not return error: \"%s\"",
				name,
			)
		}
	}
}

func TestBinarySearchFloat(t *testing.T) {
	singletonSlice := []float64{3.2}
	twoElementSlice := []float64{3.2, 10}
	evenLengthSlice := []float64{3.2, 8.1, 10, 15.2}
	oddLengthSlice := []float64{3.2, 8.1, 10, 15.2, 20}
	largeSlice := []float64{-10.1, -5, -1.2, 2.04, 4.001, 9, 15.2, 72.001, 93.92, 250.6, 1002, 99999.9}
	duplicatesSlice := []float64{-1, 2.01, 2.01, 5.32, 5.32, 5.32, 5.32, 8, 9.1, 9.1}
	allDuplicatesSlice := []float64{1.1, 1.1, 1.1, 1.1, 1.1}

	tests := map[string]floatSearchTest{
		"less than value, singleton slice":    {slice: singletonSlice, value: 0, want: 0},
		"equals value, singleton slice":       {slice: singletonSlice, value: 3.2, want: 0},
		"greater than value, singleton slice": {slice: singletonSlice, value: 20, want: 0},

		"less than all values, even length (2)":         {slice: twoElementSlice, value: 0, want: 0},
		"equals value (index 0), even length (2)":       {slice: twoElementSlice, value: 3.2, want: 0},
		"between values (index 0 & 1), even length (2)": {slice: twoElementSlice, value: 3.3, want: 1},
		"equals value (index 1), even length (2)":       {slice: twoElementSlice, value: 10, want: 1},
		"greater than all values, even length (2)":      {slice: twoElementSlice, value: 100, want: 1},

		"less than all values, even length (4)":            {slice: evenLengthSlice, value: -1, want: 0},
		"equals value (index 0), even length (4)":          {slice: evenLengthSlice, value: 3.2, want: 0},
		"between values (indexes 0 & 1), even length (4)":  {slice: evenLengthSlice, value: 3.21, want: 1},
		"equals value (index 1), even length (4)":          {slice: evenLengthSlice, value: 8.1, want: 1},
		"between values (indexes 1 & 2), even length (4)":  {slice: evenLengthSlice, value: 9.9, want: 2},
		"equals value (index 2), even length (4)":          {slice: evenLengthSlice, value: 10, want: 2},
		"between values (indexes 2 & 3), even length (4) ": {slice: evenLengthSlice, value: 14, want: 3},
		"equals value (index 3), even length (4)":          {slice: evenLengthSlice, value: 15.2, want: 3},
		"greater than all values, even length (4)":         {slice: evenLengthSlice, value: 100, want: 3},

		"less than all values, odd length (5)":            {slice: oddLengthSlice, value: -1, want: 0},
		"equals value (index 0), odd length (5)":          {slice: oddLengthSlice, value: 3.2, want: 0},
		"between values (indexes 0 & 1), odd length (5)":  {slice: oddLengthSlice, value: 3.21, want: 1},
		"equals value (index 1), odd length (5)":          {slice: oddLengthSlice, value: 8.1, want: 1},
		"between values (indexes 1 & 2), odd length (5)":  {slice: oddLengthSlice, value: 9.9, want: 2},
		"equals value (index 2), odd length (5)":          {slice: oddLengthSlice, value: 10, want: 2},
		"between values (indexes 2 & 3), odd length (5) ": {slice: oddLengthSlice, value: 14, want: 3},
		"equals value (index 3), odd length (5)":          {slice: oddLengthSlice, value: 15.2, want: 3},
		"between values (indexes 3 & 4), odd length (5) ": {slice: oddLengthSlice, value: 19.901, want: 4},
		"equals value (index 4), odd length (5)":          {slice: oddLengthSlice, value: 20, want: 4},
		"greater than all values, odd length (5)":         {slice: oddLengthSlice, value: 100, want: 4},

		"less than all values, large slice":           {slice: largeSlice, value: -1000, want: 0},
		"between values (indexes 1 & 2), large slice": {slice: largeSlice, value: -1.3, want: 2},
		"between values (indexes 5 & 6), large slice": {slice: largeSlice, value: 11, want: 6},
		"greater than all values, large slice":        {slice: largeSlice, value: 9999999, want: 11},

		"less than all values, duplicates slice":           {slice: duplicatesSlice, value: -1000, want: 0},
		"equals value (indexes 1 & 2), duplicates slice":   {slice: duplicatesSlice, value: 2.01, want: 1},
		"equals value (indexes 3 to 6), duplicates slice":  {slice: duplicatesSlice, value: 5.32, want: 3},
		"between values (indexes 2 & 3), duplicates slice": {slice: duplicatesSlice, value: 3, want: 3},
		"greater than all values, duplicates slice":        {slice: duplicatesSlice, value: 9999999, want: 9},

		"equals value, slice length 5, all duplcates":            {slice: allDuplicatesSlice, value: 1.1, want: 0},
		"less than all values, slice length 5, all duplcates":    {slice: allDuplicatesSlice, value: 0, want: 0},
		"greater than all values, slice length 5, all duplcates": {slice: allDuplicatesSlice, value: 2, want: 4},
	}

	for name, test := range tests {
		got, err := BinarySearchFloat(test.slice, test.value)
		if err != nil {
			t.Fatalf(
				"Error occurred when searching for test \"%s\": %s",
				name,
				err.Error(),
			)
		}
		if got != test.want {
			t.Fatalf(
				"Incorrect index found for test \"%s\". Wanted: %d, got: %d, test: %+v",
				name,
				test.want,
				got,
				test,
			)
		}
	}

	errorThrowingTests := map[string]floatSearchTest{
		"empty slice given": {slice: []float64{}, value: 0, want: -1},
	}

	for name, test := range errorThrowingTests {
		_, err := BinarySearchFloat(test.slice, test.value)
		if err == nil {
			t.Fatalf(
				"Test did not return error: \"%s\"",
				name,
			)
		}
	}
}

func TestBinarySearchString(t *testing.T) {
	singletonSlice := []string{"hello"}
	twoElementSlice := []string{"apple", "zebra"}
	evenLengthSlice := []string{"apple", "carrot", "lemon", "money"}
	oddLengthSlice := []string{"apple", "carrot", "lemon", "money", "zebra"}
	largeSlice := []string{"apple", "carrot", "lemon", "money", "octopus", "pumpkin", "red", "staple", "zebra"}
	duplicatesSlice := []string{"apple", "banana", "banana", "dolphin", "dolphin", "dolphin", "dolphin", "hotel", "octopus", "octopus"}
	allDuplicatesSlice := []string{"b", "b", "b", "b", "b"}

	tests := map[string]stringSearchTest{
		"less than value, singleton slice":    {slice: singletonSlice, value: "apple", want: 0},
		"equals value, singleton slice":       {slice: singletonSlice, value: "hello", want: 0},
		"greater than value, singleton slice": {slice: singletonSlice, value: "zebra", want: 0},

		"less than all values, even length (2)":         {slice: twoElementSlice, value: "a", want: 0},
		"equals value (index 0), even length (2)":       {slice: twoElementSlice, value: "apple", want: 0},
		"between values (index 0 & 1), even length (2)": {slice: twoElementSlice, value: "azure", want: 1},
		"equals value (index 1), even length (2)":       {slice: twoElementSlice, value: "zebra", want: 1},
		"greater than all values, even length (2)":      {slice: twoElementSlice, value: "zzz", want: 1},

		"less than all values, even length (4)":            {slice: evenLengthSlice, value: "a", want: 0},
		"equals value (index 0), even length (4)":          {slice: evenLengthSlice, value: "apple", want: 0},
		"between values (indexes 0 & 1), even length (4)":  {slice: evenLengthSlice, value: "azure", want: 1},
		"equals value (index 1), even length (4)":          {slice: evenLengthSlice, value: "carrot", want: 1},
		"between values (indexes 1 & 2), even length (4)":  {slice: evenLengthSlice, value: "dolphin", want: 2},
		"equals value (index 2), even length (4)":          {slice: evenLengthSlice, value: "lemon", want: 2},
		"between values (indexes 2 & 3), even length (4) ": {slice: evenLengthSlice, value: "load", want: 3},
		"equals value (index 3), even length (4)":          {slice: evenLengthSlice, value: "money", want: 3},
		"greater than all values, even length (4)":         {slice: evenLengthSlice, value: "zzz", want: 3},

		"less than all values, odd length (5)":            {slice: oddLengthSlice, value: "a", want: 0},
		"equals value (index 0), odd length (5)":          {slice: oddLengthSlice, value: "apple", want: 0},
		"between values (indexes 0 & 1), odd length (5)":  {slice: oddLengthSlice, value: "azure", want: 1},
		"equals value (index 1), odd length (5)":          {slice: oddLengthSlice, value: "carrot", want: 1},
		"between values (indexes 1 & 2), odd length (5)":  {slice: oddLengthSlice, value: "dolphin", want: 2},
		"equals value (index 2), odd length (5)":          {slice: oddLengthSlice, value: "lemon", want: 2},
		"between values (indexes 2 & 3), odd length (5) ": {slice: oddLengthSlice, value: "load", want: 3},
		"equals value (index 3), odd length (5)":          {slice: oddLengthSlice, value: "money", want: 3},
		"between values (indexes 3 & 4), odd length (5) ": {slice: oddLengthSlice, value: "octopus", want: 4},
		"equals value (index 4), odd length (5)":          {slice: oddLengthSlice, value: "zebra", want: 4},
		"greater than all values, odd length (5)":         {slice: oddLengthSlice, value: "zzz", want: 4},

		"less than all values, large slice":           {slice: largeSlice, value: "a", want: 0},
		"between values (indexes 1 & 2), large slice": {slice: largeSlice, value: "dolphin", want: 2},
		"between values (indexes 5 & 6), large slice": {slice: largeSlice, value: "queen", want: 6},
		"greater than all values, large slice":        {slice: largeSlice, value: "zzz", want: 8},

		"less than all values, duplicates slice":           {slice: duplicatesSlice, value: "a", want: 0},
		"equals value (indexes 1 & 2), duplicates slice":   {slice: duplicatesSlice, value: "banana", want: 1},
		"equals value (indexes 3 to 6), duplicates slice":  {slice: duplicatesSlice, value: "dolphin", want: 3},
		"between values (indexes 2 & 3), duplicates slice": {slice: duplicatesSlice, value: "chocolate", want: 3},
		"greater than all values, duplicates slice":        {slice: duplicatesSlice, value: "zzzzzz", want: 9},

		"equals value, slice length 5, all duplcates":            {slice: allDuplicatesSlice, value: "b", want: 0},
		"less than all values, slice length 5, all duplcates":    {slice: allDuplicatesSlice, value: "a", want: 0},
		"greater than all values, slice length 5, all duplcates": {slice: allDuplicatesSlice, value: "c", want: 4},
	}

	for name, test := range tests {
		got, err := BinarySearchString(test.slice, test.value)
		if err != nil {
			t.Fatalf(
				"Error occurred when searching for test \"%s\": %s",
				name,
				err.Error(),
			)
		}
		if got != test.want {
			t.Fatalf(
				"Incorrect index found for test \"%s\". Wanted: %d, got: %d, test: %+v",
				name,
				test.want,
				got,
				test,
			)
		}
	}

	errorThrowingTests := map[string]stringSearchTest{
		"empty slice given": {slice: []string{}, value: "hello", want: -1},
	}

	for name, test := range errorThrowingTests {
		_, err := BinarySearchString(test.slice, test.value)
		if err == nil {
			t.Fatalf(
				"Test did not return error: \"%s\"",
				name,
			)
		}
	}
}
