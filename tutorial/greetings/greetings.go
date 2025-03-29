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
  "math/rand"
)

// Hello returns a greeting for the named person.
func Hello(name string) (string, error) {
	// if no name was given, return an error with a message.
  if name == "" {
    return "", errors.New("empty name")
  }
  message := fmt.Sprintf(randomFormat(), name)
  return message, nil
}


func Hellos(names []string) (map[string]string, error) {
  messages := make(map[string]string)
  for _, name := range names {
    message, err := Hello(name)
    if err != nil {
      return nil, err
    }
    messages[name] = message
  }
  return messages, nil
}

func randomFormat() string {
  formats := []string{
    "Hi, %v. Welcome!",
    "Great to see you, %v!",
    "Hail, %v! Well met!",
  }

  return formats[rand.Intn(len(formats))]
}
