package main

import (
    "postCorr/getClusters"
    "fmt"
    "os"
)

func main(){
    filename := os.Args[1]
    text := getClusters.ReadFile(filename)
    fingerprints := getClusters.Kgram(text, 5, "md5")
    fmt.Println(fingerprints)
}
