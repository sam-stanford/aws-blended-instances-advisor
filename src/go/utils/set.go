package utils

// TODO: Doc & test

type StringSet map[string]bool

func (set StringSet) Add(s string) {
	set[s] = true
}

func (set StringSet) Remove(s string) {
	set[s] = false
}

func (set StringSet) Contains(s string) bool {
	val, contains := set[s]
	return contains && val
}
