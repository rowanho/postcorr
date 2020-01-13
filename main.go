package main

import (
	"postCorr/alignment"
	"postCorr/common"
	"postCorr/correction"
	"postCorr/flags"
	"postCorr/queries"
	"postCorr/readWrite"

	"context"
	"flag"
	"fmt"
	"time"
)

func main() {
	dirName := flag.String("dir", "test_dataset", "path to dataset")
	formatType := flag.String("format", common.Plaintext, "the dataset file format")
	alignmentTolerance := flag.Int("tolerance", 10, "Tolerance for distances between alignments to identify as similar")
	fpType := flag.String("fp", common.MinhashFP, "Fingeprinting method")
	jaccardThreshold := flag.Float64("jaccard", 0.05, "Jaccard index threshold for similarity")
	parallel := flag.Bool("parallel", false, "Whether or not to run alignments in parallel with goroutines")

	flag.Parse()

	flags.DirName = *dirName
	flags.FormatType = *formatType
	flags.AlignmentTolerance = *alignmentTolerance
	flags.FpType = *fpType
	flags.JaccardThreshold = *jaccardThreshold
	flags.Parallel = *parallel
	execute()
}

/**
* Executes the main program pipeline
**/
func execute() {
	totalCorrections := 0
	queries.CreateAlignmentIndex(common.AlignmentIndex)
	queries.CreateFingerprintIndex(common.FpIndex)
	//queries.CreateLSHFingerprintIndex(common.FpLSHIndex, 5, 7, 512)
	time.Sleep(1 * time.Second)

	docList, docsErr := readWrite.TraverseDocs()

	if docsErr != nil {
		fmt.Println("Error reading documents %s", docsErr)
		return
	}

	docMap := make(map[string]int)
	for i, doc := range docList {
		docMap[doc.ID] = i
	}
	likelyMatchingDocs := fingerprinting.GetSimilarDocuments(docsList, docsList, docMap)

	fmt.Println(likelyMatchingDocs)
	fmt.Println("Aligning")
	var alignments []common.Alignment

	if flags.Parallel {
		alignments = alignment.AlignParallel(likelyMatchingDocs, docList)
	} else {
		alignments = alignment.AlignSerial(likelyMatchingDocs, docList)
	}

	alignmentAdjacencyList := alignment.GetSimilarAlignments(alignments)
	fmt.Println(alignmentAdjacencyList)
	totalCorrections += correction.ClusterAndCorrectAlignments(alignmentAdjacencyList, 1)
	fmt.Println("Number of corrections made: ", totalCorrections)
}

func getSimilarAlignments(docIDList []string) map[string][]string {
	alignmentAdjacencyList := make(map[string][]string, 0)
	// Loop through all alignments
	for _, docID := range docIDList {
		fmt.Println(docID)
		alignments, _ := queries.GetAlignmentsByPrimID(common.AlignmentIndex, docID)
		fmt.Println(len(alignments))
		for _, al := range alignments {
			matchingAlignmentIds, _ := queries.GetMatchingAlignments(common.AlignmentIndex,
				al,
				flags.AlignmentTolerance)
			alignmentAdjacencyList[al.ID] = matchingAlignmentIds
		}
	}
	return alignmentAdjacencyList
}
