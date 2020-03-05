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

func main() {
	dirName := flag.String("input", "test_dataset", "path to dataset")
	groundTruth := flag.String("groundtruth", "", "Directory containing groundtruth data")
	writeOutput := flag.Bool("write", true, "Whether or not to write output to file")
	writeData := flag.Bool("writeData", false, "Whether to write data to file for eg: distribution plot.")
	detailedEvaluation := flag.Bool("detailedEval", false, "Whether to run detailed evaluation of edit distance.")
	fpType := flag.String("fp", common.MinhashFP, "Fingeprinting method")
	similarityProportion := flag.Float64("proportion", 0.05, "The proportion of document pairs to align")
	jaccardType := flag.String("jaccard", common.WeightedJaccard, "The type of jaccard similarity, 'regular' or 'weighted'")
	shingleSize := flag.Int("shingleSize", 7, "Length of shingle")
	runAlignment := flag.Bool("align", true, "Whether or not to run the alignment/correction phases")
	winnowingWindow := flag.Int("t", 15, "Size of winnowing window t")
	fastAlign := flag.Bool("fastAlign", false, "Whether or not to use heuristic alignment (faster but less accurate)")
	p := flag.Int("p", 5, "P to mod by when using modp")
	numAligns := flag.Int("numAligns", 2, "The number of disjoint alignments we attempt to make")
	flag.Parse()
	
	flags.WriteOutput = *writeOutput
	flags.DirName = *dirName
	flags.WriteData = * writeData
	flags.FpType = *fpType
	flags.DetailedEvaluation = *detailedEvaluation
	flags.ShingleSize = * shingleSize
	flags.SimilarityProportion = *similarityProportion
	flags.JaccardType = * jaccardType
	flags.RunAlignment = * runAlignment
	flags.WinnowingWindow = *winnowingWindow
	flags.P = *p
	flags.Groundtruth = *groundTruth
	flags.FastAlign = *fastAlign
	flags.NumAligns = *numAligns
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
		fmt.Println("Error reading documents %s", docsErr)
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
		alignments, alignmentsPerDocument := alignment.AlignParallel(documentAdjacencyList, docList)
		
		if flags.WriteData {
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
		originalStats, correctedStats, _ := evaluation.GetEvaluationStats(docList, docMap, correctedDocs)
		
		fmt.Printf("Total edit distance before correction: %d\n", originalStats.Total)
		fmt.Printf("Total edit distance after correction: %d \n", correctedStats.Total)
		
		fmt.Printf("Mean edit distance before correction: %5.2f \n", originalStats.Mean)
		fmt.Printf("Mean edit distance after correction: %5.2f \n", correctedStats.Mean)
		
		if len(correctedDocs) > 0 {
			fmt.Printf("Out of %d the corrected documents, mean edit distance changed from %5.2f to %5.2f \n", 
							len(correctedDocs), originalStats.MeanInCorrected, correctedStats.MeanInCorrected)
		} else {
			fmt.Println("No documents corrected!")
		}
	}
}

