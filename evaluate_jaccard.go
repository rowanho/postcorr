package main

import (
    "postCorr/alignment"
    "postCorr/flags"
    "postCorr/fingerprinting"
    "postCorr/readWrite"
    "postCorr/common"
    "fmt"
    "math"
    
    "sort"
)


func posInt(slice []int, value int) int {
    for p, v := range slice {
        if (v == value) {
            return p
        }
    }
    return -1
}

func  posFl(slice []float64, value float64) int {
    for p, v := range slice {
        if (v == value) {
            return p
        }
    }
    return -1
}

func aligns() (map[int]map[int]int, []int){
    docList, _ := readWrite.TraverseDocs()


    alScores := make(map[int]map[int]int)
    alScoreList := make([]int, 0)
    for i := range docList {
        alScores[i] = make(map[int]int)
        for j := i; j < len(docList); j++ {
            if i == j {
                continue;
            }
            fmt.Println("Aligning")
            score, _, _ := alignment.Gotoh(2, 5, 1, docList[i].Text, docList[j].Text)
            alScores[i][j] = score
            alScoreList = append(alScoreList, score)
        }
    }
    return alScores, alScoreList
}


func getMSE(alScores map[int]map[int]int, alScoreList []int) float64 {
    docList, _ := readWrite.TraverseDocs()
    documentScores, docScoreList := fingerprinting.GetAllPairwise(docList)
    fingerprinting.ResetRuntime()
    fmt.Println(docScoreList)
    sort.Float64s(docScoreList)    
    
    distSum := 0.0
    for i := range docList {
        for j := i; j < len(docList); j++ {
            if  i == j {
                continue;
            }
            
            posFP := posFl(docScoreList, documentScores[i][j])
            posAl := posInt(alScoreList, alScores[i][j])
            distSum += math.Pow(float64(posFP) - float64(posAl), 2)
        }
    }
    l := len(docList)
    return distSum / float64(l*l - l)
}

func EvaluateJaccard() {
    flags.ShingleSize = 5
    flags.NumAligns = 1
    flags.FpType = common.ModFP
    flags.DirName = "synthetic_data/benchmark_align/500_chars/err"
    flags.P = 1
    flags.Affine = true

    flags.JaccardType = common.Jaccard
    als, alList := aligns()
    sort.Ints(alList)
    fmt.Println(alList)
    mse := getMSE(als, alList)
    fmt.Printf("MSE non weighted %5.2f\n", mse)
    flags.JaccardType = common.WeightedJaccard
    // Do the same for the weighted jaccard
    mseWeighted := getMSE(als, alList)

    
    fmt.Printf("MSE weighted %5.2f \n", mseWeighted)
}

