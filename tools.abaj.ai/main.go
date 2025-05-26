package main

import (
	"fmt"
	"math"
	"math/rand/v2"
)

func AreaCalculator(radius float64) float64 {
  return math.Pow(math.Pi*radius, 2)
}

func RandomIntGenerator(min, max int) int {
  return rand.IntN(max+1-min) + min
}

func main() {
  fmt.Println(RandomIntGenerator(4,9))
  fmt.Printf("Area of circle with radius: %v is: %v", 5, AreaCalculator(5))
}
