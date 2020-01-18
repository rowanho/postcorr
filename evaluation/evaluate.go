package evaluation

import (
    "postCorr/readWrite"
    "postCorr/flags"
    "postCorr/common"
    
    "os"
    "path/filepath"
    "fmt"
    
    "github.com/agnivade/levenshtein"
)


/**
* Traverses the input and output directories, and sums the edit distances of each file
* Between the data and the ground truth data
**/
func editDistances(docs []common.Document, docMap map[string]int, correctedDocs map[string]bool) ([]string, []int, []int, error) {
    originalDistances := make([]int, 0)
    correctedDistances := make([]int, 0)
    docIds := make([]string, 0)
	err := filepath.Walk(flags.DirName,
		func(path string, info os.FileInfo, err error) error {
            right := ""
            if path != flags.DirName {
                right = path[len(flags.DirName) + 1:]
            }
			if err != nil {
				return err
			}
			if info.IsDir() == false {
                _, correctable := correctedDocs[right]

                var groundTruth string
                var corrected string
                var original string
                var readErr error
                
				if flags.FormatType == common.Plaintext {
					original, readErr = readWrite.ReadString(path)
                    groundTruth, readErr = readWrite.ReadString(flags.Groundtruth + "/" + right)
                    if correctable {
                        corrected = string(docs[docMap[right]].Text)
                    }
				}
                
				if readErr != nil {
                    fmt.Println(readErr)
					return readErr
				}
                
                originalDist := levenshtein.ComputeDistance(original, groundTruth)
                originalDistances = append(originalDistances, originalDist)
                
                if correctable {
                    correctedDist := levenshtein.ComputeDistance(corrected, groundTruth)
                    correctedDistances = append(correctedDistances, correctedDist)
                } else {
                    correctedDistances = append(correctedDistances, originalDist)                    
                }
                
                docIds = append(docIds, right)
			}
			return nil
		},
	)
    
    if err != nil {
        return []string{}, []int{}, []int{}, err
    }
	return  docIds, originalDistances, correctedDistances, err
}

type EvalStats = struct {
    Mean float64
    Total int
    MeanInCorrected float64
}

func sum_mean(slice []int, docIds []string, correctedDocs map[string]bool) (int, float64, float64) {
    total := 0
    totalCorrected := 0
    for i, ed := range slice {
        total += ed
        if _, exists := correctedDocs[docIds[i]]; exists {
            totalCorrected += ed
        }
    }  
    return total, float64(total) / float64(len(slice)), float64(totalCorrected) / float64(len(correctedDocs))
}

func GetEvaluationStats(docs []common.Document, docMap map[string]int, correctedDocs map[string]bool) (EvalStats,  EvalStats, error) {
    docIds, originalDistances, correctedDistances, err := editDistances(docs, docMap, correctedDocs)
    if err != nil {
        return EvalStats{}, EvalStats{},  err
    }
    
    sum, mean, mean2:= sum_mean(originalDistances, docIds, correctedDocs)
    originalStats := EvalStats{
        Mean: mean,
        Total: sum,
        MeanInCorrected: mean2,
    }
    sum, mean, mean2 = sum_mean(correctedDistances, docIds, correctedDocs)
    correctedStats := EvalStats{
        Mean: mean,
        Total: sum,
        MeanInCorrected: mean2,
    }
    
    return originalStats, correctedStats, nil
}
