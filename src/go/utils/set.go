package utils

// StringSet is a type which contains no duplicate elements.
type StringSet map[string]bool

// Add adds a string to a StringSet.
func (set StringSet) Add(s string) {
	set[s] = true
}

// Removes removes a string from a StringSet.
func (set StringSet) Remove(s string) {
	set[s] = false
}

// Contains returns a boolean representing whether a
// StringSet contains the given value.
func (set StringSet) Contains(s string) bool {
	val, contains := set[s]
	return contains && val
}
