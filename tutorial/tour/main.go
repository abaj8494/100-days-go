package main

import (
  "fmt"
  "math/cmplx"
)


// note that the return values are named!
func split(sum int) (x, y int) {
  x = sum * 4 / 9
  y = sum - x
  return // a _naked_ return. should only be used in short functions.
}

func vars() (int, rune, bool, int, bool, string) {
  var c, python, java = true, false, "no!" // note the hetoregeniety of types and lowercase bools.
  // inside functions, the shorthand := is available:
  d := 5
  a, b := 1, 'a'
  return a, b, c, d, python, java
}

func main() {
  // same as import blocks, you can do with vars:
  var (
    ToBe   bool       = false
    MaxInt uint64     = 1<<64 - 1
    z      complex128 = cmplx.Sqrt(-5 + 12i)
  )
	fmt.Printf("Type: %T Value: %v\n", ToBe, ToBe)
	fmt.Printf("Type: %T Value: %v\n", MaxInt, MaxInt)
	fmt.Printf("Type: %T Value: %v\n", z, z)
  // you can do type inferring:
  i := 42           // int
  f := 3.142        // float64
  g := 0.867 + 0.5i // complex128
  fmt.Println(i,f,g)

  // constants cannot be declared using the := syntax
  const World = "Pen15"
  fmt.Println(split(17))
  fmt.Println(vars())
}
