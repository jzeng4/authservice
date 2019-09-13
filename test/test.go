package main

import (
        "fmt"
        "log"
        "io/ioutil"

        "gopkg.in/yaml.v2"
)

type User struct {
    name string
    secret string
    auth []string
}

type Config struct {
    users []User `users`
}

func main() { 
/*
  var config Config
  filename := "secrete.yaml"
  source, err := ioutil.ReadFile(filename)
  
  if err != nil {
      panic(err)
  }
  
  err2 := yaml.Unmarshal(source, &config)
  if err2 != nil {
      log.Fatalf("error: %v", err2)
  }
  fmt.Printf("--- config:\n%s%v\n\n", source, config)
  */
}
