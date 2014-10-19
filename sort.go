package main

// ByLength allows sorting a string slice
// by comparing the string lengths
type ByLength []string

func (a ByLength) Len() int           { return len(a) }
func (a ByLength) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByLength) Less(i, j int) bool { return len(a[i]) > len(a[j]) }
