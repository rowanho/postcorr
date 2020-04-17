package correction

import (
	"postCorr/common"
	"postCorr/flags"
	"postCorr/readWrite"

	"fmt"
)

type alignMap = struct {
	Mapping             map[int]int
	PrimaryDocumentID   string
	SecondaryDocumentID string
	Start               int
	End                 int
}

/**
* There needs to be a function here that takes in the alignment graph and produces clusters
* We can ideally produce 1 cluster per alignment, if it's too small, we can stop
* The max distance level is how far we want to traverse the neighbours of the master's neighbours
* High max distances can lead to worse time complexity
**/

func ClusterAndCorrectAlignments(clustersList [][]string, alignments map[string]common.Alignment, documents []common.Document, docMap map[string]int) (map[string]bool, int) {

	totalCorrections := 0
	correctedDocs := make(map[string]bool)
	// Loop through the cluster list
	for _, cluster := range clustersList {
		// Attempt to correct the primary document of the cluster
		if len(cluster) > 1 {
			alignmentMaps := make([]alignMap, len(cluster))
			primaryDocumentID := alignments[cluster[0]].PrimaryDocumentID
			for i, alignmentId := range cluster {
				alignmentMaps[i] = getAlignmentMap(alignments[alignmentId])
			}
			correctedDocText, noCorrections := MajorityVote(primaryDocumentID, alignmentMaps, documents, docMap)
			documents[docMap[primaryDocumentID]].Text = correctedDocText
			totalCorrections += noCorrections
			correctedDocs[primaryDocumentID] = true
		}
	}
	if flags.WriteOutput {
		for  docID := range correctedDocs {
			readWrite.PlaintextWrite(docID, documents[docMap[docID]].Text)
		}
	}

	if flags.Logging {
		readWrite.SerialiseVote(reuseGraph)
		readWrite.SerialiseStartEnds(reuseStartEndGraph)
		readWrite.SerialiseEdits(correctionGraph)
	}
	if flags.UseLM {
		fmt.Printf("Prevented %d\n", prevCount)
	}
	return correctedDocs, totalCorrections
}

func getAlignmentMap(al common.Alignment) alignMap {
	m := map[int]int{}
	for i, ind := range al.PrimaryAl {
		m[ind] = al.SecondaryAl[i]
	}
	a := alignMap{
		Mapping:             m,
		PrimaryDocumentID:   al.PrimaryDocumentID,
		SecondaryDocumentID: al.SecondaryDocumentID,
		Start:               al.PrimaryStartIndex,
		End:                 al.PrimaryEndIndex,
	}
	return a
}
