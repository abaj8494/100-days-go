package main

import (
  "fmt" 
  "example/user/hello/morestrings"
  "github.com/google/go-cmp/cmp"
)

func main() {
  fmt.Println(morestrings.ReverseRunes("Hello, world."))
  fmt.Println(cmp.Diff("Hello World", "Hello Go"))
}
