package main

import (
	"postCorr/readWrite"
	"postCorr/common"
	"postCorr/alignment"
	"postCorr/queries"
	"postCorr/correction"
	"postCorr/flags"
	
	"fmt"
	"flag"
	"time"
	"context"
)

func main() {
	dirName  := flag.String("dir","test_dataset","path to dataset")
	formatType := flag.String("format", common.Plaintext, "the dataset file format")
	alignmentTolerance := flag.Int("tolerance", 10, "Tolerance for distances between alignments to identify as similar" )
	fpType := flag.String("fp", common.MinhashFP, "Fingeprinting method")
	jaccardThreshold := flag.Float64("jaccard", 0.05, "Jaccard index threshold for similarity")
	parallel := flag.Bool("parallel",false, "Whether or not to run alignments in parallel with goroutines")
	
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

	docIDList, docsErr := readWrite.TraverseAndIndexDocs()

	if docsErr != nil {
		fmt.Println("Error indexing documents %s", docsErr)
		return
	}
	
	fmt.Println("Waiting for documents to index")
	numDocs := len(docIDList)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	for {
	    count, err := queries.CountDocs(common.DocumentIndex) 
	    if count == numDocs  || ctx.Err() != nil || err != nil {
	        break
	    }
	}
		
	likelyMatchingDocs := getSimilarDocuments(docIDList)
	
	fmt.Println(likelyMatchingDocs)
	fmt.Println("Aligning")
	var alignmentCount int
	
	if flags.Parallel {
		alignmentCount = alignAndIndexParallel(likelyMatchingDocs)
	} else {
		alignmentCount = alignAndIndexSerial(likelyMatchingDocs) 
	}
	

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	for {
	    count, err := queries.CountDocs(common.AlignmentIndex) 
	    if count == alignmentCount || ctx.Err() != nil || err != nil {
	        break
	    }
	}
	
	alignmentAdjacencyList := getSimilarAlignments(docIDList)
	fmt.Println(alignmentAdjacencyList)
	totalCorrections += correction.ClusterAndCorrectAlignments(alignmentAdjacencyList, 1)
	fmt.Println("Number of corrections made: ", totalCorrections)
	queries.DeleteIndexes([]string{common.AlignmentIndex, common.FpIndex, common.DocumentIndex, common.MinHashIndex})
}

func getSimilarDocuments(docIDList []string) map[string]map[string]bool{
	similarityScore := 0
	likelyMatchingDocs := make(map[string]map[string]bool, 0)
	for _, docID := range docIDList {
		if (flags.FpType == common.ModFP) {
			similarDocIDs, _ := queries.GetSimilarFps(common.FpIndex, docID, docIDList, flags.JaccardThreshold)
			likelyMatchingDocs[docID] = similarDocIDs	
			similarityScore += len(similarDocIDs)		
		} else if (flags.FpType == common.MinhashFP) {
			similarDocIDs, _ := queries.GetSimilarMinHashes(common.MinHashIndex, docID, docIDList)
			likelyMatchingDocs[docID] = similarDocIDs		
			similarityScore += len(similarDocIDs)	
		}
	}
	fmt.Println("Found ", similarityScore/2, " high similarity pairs")
	return likelyMatchingDocs
}

func alignAndIndexParallel(likelyMatchingDocs map[string]map[string]bool) int {
	alignmentCount := 0
	for primID, secIDs := range likelyMatchingDocs {
			primDoc, _ := queries.GetDocByID(common.DocumentIndex, primID)
			for secID, _:= range secIDs {
				if _, exists := likelyMatchingDocs[secID][primID]; exists {
						delete(likelyMatchingDocs[secID],primID)
				}
			}
			
			alignmentChannel := make(chan []common.Alignment)
			
			for secID, _:= range secIDs {
				secDoc, _ := queries.GetDocByID(common.DocumentIndex, secID)
				  go func(channel chan []common.Alignment, primDoc common.Document, secDoc common.Document, secID string){
					alignments, inverseAlignments  := alignment.GetAlignments(1.0, 2.0, primDoc, secDoc, 1, 0.0)
					channel <- alignments
					channel <- inverseAlignments
				}(alignmentChannel, primDoc, secDoc, secID)
				
			}
			
			for i := 0; i < len(secIDs) * 2; i++ {
				alignments := <- alignmentChannel
				for _, al := range alignments {
					queries.IndexAlignment(common.AlignmentIndex, al)
				}
				alignmentCount += len(alignments)
			}
			
	}
	return alignmentCount
}

func alignAndIndexSerial(likelyMatchingDocs map[string]map[string]bool) int {
	alignmentCount := 0
	for primID, secIDs := range likelyMatchingDocs {
			primDoc, _ := queries.GetDocByID(common.DocumentIndex, primID)
			for secID, _:= range secIDs {
				if _, exists := likelyMatchingDocs[secID][primID]; exists {
						delete(likelyMatchingDocs[secID],primID)
				}
				secDoc, _ := queries.GetDocByID(common.DocumentIndex, secID)
				alignments, inverseAlignments := alignment.GetAlignments(1.0, 2.0, primDoc, secDoc, 1, 0.0)
				for i, al := range alignments {
					queries.IndexAlignment(common.AlignmentIndex, al)
					queries.IndexAlignment(common.AlignmentIndex, inverseAlignments[i])				
				}
				alignmentCount += 1
			}
	}
	return alignmentCount
	
}

func getSimilarAlignments(docIDList []string) map[string][]string {
	alignmentAdjacencyList := make(map[string][]string, 0)
	// Loop through all alignments
	for _, docID := range docIDList {
		fmt.Println(docID)
		alignments,_ := queries.GetAlignmentsByPrimID(common.AlignmentIndex, docID)
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


