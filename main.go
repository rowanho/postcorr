package main

import (
    "postCorr/getClusters"
    "fmt"
    "os"
)

func main(){
    filename := os.Args[1]
    text := getClusters.ReadFile(filename)
    fmt.Println(text)
}
