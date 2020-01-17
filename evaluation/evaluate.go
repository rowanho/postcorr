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
func editDistances() ([]int, []int, error) {
    originalDistances := make([]int, 0)
    correctedDistances := make([]int, 0)
	err := filepath.Walk(flags.DirName,
		func(path string, info os.FileInfo, err error) error {
            right := path[len(flags.DirName):]
			if err != nil {
				return err
			}
			if info.IsDir() == false {
                correctable := true
                _, err := os.Stat(flags.OutDir  + right)
                if os.IsNotExist(err){
                    correctable = false
                }

                var groundTruth string
                var corrected string
                var original string
                var readErr error
                
				if flags.FormatType == common.Plaintext {
					original, readErr = readWrite.ReadString(path)
                    groundTruth, readErr = readWrite.ReadString(flags.Groundtruth  + right)
                    if correctable{
                        corrected, readErr = readWrite.ReadString(flags.OutDir  + right)
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
                    correctedDistances = append(correctedDistances, originalDistances)                    
                }
			}
			return nil
		},
	)
    
    if err != nil {
        return []int{}, []int{}, err
    }
	return originalDistances, correctedDistances, err
}

type EvalStats = struct {
    Mean float64
    Total int
}

func sum_mean(slice []int) (int, float64) {
    total := 0
    for _, i := range slice {
        total += i
    }  
    return total, float64(total) / float64(len(slice))
}

func GetEvaluationStats() (EvalStats,  EvalStats, error) {
    originalDistances, correctedDistances, err := editDistances()
    if err != nil {
        return EvalStats{}, EvalStats{},  err
    }
    
    sum, mean := sum_mean(originalDistances)
    originalStats := EvalStats{
        Mean: mean,
        Total: sum,
    }
    
    sum, mean = sum_mean(correctedDistances)
    correctedStats := EvalStats{
        Mean: mean,
        Total: sum,
    }
    
    return originalStats, correctedStats, nil
}
