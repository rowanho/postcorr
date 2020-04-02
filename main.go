package main

import (
	"postCorr/alignment"
	"postCorr/common"
	"postCorr/correction"
	"postCorr/flags"
	"postCorr/fingerprinting"
	"postCorr/readWrite"
	"postCorr/evaluation"
	"flag"
	"fmt"
)
func main2(){
	EvaluateK()
	//EvaluateJaccard()
}
func main() {
	dirName := flag.String("input", "test_dataset", "path to dataset")
	groundTruth := flag.String("groundtruth", "", "Directory containing groundtruth data")
	writeOutput := flag.Bool("write", true, "Whether or not to write output to file")
	logLevel := flag.Int("logLevel", 0, "The level of logging used, 0 (no logging), 1 or 2")
	detailedEvaluation := flag.Bool("detailedEval", false, "Whether to run detailed evaluation of edit distance.")
	fpType := flag.String("fp", common.MinhashFP, "Fingeprinting method")
	similarityProportion := flag.Float64("proportion", 0.05, "The proportion of document pairs to align")
	jaccardType := flag.String("jaccard", common.WeightedJaccard, "The type of jaccard similarity, 'regular' or 'weighted'")
	shingleSize := flag.Int("shingleSize", 7, "Length of shingle")
	runAlignment := flag.Bool("align", true, "Whether or not to run the alignment/correction phases")
	winnowingT := flag.Int("t", 15, "Size of winnowing threshold t, must be >= k")
	affine := flag.Bool("affine", false, "Whether or not to use affine gap scoring")
	fastAlign := flag.Bool("fastAlign", false, "Whether or not to use heuristic alignment (faster but less accurate)")
	p := flag.Int("p", 5, "P to mod by when using modp")
	numAligns := flag.Int("numAligns", 2, "The number of disjoint alignments we attempt to make")
	useLM := flag.Bool("useLM", false, "Whether to use a language model to inform correction")
	lmThreshold := flag.Float64("lmThreshold", 0.1, "The probability score under which language model permits correction")
	flag.Parse()
	
	flags.WriteOutput = *writeOutput
	flags.DirName = *dirName
	flags.LogLevel = *logLevel
	flags.FpType = *fpType
	flags.DetailedEvaluation = *detailedEvaluation
	flags.ShingleSize = * shingleSize
	flags.SimilarityProportion = *similarityProportion
	flags.JaccardType = * jaccardType
	flags.RunAlignment = * runAlignment
	flags.WinnowingT = *winnowingT
	flags.P = *p
	flags.Groundtruth = *groundTruth
	flags.FastAlign = *fastAlign
	flags.NumAligns = *numAligns
	flags.Affine = *affine
	flags.UseLM = *useLM
	flags.LmThreshold = *lmThreshold
	execute()
	
}

/**
* Executes the main program pipeline
**/
func execute() {
	var totalCorrections int
 	var correctedDocs map[string] bool
	docList, docsErr := readWrite.TraverseDocs()

	if docsErr != nil {
		fmt.Printf("Error reading documents %s", docsErr)
		return
	}

	docMap := make(map[string]int)
	for i, doc := range docList {
		docMap[doc.ID] = i
	}
	
	documentAdjacencyList := fingerprinting.GetSimilarDocuments(docList)
	
	numPairs := 0
	
	for _, similarDocs := range documentAdjacencyList {
		numPairs += len(similarDocs)
	}
	fmt.Printf("Found %d high scoring pairs \n", numPairs / 2)
	
	if flags.RunAlignment {
		fmt.Println("Aligning")
		var alignments map[string]common.Alignment
		var alignmentsPerDocument map[string][]string
		fmt.Println(flags.Affine)
		if flags.Affine && !flags.FastAlign {
			// We might run  out of memory if we create a lot of go routines, best to use serial methods
			alignments, alignmentsPerDocument = alignment.AlignSerial(documentAdjacencyList, docList)
		} else {
			alignments, alignmentsPerDocument = alignment.AlignParallel(documentAdjacencyList, docList)
		}
		
		if flags.LogLevel > 1 {
			readWrite.SerialiseGraph(alignments, alignmentsPerDocument)
		}
		scoreSum := 0
		for _, al := range alignments {
			scoreSum += al.Score 
		}

		fmt.Printf("Score sum: %d \n", scoreSum)
		alignmentAdjacencyList := alignment.GetSimilarAlignments(alignments, alignmentsPerDocument)
		correctedDocs, totalCorrections = correction.ClusterAndCorrectAlignments(alignmentAdjacencyList, alignments, docList, docMap)
		fmt.Println("Number of corrections made: ", totalCorrections)
	}
	
	// Evaluation
	if flags.RunAlignment &&  len(flags.Groundtruth) > 0 {
		originalStats, correctedStats,  originalWordStats, correctedWordStats, _ := evaluation.GetEvaluationStats(docList, docMap, correctedDocs)
		
		fmt.Printf("Total character distance before correction: %d\n", originalStats.Total)
		fmt.Printf("Total character distance after correction: %d \n", correctedStats.Total)
		
		fmt.Printf("Mean character distance before correction: %5.2f \n", originalStats.Mean)
		fmt.Printf("Mean character distance after correction: %5.2f \n", correctedStats.Mean)
		
		fmt.Printf("Total word distance before correction: %d\n", originalWordStats.Total)
		fmt.Printf("Total word distance after correction: %d \n", correctedWordStats.Total)
		
		fmt.Printf("Mean word distance before correction: %5.2f \n", originalWordStats.Mean)
		fmt.Printf("Mean word distance after correction: %5.2f \n", correctedWordStats.Mean)
	
		if len(correctedDocs) > 0 {
			fmt.Printf("Out of %d the corrected documents, mean edit distance changed from %5.2f to %5.2f \n", 
							len(correctedDocs), originalStats.MeanInCorrected, correctedStats.MeanInCorrected)
		} else {
			fmt.Println("No documents corrected!")
		}
	}
}

