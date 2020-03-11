package fingerprinting

import (
	"postCorr/common"
	"postCorr/flags"
	"postCorr/readWrite"
	
	"fmt"
	"sort"
	
	minhash "github.com/rowanho/go-minhash"
	spooky "github.com/dgryski/go-spooky"
	metro "github.com/dgryski/go-metro"
	
	inverted "github.com/rowanho/Inverted-Index-Generator/invertedindex"
)

var total int = 0 
var totalSum float64 = 0.0
var scores []float64 = []float64{}
var score map[int]map[int]float64 = make(map[int]map[int]float64, 0)
var bools map[int]map[int]bool = make(map[int]map[int]bool, 0)

func ResetRuntime() {
	total = 0
	totalSum = 0.0
	scores = []float64{}
	score = make(map[int]map[int]float64, 0)
	bools = make(map[int]map[int]bool, 0)

}


func max_int(x int, y int) int {
	if x > y {
		return x
	} else {
		return y
	}
}

func min_int(x int, y int) int {
	if x < y {
		return x
	} else {
		return y
	}	
}

/**
* Using the inverted index, outputs the documents that have higher matches than the threshold docs
**/
func invertedIndexHighScores(fpList []map[uint64]int, targetDoc int, invertedIndex inverted.InvertedIndex) {
	numMatches := make([]int, len(fpList))
	
	for fp := range fpList[targetDoc] {
		contains := inverted.Find(invertedIndex, fp)
		if flags.JaccardType == common.WeightedJaccard {
			for _, c := range contains {
				// Add the minimum of the two counts
					numMatches[c] += min_int(fpList[c][fp], fpList[targetDoc][fp])
			} 
		} else {
			for _, c := range contains {
				numMatches[c] += 1
			}						
		}
	}
	
	score[targetDoc] = make(map[int]float64)
	bools[targetDoc] = make(map[int]bool)
	for i, n := range numMatches {
		if i == targetDoc {
			continue;
		}
		jaccard := 0.0
		if flags.JaccardType == common.WeightedJaccard {
			maxSum := 0
			for fp, freq := range fpList[i] {
				if _, exists := fpList[targetDoc][fp]; !exists {
					maxSum += freq
				} else {
					maxSum += max_int(fpList[targetDoc][fp], freq)					
				}
			}
			for fp, freq := range fpList[targetDoc] {
				if _, exists := fpList[i][fp]; !exists {
					maxSum +=  freq
				}
			}
			if maxSum > 0 {
				jaccard = float64(n) / float64(maxSum)
			}
			
		} else {
			// Jaccard Index
			l := len(fpList[i]) + len(fpList[targetDoc]) - n
			if l > 0 {
				jaccard = (float64(n) / float64(len(fpList[i]) + len(fpList[targetDoc]) - n))
			}
		}
		
		totalSum += jaccard
		total += 1
		
		scores = append(scores, jaccard)
		score[targetDoc][i] = jaccard
		bools[targetDoc][i] = true
	}
}

func mhash(b []byte) uint64 { return metro.Hash64(b, 0) } 
func getSimilarLsh(docs []common.Document) {
	ms := make([]*minhash.MinWise, len(docs))
	for i, doc := range(docs) {
		ms[i] = minhash.NewMinWise(spooky.Hash64, mhash, 100)
		for j := 0; j+flags.ShingleSize < len(doc.Text) + 1; j++ {
			ms[i].Push([]byte(string(doc.Text[j : j+flags.ShingleSize])))				
		}
	}
	
	for i := range(docs) {
		score[i] = make(map[int]float64)
		bools[i] = make(map[int]bool)
		for j := range(docs) {
			if i == j {
				continue;
			}
			s := ms[i].Similarity(ms[j])
			score[i][j] = s
			bools[i][j] = true
			total += 1
			totalSum += s
			scores = append(scores, s)
		}
	}		
}

func getSimilarModP(docs []common.Document) {
	fps := make([]map[uint64]int, len(docs))
	for i, doc := range docs {
		fp := ModP(preProcess(string(doc.Text)), flags.ShingleSize, flags.P)
		fps[i] = fp
	}
	invertedIndex := inverted.GenerateInvertedIndex(fps)
	for i := range fps {
		invertedIndexHighScores(fps, i, invertedIndex)
	}
}

func getSimilarWinnowing(docs []common.Document) {
	fps := make([]map[uint64]int, len(docs))
	for i, doc := range docs {
		fp := Winnowing(preProcess(string(doc.Text)), flags.ShingleSize, flags.WinnowingWindow)	
		fps[i] = fp
	}
	invertedIndex := inverted.GenerateInvertedIndex(fps)
	for i := range fps {
		invertedIndexHighScores(fps, i, invertedIndex)
	}
}


func GetAllPairwise(docs []common.Document) (map[int]map[int]float64, []float64) {
	fmt.Println(flags.JaccardType)
	if flags.FpType == common.MinhashFP {
		getSimilarLsh(docs)
	} else if flags.FpType == common.ModFP {
		getSimilarModP(docs)
	} else if flags.FpType == common.Winnowing {
		getSimilarWinnowing(docs)
	}	
	return score, scores
}
	
func GetSimilarDocuments(docs []common.Document) map[int]map[int]bool  {
	var documentAdjacencyList map[int]map[int]bool
	if flags.FpType == common.MinhashFP {
		getSimilarLsh(docs)
	} else if flags.FpType == common.ModFP {
		getSimilarModP(docs)
	} else if flags.FpType == common.Winnowing {
		getSimilarWinnowing(docs)
	}
	//fmt.Println(total)
	//fmt.Printf("Average jaccard index was %6.3f \n", totalSum / float64(total))
	
	pos := 0
	proportion := flags.SimilarityProportion
	sort.Float64s(scores)
	l := len(docs)
	numPairs := l*l - l
	threshold := 0.0
	numP := int(proportion * float64(numPairs))
	if  numP < len(scores) {
		threshold = scores[len(scores) - 1]
	
		threshold = scores[len(scores) - 1 - numP]
	
		for doc1 := range(score) {
			for doc2, s := range(score[doc1]) {
				if s < threshold {
					delete(bools[doc1], doc2)					
				} 
			}
		}
	}
	
	documentAdjacencyList = bools

	if flags.WriteData {
		readWrite.SerialiseJaccards(scores[pos:])
	}

	return documentAdjacencyList		
}
