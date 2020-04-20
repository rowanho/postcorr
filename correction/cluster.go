package correction

import (
	"postCorr/common"
	"postCorr/flags"
	"postCorr/readWrite"

	"fmt"
	"path"

	"github.com/rowanho/levenshtein"
)

type alignMap = struct {
	Mapping             map[int]int
	PrimaryDocumentID   string
	SecondaryDocumentID string
	Start               int
	End                 int
}

func modifyText(primaryDocumentID string, text []rune) []rune{
	var groundText []rune
	if flags.Logging && flags.Groundtruth != "" {
			groundText, _ = readWrite.ReadRunes(path.Join(flags.Groundtruth, primaryDocumentID))
	}
	subEdits := make(map[int]string)
	delEdits := make(map[int]string)
	newText := make([]rune, 0)
	endPoint := 0
	modified := false
	sub := true
	for i := 0; i < len(text); i++ {
		if _, exists := removeIndices[primaryDocumentID][i]; exists {
			modified = true
			sub = false
			endPoint = len(newText)
		} else if _, exists := editIndices[primaryDocumentID][i]; exists {
			modified = true
			sub = true
			endPoint = i + 1
			newText = append(newText, editIndices[primaryDocumentID][i])
		} else {
			modified = false
			newText = append(newText, text[i])
		}

		if flags.Logging && flags.Groundtruth != "" && modified {
			before := levenshtein.ComputeDistance(groundText, append(newText[:endPoint], text[endPoint:]...))
			after := levenshtein.ComputeDistance(groundText, append(newText[:endPoint+1], text[endPoint+1:]...))
			fmt.Println(before, after)
			if sub {
				if before < after {
					subEdits[endPoint] = "worse"
					} else if before == after {
						subEdits[endPoint] = "same"
					} else {
						subEdits[endPoint] = "better"
					}
			} else {
				if before < after {
					delEdits[i] = "worse"
					} else if before == after {
						delEdits[i] = "same"
					} else {
						delEdits[i] = "better"
					}
			}
		}
	}
	substitutionGraph[primaryDocumentID] = subEdits
	deletionGraph[primaryDocumentID] = delEdits
	return newText
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
			noCorrections := MajorityVote(primaryDocumentID, alignmentMaps, documents, docMap)
			correctedDocText := modifyText(primaryDocumentID, documents[docMap[primaryDocumentID]].Text)
			documents[docMap[primaryDocumentID]].Text = correctedDocText
			totalCorrections += noCorrections
			if noCorrections > 0 {
				correctedDocs[primaryDocumentID] = true
			}
		}
	}
	if flags.WriteOutput {
		for  docID := range correctedDocs {
			readWrite.PlaintextWrite(docID, documents[docMap[docID]].Text)
		}
	}



	if flags.Logging {
		readWrite.SerialiseVote(reuseGraph)
		readWrite.SerialiseStartEnds(oldStartEndGraph, "old")
		readWrite.SerialiseStartEnds(reuseStartEndGraph, "new")
		readWrite.SerialiseEdits(substitutionGraph, "sub")
		readWrite.SerialiseEdits(deletionGraph, "del")
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
