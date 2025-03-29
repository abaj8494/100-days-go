/*
Create a module -- Write a small module with functions you can call from another module.
Call your code from another module -- Import and use your new module.
Return and handle an error -- Add simple error handling.
Return a random greeting -- Handle data in slices (Go's dynamically-sized arrays).
Return greetings for multiple people -- Store key/value pairs in a map.
Add a test -- Use Go's built-in unit testing features to test your code.
Compile and install the application -- Compile and install your code locally.
*/

package greetings

import (
  "errors"
  "fmt"
)

// Hello returns a greeting for the named person.
func Hello(name string) (string, error) {
	// if no name was given, return an error with a message.
  if name == "" {
    return "", errors.New("empty name")
  }
  
  

  // if a name was received, return a value that embeds the name
  // in a greeting message.
  message := fmt.Sprintf("Hi, %v. Welcome!", name)
  return message, nil
}
