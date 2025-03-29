package main

import (
	"fmt"
  "log"
	"example.com/greetings"
)

func main() {
  log.SetPrefix("greetings: ")
  log.SetFlags(0)

	message, err := greetings.Hello("")
	fmt.Println(message)
  if err != nil {
    log.Fatal(err)
  }
	//fmt.Println("Hello, World!")
	fmt.Println(message)
}
