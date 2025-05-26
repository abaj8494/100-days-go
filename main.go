package main

import "math/rand/v2"

func RandomIntGenerator(min, max int) int {
  return rand.IntN(max-min) + min
}

func main() {
  RandomIntGenerator(4,9)
}
