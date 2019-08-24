package main

import (
  "github.com/rowanho/msa/get_clusters"
  "fmt"
  "os"
)

func main(){
  filename := os.Args[1]
  text := readFile()
  fmt.printLn(text)
}
