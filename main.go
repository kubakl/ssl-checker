package main

import (
  "fmt"
  "flag"
  _ "os"
)

var (
  email = flag.String("e", "", "Email address for sending expiry date of provided domains")
  domain = flag.String("d", "", "Single domain to be checked, in case you don't have a file with domain names")
  filename = flag.String("f", "", "Route to a file with domains to check")
)

func main() {

  flag.Parse()
  if *filename == "" && *email == ""{
  } else {
    fmt.Println("Provided filename:", *filename)  
    fmt.Println("Provided email:", *email)  
  }
}
