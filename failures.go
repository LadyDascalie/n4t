package main

import "fmt"

// Failures is used to store how many download errors happened
type Failures struct {
	Get  uint64
	Copy uint64
}

// String stringer implementation for Failures struct
func (f Failures) String() string {
	return fmt.Sprintf("%b errors while downloading\n%b errors while writing to disk", f.Get, f.Copy)
}
