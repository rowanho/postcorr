package main

import (
    "postCorr/alignment"
    "postCorr/flags"
    "postCorr/fingerprinting"
    "postCorr/readWrite"
    "postCorr/common"
    "fmt"
    
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
        for j := i + 1; j < len(docList); j++ {
            fmt.Println("Aligning")
            score, _, _ := alignment.Gotoh(2, 5, 1, docList[i].Text, docList[j].Text)
            alScores[i][j] = score
            alScoreList = append(alScoreList, score)
        }
    }
    return alScores, alScoreList
}


func getTopKPrecision(alScores map[int]map[int]int, alScoreList []int, kProportion float64) float64 {
    cutoff_a := (1.0 - kProportion) * float64((len(alScoreList) - 1))
    topKScore := alScoreList[int(cutoff_a)]
    docList, _ := readWrite.TraverseDocs()
    documentScores, docScoreList := fingerprinting.GetAllPairwise(docList)
    fingerprinting.ResetRuntime()
    sort.Float64s(docScoreList)   
    topFScore := docScoreList[int(cutoff_a)] 
    if topFScore == 0.0 {
        return 0.0
    }  
    fmt.Println(topFScore)
    top := 0
    bot := 0
    for i := range docList {
        for j := i + 1; j < len(docList); j++ {
             if alScores[i][j] > topKScore {
                 if documentScores[i][j] > topFScore {
                     top += 1
                 } else {
                     bot += 1
                 }
            }
        }
    }
    return float64(top) / float64(top + bot)
}

func EvaluateJaccard() {
    downsamples := []int{1, 30, 60, 90, 120, 150, 180}
    flags.NumAligns = 1
    flags.DirName = "real_datasets/copyright/ocr_min"
    flags.Affine = false
    flags.ShingleSize = 5
    flags.JaccardType = common.Jaccard
    als, alList := aligns()
    sort.Ints(alList)
    for _, d := range downsamples  {
        flags.P = d
        flags.FpType = common.ModFP
        fmt.Printf("Downsampling rate %d\n", flags.P)
        fmt.Println("ModP")
        proportion := 0.1
        flags.JaccardType = common.Jaccard
        topk := getTopKPrecision(als, alList, proportion)
        fmt.Printf("Top k precision %5.2f\n", topk)
        // Do the same for the weighted jaccard
        flags.JaccardType = common.WeightedJaccard
        topkWeighted := getTopKPrecision(als, alList, proportion)
        fmt.Printf("Top k precision weighted %5.2f \n", topkWeighted)
        flags.FpType = common.Winnowing
        flags.WinnowingT = d + flags.ShingleSize - 1
        fmt.Println("Winnowing")
        flags.JaccardType = common.Jaccard
        topk = getTopKPrecision(als, alList, proportion)
        fmt.Printf("Top k precision %5.2f\n", topk)
        // Do the same for the weighted jaccard
        flags.JaccardType = common.WeightedJaccard
        topkWeighted = getTopKPrecision(als, alList, proportion)
        fmt.Printf("Top k precision weighted %5.2f \n", topkWeighted)

    }
    
}

