package main

import "fmt"

type Number interface {
  int64 | float64
}

func main() {
  // Initialise a map for the integer values
  ints := map[string]int64 {
    "first": 34,
    "second": 12,
  }

  // Initialise a map for the float values
  floats := map[string]float64 {
    "first": 35.98,
    "second": 26.99,
  }
  fmt.Printf("Non-Generic Sums: %v and %v\n",
    SumInts(ints),
    SumFloats(floats))

  fmt.Printf("Generic Sums: %v and %v\n",
    SumIntsOrFloats[string, int64](ints), // go can infer these
    SumIntsOrFloats[string, float64](floats))


  fmt.Printf("Generic Sums, type parameters inferred: %v and %v\n",
    SumIntsOrFloats(ints), // go can infer these
    SumIntsOrFloats(floats))
}
// SumInts adds together the values of m.
func SumInts(m map[string]int64) int64 {
  var s int64
  for _, v := range m {
    s += v
  }
  return s
}

// SumFloats adds together the values of m.
func SumFloats(m map[string]float64) float64 {
  var s float64
  for _, v := range m {
    s += v
  }
  return s
}

// the comparable is predefined in Go
// typesafe language (btw)
func SumIntsOrFloats[K comparable, V Number](m map[K]V) V {
  var s V // declaring variable s of type (variable) V as defined in interface^
  for _, v := range m {
    s += v
  }
  return s
}
