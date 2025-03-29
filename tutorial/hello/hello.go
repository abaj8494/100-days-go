package main

import (
	"fmt"
	"example.com/greetings"
)

func main() {
	message := greetings.Hello("Gladys")
	fmt.Println(message)
	//fmt.Println("Hello, World!")
	//fmt.Println(quote.Go())
}
