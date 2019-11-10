package main

import (
    "postCorr/reader"
    "postCorr/fingerprinting"
    "fmt"
)

func getFingerprints(){
    filename := "test.txt"
    text := reader.ReadFile(filename)
    fingerprints := fingerprinting.ModP(text, 5,4)
    fmt.Println(len(fingerprints))
    fmt.Println(fingerprints[0])

}



func main(){
    getFingerprints()
}
