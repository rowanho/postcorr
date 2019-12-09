package main

import (
	"postCorr/readWrite"
	"postCorr/common"
	"postCorr/alignment"
	"postCorr/queries"
	"postCorr/correction"
	
	"fmt"
	"flag"
	"time"
)

func main() {
	dirName := flag.String("dir","test_dataset","path to dataset")
	formatType := flag.String("format", common.Plaintext, "the dataset file format")
	alignmentTolerance := flag.Int("tolerance", 10, "Tolerance for distances between alignments to identify as similar" )
	flag.Parse()
	
	execute(*dirName, *formatType, *alignmentTolerance)
}


/**
* Executes the main program pipeline
**/
func execute(dirName string, formatType string, alignmentTolerance int) {
	
	queries.CreateAlignmentIndex(common.AlignmentIndex)
	queries.CreateLSHFingerprintIndex(common.FpLSHIndex, 5, 7, 512)
	time.Sleep(1 * time.Second)

	docIDList, docsErr := readWrite.TraverseAndIndexDocs(dirName, formatType)

	if docsErr != nil {
		fmt.Println("Error indexing documents %s", docsErr)
		return
	}
	time.Sleep(1 * time.Second)
	
	likelyMatchingDocs := getSimilarDocuments(docIDList)
	
	fmt.Println(likelyMatchingDocs)
	alignAndIndex(likelyMatchingDocs)
	time.Sleep(1 * time.Second)

	alignmentAdjacencyList := getSimilarAlignments(docIDList, alignmentTolerance)
	fmt.Println(alignmentAdjacencyList)
	correction.ClusterAndCorrectAlignments(alignmentAdjacencyList, 2)
}

func getSimilarDocuments(docIDList []string) map[string][]string{
	likelyMatchingDocs := make(map[string][]string, 0)
	
	for _, docID := range docIDList {
		similarDocIDs, _ := queries.GetSimilarFpsLSH(common.FpLSHIndex, docID)
		likelyMatchingDocs[docID] = similarDocIDs
	}
	return likelyMatchingDocs	
}

func alignAndIndex(likelyMatchingDocs map[string][]string) {
	for primID, secIDs := range likelyMatchingDocs {
		primDoc, _ := queries.GetDocByID(common.DocumentIndex, primID)
		for _, secID := range secIDs {
			secDoc, _ := queries.GetDocByID(common.DocumentIndex, secID)
			alignments := alignment.GetAlignments(1.0, 1.0, primDoc, secDoc, 3)
			for _, al := range alignments {
				queries.IndexAlignment(common.AlignmentIndex, al)
			}
		}
	}	
}

func getSimilarAlignments(docIDList []string, tolerance int) map[string][]string {
	
	alignmentAdjacencyList := make(map[string][]string, 0)
	// Loop through all alignments
	for _, docID := range docIDList {
		fmt.Println(docID)
		alignments,_ := queries.GetAlignmentsByPrimID(common.AlignmentIndex, docID)
		fmt.Println(len(alignments))
		for _, al := range alignments {
			matchingAlignmentIds, _ := queries.GetMatchingAlignments(common.AlignmentIndex, 
																al, 
																tolerance)
			alignmentAdjacencyList[al.ID] = matchingAlignmentIds
		}
	}
	return alignmentAdjacencyList
}


