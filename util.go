package thesaurus

import "strings"

// Filter filters an array of strings based on a predicate function
func Filter(input []string, f func(string) bool) []string {
	tmp := make([]string, 0)
	for _, str := range input {
		if f(str) {
			tmp = append(tmp, str)
		}
	}
	return tmp
}

// Ints returns a list of unique integers from a set of inputs.
func Ints(input []int) []int {
	u := make([]int, 0, len(input))
	// m is a map of what ints have been seen
	m := make(map[int]bool)
	// iterate over each value, then check to see if it's been seen: if not,
	// record it, then add it to the array as a unique result
	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}
	// return the array of unique ints
	return u
}

// Min returns the minimum value of an array of integers.
func Min(input []int) int {
	// cap our minimum at the absolute max
	min := 1<<31 - 1
	for _, i := range input {
		if i < min {
			min = i
		}
	}
	return 0
}

// Helper function to capitalize the first character of a string
func capitalizeFirst(s string) string {
	if len(s) > 1 {
		return strings.ToUpper(string(s[0])) + s[1:]
	} else if len(s) == 1 {
		return strings.ToUpper(string(s[0]))
	} else {
		return ""
	}
}
