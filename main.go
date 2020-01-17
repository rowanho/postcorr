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
	outDir := flag.String("output", "corrected_data", "Folder to write the output data to")
	groundTruth := flag.String("groundtruth", "", "Directory containing groundtruth data")
	writeOutput := flag.Bool("write", true, "Whether or not to write output to file")
	formatType := flag.String("format", common.Plaintext, "the dataset file format")
	
	fpType := flag.String("fp", common.MinhashFP, "Fingeprinting method")
	jaccardThreshold := flag.Float64("jaccard", 0.05, "Jaccard index threshold for similarity")
	shingleSize := flag.Int("shingleSize", 7, "Length of shingle")
	parallel := flag.Bool("parallel", false, "Whether or not to run alignments in parallel with goroutines")
	runAlignment := flag.Bool("align", true, "Whether or not to run the alignment/correction phases")
	winnowingWindow := flag.Int("t", 15, "Size of winnowing window t")
	p := flag.Int("p", 5, "P to mod by when using modp")
	flag.Parse()
	
	flags.WriteOutput = *writeOutput
	flags.DirName = *dirName
	flags.OutDir = *outDir
	flags.FormatType = *formatType
	flags.FpType = *fpType
	flags.ShingleSize = * shingleSize
	flags.JaccardThreshold = *jaccardThreshold
	flags.Parallel = *parallel
	flags.RunAlignment = * runAlignment
	flags.WinnowingWindow = *winnowingWindow
	flags.P = *p
	flags.Groundtruth = *groundTruth
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
		var alignments map[string]common.Alignment
		var alignmentsPerDocument  map[string][]string
		if flags.Parallel {
			alignments, alignmentsPerDocument = alignment.AlignParallel(documentAdjacencyList, docList)
		} else {
			alignments, alignmentsPerDocument = alignment.AlignSerial(documentAdjacencyList, docList)
		}
		scoreSum := 0.0
		for _, al := range alignments {
			scoreSum += al.Score 
		}

		fmt.Printf("Score sum: %5.1f \n", scoreSum)
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
		
		fmt.Printf("Out of %d the corrected documents, mean edit distance improved from %5.2f to %5.2f \n", len(correctedDocs), originalStats.MeanInCorrected, correctedStats.MeanInCorrected)

		fmt.Println(originalStats)
		fmt.Println(correctedStats)
	}
}

