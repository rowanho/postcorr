package main

import (
	"postCorr/fingerprinting"
	"postCorr/reader"

	"fmt"
)

func getFingerprints() {
	filename := "test.txt"
	text := reader.ReadFile(filename)
	fingerprints := fingerprinting.ModP(text, 5, 4)
	fmt.Println(len(fingerprints))
	fmt.Println(fingerprints[0])

}

func main() {
	getFingerprints()
}
