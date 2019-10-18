package main

import (
    "postCorr/reader"
    "postCorr/fingerprinting"
    "postCorr/alignment"
    "fmt"
    "os"
)

func main(){
    filename := os.Args[1]
    text := getClusters.ReadFile(filename)
    fingerprints := getClusters.Kgram(text, 5, "md5")
    fmt.Println(len(fingerprints))
    fmt.Println(fingerprints[0])
}
