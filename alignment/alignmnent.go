package alignment

import (
	"postCorr/common"
	"postCorr/flags"

	"github.com/schollz/progressbar"
)

func AlignParallel(documentAdjacencyList map[int]map[int]bool, docs []common.Document) (map[string]common.Alignment, map[string][]string) {
	bar := progressbar.New(100)
	alignments := make(map[string]common.Alignment, 0)
	alignmentDocIdMap := make(map[string][]string)

	for _, doc := range docs {
		alignmentDocIdMap[doc.ID] = make([]string, 0)
	}

	total := 0
	for _, secIDs := range documentAdjacencyList {
		total += len(secIDs)
	}
	prev := 0
	c := 0
	for primID, secIDs := range documentAdjacencyList {
		primDoc := docs[primID]
		for secID := range secIDs {
			if _, exists := documentAdjacencyList[secID][primID]; exists {
				delete(documentAdjacencyList[secID], primID)
			}
		}

		alignmentChannel := make(chan []common.Alignment)
		inverseAlignmentChannel := make(chan []common.Alignment)
		for secID := range secIDs {
			secDoc := docs[secID]
			go func(channel1 chan []common.Alignment, channel2 chan []common.Alignment, primDoc common.Document, secDoc common.Document) {
				alignments, inverseAlignments := GetAlignments(2, 4, primDoc, secDoc, flags.NumAligns, flags.AlignThreshold)
				channel1 <- alignments
				channel2 <- inverseAlignments
			}(alignmentChannel, inverseAlignmentChannel, primDoc, secDoc)

		}
		for i := 0; i < len(secIDs); i++ {
			als := <-alignmentChannel
			for _, al := range als {
				alignments[al.ID] = al
				alignmentDocIdMap[al.PrimaryDocumentID] = append(alignmentDocIdMap[al.PrimaryDocumentID], al.ID)
			}
			c += 1
			prog := (c * 100) / total
			if prog > prev {
				bar.Add(1)
				prev = prog
			}
		}

		for i := 0; i < len(secIDs); i++ {
			als := <-inverseAlignmentChannel
			for _, al := range als {
				alignments[al.ID] = al
				alignmentDocIdMap[al.PrimaryDocumentID] = append(alignmentDocIdMap[al.PrimaryDocumentID], al.ID)
			}
			c += 1
			prog := (c * 100) / total
			if prog > prev {
				bar.Add(1)
				prev = prog
			}
		}
	}
	return alignments, alignmentDocIdMap
}

func AlignSerial(documentAdjacencyList map[int]map[int]bool, docs []common.Document) (map[string]common.Alignment, map[string][]string) {
	bar := progressbar.New(100)
	alignments := make(map[string]common.Alignment, 0)
	alignmentDocIdMap := make(map[string][]string)
	for _, doc := range docs {
		alignmentDocIdMap[doc.ID] = make([]string, 0)
	}

	total := 0
	for _, secIDs := range documentAdjacencyList {
		total += len(secIDs)
	}
	total = total / 2
	prev := 0
	c := 0
	for primID, secIDs := range documentAdjacencyList {
		primDoc := docs[primID]
		for secID := range secIDs {
			if _, exists := documentAdjacencyList[secID][primID]; exists {
				delete(documentAdjacencyList[secID], primID)
			}
		}

		for secID := range secIDs {
			secDoc := docs[secID]
			als, inverseAls := GetAlignments(2, 4, primDoc, secDoc, flags.NumAligns, flags.AlignThreshold)
			for i, al := range als {
				alignments[al.ID] = al
				alignments[inverseAls[i].ID] = inverseAls[i]
				alignmentDocIdMap[al.PrimaryDocumentID] = append(alignmentDocIdMap[al.PrimaryDocumentID], al.ID)
				alignmentDocIdMap[al.SecondaryDocumentID] = append(alignmentDocIdMap[al.SecondaryDocumentID], inverseAls[i].ID)
			}
			c += 1
		}
		prog := (c * 100) / total
		if prog > prev {
			bar.Add(1)
			prev = prog
		}
	}
	return alignments, alignmentDocIdMap

}
