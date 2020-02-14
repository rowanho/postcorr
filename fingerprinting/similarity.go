package fingerprinting

import (
	"postCorr/common"
	"postCorr/flags"
	
	"fmt"
	"os"
	
	inverted "github.com/rowanho/Inverted-Index-Generator/invertedindex"
)
var total int = 0 
var totalSum float64 = 0.0
var scores []float64 = []float64{}

/**
* Using the inverted index, outputs the documents that have higher matches than the threshold docs
**/
func invertedIndexHighScores(fpList []map[uint64]bool, targetDoc int, invertedIndex inverted.InvertedIndex, threshold float64) map[int]bool {
	numMatches := make([]int, len(fpList))
	
	for fp := range fpList[targetDoc] {
		contains := inverted.Find(invertedIndex, fp)
		for _, c := range contains {
			numMatches[c] += 1
		}
	}
	
	highScoring := make(map[int]bool)
	for i, n := range numMatches {
		if i == targetDoc {
			continue;
		}
		// Jaccard Index
		l := len(fpList[i]) + len(fpList[targetDoc]) - n
		jaccard := 0.0
		if l > 0 {
			jaccard = (float64(n) / float64(len(fpList[i]) + len(fpList[targetDoc]) - n))
		}
		
		if  jaccard >= threshold {
			highScoring[i] = true
		}  
		totalSum += jaccard
		total += 1
		scores = append(scores, jaccard)
	}
	
	return highScoring
}

func writeScores() {
	f, _ := os.Create(fmt.Sprintf("%s_jaccard_indexes%d.txt", flags.DirName, flags.ShingleSize))
	defer f.Close()
	
	for _, j := range scores {
		f.WriteString(fmt.Sprintf("%f", j) + "\n")
	}

}

func getSimilarLsh(docs []common.Document) map[int]map[int]bool {
	GetLSHObject(100, flags.JaccardThreshold, len(docs))
	fps := make([]common.LSH_fp, len(docs))
	for i, doc := range docs {
		fp := MinHash(i, preProcess(string(doc.Text)), 7)
		fps[i] = fp
	}
	
	IndexMinHashObject()

	documentAdjacencyList := make(map[int]map[int]bool)
	for i, fp := range fps {
		documentAdjacencyList[i] = make(map[int]bool)
		sameBucketIds := SameBucketIds(fp.Signature)
		for _, id := range sameBucketIds {
			if id != i {
				documentAdjacencyList[i][id] = true
			}
		}
	}
	return documentAdjacencyList
}

func getSimilarModP(docs []common.Document) map[int]map[int]bool {
	fps := make([]map[uint64]bool, len(docs))
	for i, doc := range docs {
		fp := ModP(preProcess(string(doc.Text)), flags.ShingleSize, flags.P)
		fps[i] = fp
	}
	invertedIndex := inverted.GenerateInvertedIndex(fps)
	documentAdjacencyList := make(map[int]map[int]bool)
	for i := range fps {
		documentAdjacencyList[i] = invertedIndexHighScores(fps, i, invertedIndex, flags.JaccardThreshold)
	}
	return documentAdjacencyList
}

func getSimilarWinnowing(docs []common.Document) map[int]map[int]bool {
	fps := make([]map[uint64]bool, len(docs))
	for i, doc := range docs {
		fp := Winnowing(preProcess(string(doc.Text)), flags.ShingleSize, flags.WinnowingWindow)	
		fps[i] = fp
	}
	invertedIndex := inverted.GenerateInvertedIndex(fps)
	documentAdjacencyList := make(map[int]map[int]bool)
	for i := range fps {
		documentAdjacencyList[i] = invertedIndexHighScores(fps, i, invertedIndex, flags.JaccardThreshold)
	}
	return documentAdjacencyList
}


func GetSimilarDocuments(docs []common.Document) map[int]map[int]bool {
	var documentAdjacencyList map[int]map[int]bool
	if flags.FpType == common.MinhashFP {
		documentAdjacencyList = getSimilarLsh(docs)
	} else if flags.FpType == common.ModFP {
		documentAdjacencyList = getSimilarModP(docs)
	} else if flags.FpType == common.Winnowing {
		documentAdjacencyList = getSimilarWinnowing(docs)
	}
	fmt.Println(total)
	fmt.Printf("Average jaccard index was %6.3f ", totalSum / float64(total))
	
	if flags.WriteData {
		writeScores()
	}
	
	return documentAdjacencyList
}
