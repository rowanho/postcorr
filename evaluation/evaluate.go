package evaluation

import (
    "postCorr/readWrite"
    "postCorr/flags"
    "postCorr/common"
    
    "os"
    "path"
    "path/filepath"
    "fmt"
    "sort"
    
    "github.com/rowanho/levenshtein"
)

/**
* Traverses the input and output directories, and sums the edit distances of each file
* Between the data and the ground truth data
**/
func editDistances(docs []common.Document, docMap map[string]int, correctedDocs map[string]bool) ([]string, []int, []int, []int, []int, error) {
    originalDistances := make([]int, 0)
    correctedDistances := make([]int, 0)
    ogWordDistances := make([]int, 0)
    corrWordDistances := make([]int, 0)
    
    docIds := make([]string, 0)
    changedStatsTotal := levenshtein.NewEditStats()
    changedStats := levenshtein.NewEditStats()
    var changeDist int
	err := filepath.Walk(flags.DirName,
		func(pth string, info os.FileInfo, err error) error {
            right := ""
            if pth != flags.DirName {
                right = pth[len(flags.DirName) + 1:]
            }
			if err != nil {
				return err
			}
			if info.IsDir() == false {
                _, correctable := correctedDocs[right]

                var groundTruth []rune
                var corrected []rune
                
				original, readErr := readWrite.ReadRunes(pth)
                groundTruth, readErr = readWrite.ReadRunes(path.Join(flags.Groundtruth, right))
                if correctable {
                    corrected = docs[docMap[right]].Text
                }
				
				if readErr != nil {
                    fmt.Println(readErr)
					return readErr
				}
                        
                if flags.DetailedEvaluation {
                    if correctable {
                        changeDist, changedStats =  levenshtein.ComputeDistanceWithConstruction(original, corrected)
                    }
                } 
                mergeStats(changedStats, changedStatsTotal)
                originalDist := levenshtein.ComputeDistance(original, groundTruth)
                originalDistances = append(originalDistances, originalDist)
                ogWordDist := levenshtein.ComputeWordDistance(original, groundTruth)
                ogWordDistances = append(ogWordDistances, ogWordDist)
                
                if correctable {
                    correctedDist := levenshtein.ComputeDistance(corrected, groundTruth)
                    correctedDistances = append(correctedDistances, correctedDist)
                    correctedWordDist := levenshtein.ComputeWordDistance(corrected, groundTruth)
                    corrWordDistances = append(corrWordDistances, correctedWordDist)
                } else {
                    correctedDistances = append(correctedDistances, originalDist) 
                    corrWordDistances = append(corrWordDistances, ogWordDist)                   
                }                                
                docIds = append(docIds, right)
			}
			return nil
		},
	)
    
    if err != nil {
        return []string{}, []int{}, []int{}, []int{}, []int{}, err
    }
    if flags.DetailedEvaluation {
        printStats(changedStatsTotal)
    }
	return  docIds, originalDistances, correctedDistances, ogWordDistances, corrWordDistances, err
}


func printStats(stats levenshtein.EditStats) {
    topDels := max_n_dict(stats.Dels, 20)
    topIns := max_n_dict(stats.Ins, 20)
    topSubs := max_n_dict(stats.Subs, 20)
    fmt.Println("Top deleted characters", topDels)
    fmt.Println("Top inserted characters", topIns)
    fmt.Println("Top substituted characters", topSubs)
}

func mergeStats(updateStats, totalStats levenshtein.EditStats) {
    mergeDict(updateStats.Dels, totalStats.Dels)
    mergeDict(updateStats.Ins, totalStats.Ins)
    mergeDict(updateStats.Subs, totalStats.Subs)
}

func mergeDict(update, total map[string]int) {
    for s := range update {
        if _, exists := total[s]; exists {
            total[s] += update[s]
        } else {
            total[s] = update[s]
        }
    }    
}

type kv struct {
    Key   string
    Value int
}

func max_n_dict(counts map[string]int, n int) map[string]int {
    
    if n > len(counts) {
        n = len(counts)
    }
    ss := make([]kv, len(counts))
    i := 0
    for k, v := range counts {
        ss[i] =  kv{k, v} 
        i += 1
    }

    sort.Slice(ss, func(i, j int) bool {
        return ss[i].Value > ss[j].Value
    })
    
    res := make(map[string]int)
    for _, kv := range ss[:n] {
        res[kv.Key] = kv.Value
    }
    return res
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

type EvalStats = struct {
    Mean float64
    Total int
    MeanInCorrected float64
}

func GetEvaluationStats(docs []common.Document, docMap map[string]int, correctedDocs map[string]bool) (EvalStats,  EvalStats, EvalStats, EvalStats, error) {
    docIds, originalDistances, correctedDistances, 
    ogWordDistances, corrWordDistances, err := editDistances(docs, docMap, correctedDocs)
    if err != nil {
        return EvalStats{}, EvalStats{}, EvalStats{}, EvalStats{}, err
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
    
    sum, mean, mean2 = sum_mean(ogWordDistances, docIds, correctedDocs)
    originalWordStats := EvalStats{
        Mean: mean,
        Total: sum,
        MeanInCorrected: mean2,
    }
    sum, mean, mean2 = sum_mean(corrWordDistances, docIds, correctedDocs)
    correctedWordStats := EvalStats{
        Mean: mean,
        Total: sum,
        MeanInCorrected: mean2,
    }
    
    return originalStats, correctedStats, originalWordStats, correctedWordStats, nil
}
