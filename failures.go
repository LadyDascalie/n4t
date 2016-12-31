package main

import "fmt"

type Failures struct {
	Get  uint64
	Copy uint64
}

func (f Failures) String() string {
	return fmt.Sprintf("%b errors while downloading\n%b errors while writing to disk", f.Get, f.Copy)
}
