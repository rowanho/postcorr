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
    fingerprints := fingerprinting.Kgram(text, 5, "md5")
    fmt.Println(len(fingerprints))
    fmt.Println(fingerprints[0])

}

func getAlignments(){
    r1 := []rune("gattcaa")
    r2 := []rune("gatctca")
    i, j, score, a := alignment.SmithWaterman(r1,r2,1.0,1.0)
    solution := alignment.ReconstructSolution(r1,r2,i,j,a)
    fmt.Println(string(solution))
    fmt.Println(score)
    
}

func main(){
    getFingerprints()
    getAlignments()
}
