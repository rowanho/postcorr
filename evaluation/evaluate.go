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
func editDistances(docs []common.Document, docMap map[string]int, correctedDocs map[string]bool) ([]string, []int, []int, error) {
    originalDistances := make([]int, 0)
    correctedDistances := make([]int, 0)
    docIds := make([]string, 0)
    improvedStats := levenshtein.NewEditStats()
    worsenedStats := levenshtein.NewEditStats()
    sameStats := levenshtein.NewEditStats()
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
                    originalDist, ogEditStats :=  levenshtein.ComputeDistanceWithConstruction(original, groundTruth)
                    originalDistances = append(originalDistances, originalDist)
                    if correctable {
                        correctedDist, correctedEditStats := levenshtein.ComputeDistanceWithConstruction(corrected, groundTruth)
                        correctedDistances = append(correctedDistances, correctedDist)
                        mergeStats(improvedStats, worsenedStats, sameStats, ogEditStats, correctedEditStats)
                    } else {
                        correctedDistances = append(correctedDistances, originalDist)                    
                    }                    
                } else {
                    originalDist := levenshtein.ComputeDistance(original, groundTruth)
                    originalDistances = append(originalDistances, originalDist)
                    
                    if correctable {
                        correctedDist := levenshtein.ComputeDistance(corrected, groundTruth)
                        correctedDistances = append(correctedDistances, correctedDist)
                    } else {
                        correctedDistances = append(correctedDistances, originalDist)                    
                    }                    
                }
                docIds = append(docIds, right)
			}
			return nil
		},
	)
    
    if err != nil {
        return []string{}, []int{}, []int{}, err
    }
    if flags.DetailedEvaluation {
        statsBreakDown(improvedStats, worsenedStats, sameStats)
    }
	return  docIds, originalDistances, correctedDistances, err
}


func statsBreakDown(improvedStats, worsenedStats, sameStats levenshtein.EditStats) {
    printStats("corrected", improvedStats)
    printStats("unchanged", sameStats)
    printStats("degraded", worsenedStats)
}

func printStats(printString string, stats levenshtein.EditStats) {
    topDels := max_n_dict(stats.Dels, 20)
    topIns := max_n_dict(stats.Ins, 20)
    topSubs := max_n_dict(stats.Subs, 20)
    fmt.Println("Top deleted", printString, "characters", topDels)
    fmt.Println("Top inserted", printString, "characters", topIns)
    fmt.Println("Top substituted", printString, "characters", topSubs)
}

func mergeStats(improvedStats, worsenedStats, sameStats, ogEditStats, correctedEditStats levenshtein.EditStats) {
    mergeDict(improvedStats.Dels, worsenedStats.Dels, sameStats.Dels, ogEditStats.Dels, correctedEditStats.Dels)
    mergeDict(improvedStats.Ins, worsenedStats.Ins, sameStats.Ins, ogEditStats.Ins, correctedEditStats.Ins)
    mergeDict(improvedStats.Subs, worsenedStats.Subs, sameStats.Subs, ogEditStats.Subs, correctedEditStats.Subs)
}

func mergeDict(improved, worsened, same, og, corrected map[string]int) {
    for s := range og {
        if _, exists := corrected[s]; exists {
            if og[s] == corrected[s] {
                same[s] += og[s]
            } else if og[s] <  corrected[s] {
                worsened[s] += og[s] - corrected[s]
            } else {
                improved[s] += corrected[s] - og[s]
            }
        } else {
            improved[s] += og[s]
        }
    }
    
    for s := range corrected {
        if _, exists := og[s]; !exists {
            worsened[s] += corrected[s]
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
