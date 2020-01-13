package correction

import (
	"postCorr/common"

	"postCorr/readWrite"
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

func ClusterAndCorrectAlignments(clustersList [][]string, alignments map[string]common.Alignment, documents []common.Document, docMap map[string]int) int {

	totalCorrections := 0
	// Loop through the adjancency list
	for _, clusterList := range clustersList {
		// Our key alignment is the 'master' alignment, we produce a cluster centred around it
		// Attempt to correct the primary alignment in the master
		if len(clusterList) >= 1 {
			alignmentMaps := make([]alignMap, len(clusterList))
			primaryDocumentID := alignments[clusterList[0]].PrimaryDocumentID
			for i, alignmentId := range clusterList {
				alignmentMaps[i] = getAlignmentMap(alignments[alignmentId])
			}
			correctedDocText, noCorrections := MajorityVote(primaryDocumentID, alignmentMaps, documents, docMap)
			totalCorrections += noCorrections
			readWrite.PlaintextWrite(primaryDocumentID, correctedDocText)
		}
	}

	return totalCorrections
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
