package main

import (
  "fmt"
  "log"
  "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:]) //slicing drops the leading /
}

func main() {
  http.HandleFunc("/", handler) // http package handles requests to web root
  log.Fatal(http.ListenAndServe(":8080",nil)) // wrapped with log.Fatal. blocking
}
