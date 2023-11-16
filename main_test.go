package main

import (
	"fmt"
	"math"
	"testing"
)

func Test_main(t *testing.T) {
	var v float64 = 100
	var a float64 = 200
	var d float64 = 100

	root1 := (-v + math.Sqrt(v*v-4*0.5*a*(-d))) / (2 * 0.5 * a)

	fmt.Println(root1)
}
