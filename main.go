package main

import (
    "postCorr/reader"
    "postCorr/fingerprinting"
    "postCorr/alignment"
    "fmt"
)

func getFingerprints(){
    filename := "test.txt"
    text := reader.ReadFile(filename)
    fingerprints := fingerprinting.ModP(text, 5,4)
    fmt.Println(len(fingerprints))
    fmt.Println(fingerprints[0])

}

func getAlignments(){
    r1 := []rune("gattc")
    r2 := []rune("gattc")
    i, j, score, a := alignment.SmithWaterman(r1,r2,1.0,1.0)
    solution := alignment.ReconstructSolution(r1,r2,i,j,a)
    fmt.Println(string(solution))
    fmt.Println(score)
    
}

func main(){
    getFingerprints()
    getAlignments()
}
