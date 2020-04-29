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
	dirName := flag.String("input", "", "Path to directory containing OCR dataset")
	groundTruth := flag.String("groundtruth", "", "Path to directory containing groundtruth dataset")
	writeOutput := flag.Bool("write", true, "Whether or not to write output to file in the folder 'corrected'")
	logging := flag.Bool("logging", true, "Whether to generate log files in the folder 'logs'")
	fpType := flag.String("fp", common.ModFP, "Fingerprinting method: 'minhash', 'modp' or 'winnowing'")
	similarityProportion := flag.Float64("candidate_proportion", 0.05, "The proportion of document pairs to align")
	jaccardType := flag.String("jaccard", common.WeightedJaccard, "The type of jaccard similarity, 'regular' or 'weighted'")
	shingleSize := flag.Int("k", 7, "Length of k-grams used for shingling")
	winnowingT := flag.Int("t", 15, "Size of winnowing threshold t, must be >= k")
	affine := flag.Bool("affine", false, "Whether or not to use affine gap scoring")
	fastAlign := flag.Bool("fast_align", false, "Whether or not to use heuristic alignment (faster but less accurate)")
	bandWidth := flag.Int("band_width", 200, "The dynamic programming band width w for the heuristic algorithm.")
	p := flag.Int("p", 3, "P to mod by when using modp")
	numAligns := flag.Int("num_aligns", 2, "The number of disjoint alignments we attempt to make")
	alignThreshold := flag.Int("align_threshold", 0, "The minimum previous alignment score with which to keep aligning 2 documents.")
	useLM := flag.Bool("use_lm", false, "Whether to use a language model to inform correction")
	lmThreshold := flag.Float64("lm_threshold", 0.1, "The probability score under which language model permits correction")
	handleInsertionDeletion := flag.Bool("insert_delete", true, "The correction algorithm tries to handle insertion and deletion errors.")
	lInsert := flag.Int("l_insert", 2, "The maximum length of character sequence that the algorithm will attempt considers an erroneous insertion in consensus vote.")
	lDelete := flag.Int("l_delete", 2, "The maximum length of character sequence that the algorithm will attempt considers an erroneous deletion in consensus vote. ")
	flag.Parse()

	flags.WriteOutput = *writeOutput
	flags.DirName = *dirName
	flags.Logging = *logging
	flags.FpType = *fpType
	flags.K = *shingleSize
	flags.SimilarityProportion = *similarityProportion
	flags.JaccardType = *jaccardType
	flags.WinnowingT = *winnowingT
	flags.P = *p
	flags.AlignThreshold = *alignThreshold
	flags.Groundtruth = *groundTruth
	flags.FastAlign = *fastAlign
	flags.BandWidth = * bandWidth
	flags.NumAligns = *numAligns
	flags.Affine = *affine
	flags.UseLM = *useLM
	flags.LmThreshold = *lmThreshold
	flags.HandleInsertionDeletion = *handleInsertionDeletion
	flags.LInsert= *lInsert
	flags.LDelete = *lDelete
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

	fmt.Println("Running candidate selection...")
	documentAdjacencyList := fingerprinting.GetSimilarDocuments(docList)

	numPairs := 0

	for _, similarDocs := range documentAdjacencyList {
		numPairs += len(similarDocs)
	}
	fmt.Println("")
	fmt.Printf("Found %d high scoring pairs \n", numPairs / 2)

	fmt.Println("Running Alignment...")
	var alignments map[string]common.Alignment
	var alignmentsPerDocument map[string][]string
	fmt.Println(flags.Affine)
	if flags.Affine && !flags.FastAlign {
		// We might run  out of memory if we create a lot of go routines, best to use serial methods
		alignments, alignmentsPerDocument = alignment.AlignSerial(documentAdjacencyList, docList)
	} else {
		alignments, alignmentsPerDocument = alignment.AlignParallel(documentAdjacencyList, docList)
	}

	if flags.Logging {
		readWrite.SerialiseGraph(alignments, alignmentsPerDocument)
	}
	fmt.Println("")
	fmt.Println("Running Correction")

	alignmentAdjacencyList := alignment.GetSimilarAlignments(alignments, alignmentsPerDocument)
	correctedDocs, totalCorrections = correction.ClusterAndCorrectAlignments(alignmentAdjacencyList, alignments, docList, docMap)
	fmt.Println("Number of corrections made: ", totalCorrections)


	// Evaluation
	if len(flags.Groundtruth) > 0 {
		fmt.Println("Running Evaluation...")
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
