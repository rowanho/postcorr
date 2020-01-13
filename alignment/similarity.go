package alignment

import (
	"postCorr/common"
	
	"fmt"
)

func AlignParallel(documentAdjacencyList map[string]map[string]bool, docs []common.Document, docMap map[string]int) (map[string]common.Alignment, map[string][]string)  {

	alignments := make(map[string]common.Alignment, 0)
	alignmentDocIdMap := make(map[string][]string)
	
	for  primID,_ := range documentAdjacencyList{
		alignmentDocIdMap[primID] = make([]string, 0)
	}
	
	for primID, secIDs := range documentAdjacencyList {
		primDoc := docs[docMap[primID]]
		for secID, _ := range secIDs {
			if _, exists := documentAdjacencyList[secID][primID]; exists {
				delete(documentAdjacencyList[secID], primID)
			}
		}

		alignmentChannel := make(chan []common.Alignment)
		inverseAlignmentChannel := make(chan []common.Alignment)
		for secID, _ := range secIDs {
			secDoc := docs[docMap[secID]]
			go func(channel1 chan []common.Alignment, channel2 chan []common.Alignment, primDoc common.Document, secDoc common.Document, secID string) {
				alignments, inverseAlignments := GetAlignments(1.0, 2.0, primDoc, secDoc, 1, 0.0)
				channel1 <- alignments
				channel2 <- inverseAlignments
			}(alignmentChannel, inverseAlignmentChannel, primDoc, secDoc, secID)

		}
	for i := 0; i < len(secIDs); i++ {
		als := <-alignmentChannel
	  	for _, al := range als {
	        alignments[al.ID] = al
			alignmentDocIdMap[primID] = append(alignmentDocIdMap[primID], al.ID)
		}
	}
	for i := 0; i < len(secIDs); i++ {
		als := <-inverseAlignmentChannel
		for _, al := range als {
	        alignments[al.ID] = al
			p := al.PrimaryDocumentID
			alignmentDocIdMap[p] = append(alignmentDocIdMap[p], al.ID)
		}
	}
	}
  return alignments, alignmentDocIdMap
}

func AlignSerial(documentAdjacencyList map[string]map[string]bool, docs []common.Document, docMap map[string]int) (map[string]common.Alignment,map[string][]string)  {

	alignments := make(map[string]common.Alignment, 0)
  	alignmentDocIdMap := make(map[string][]string)
	for  primID,_ := range documentAdjacencyList{
		alignmentDocIdMap[primID] = make([]string, 0)
	}
	
	for primID, secIDs := range documentAdjacencyList {
		primDoc := docs[docMap[primID]]
		for secID, _ := range secIDs {
			if _, exists := documentAdjacencyList[secID][primID]; exists {
				delete(documentAdjacencyList[secID], primID)
			}
		}
    
		for secID, _ := range secIDs {
			secDoc := docs[docMap[secID]]
			als, inverseAls := GetAlignments(1.0, 2.0, primDoc, secDoc, 1, 0.0)
			for i, al := range als {
				alignments[al.ID] = al
				alignments[inverseAls[i].ID] = inverseAls[i]
				alignmentDocIdMap[primID] = append(alignmentDocIdMap[primID], al.ID)
				alignmentDocIdMap[secID] = append(alignmentDocIdMap[secID], inverseAls[i].ID)
				 
		     }
		}
	}

	return alignments, alignmentDocIdMap

}


func checkOverlap (al1 common.Alignment, al2 common.Alignment) bool{
	if al1.PrimaryStartIndex  > al2.PrimaryEndIndex {
		return false
	} else if  al1.PrimaryEndIndex < al2.PrimaryStartIndex{
		return false
	}
	return true
}

// Returns boolean if any alignment in cluster1 overlaps with ones in cluster 2
func anyOverlap(alignments map[string]common.Alignment, cluster1 []string, cluster2 []string) bool{
	for _, alId1 := range cluster1 {
		for _, alId2 := range cluster2 {
			if checkOverlap(alignments[alId1], alignments[alId2]){
				return true
			}
		}
	}
	return false
}

// Merges overlapping cluster
func merge(overlapping []int, clusters [][]string) [][]string {
	
	overlaps := 0
	for _, o := range overlapping {
		if o != -1 {
			overlaps += 1
		}
	}
	
	skip := make(map[int]bool)
	newClusters := make([][]string, len(clusters) - (overlaps/2))
	j := 0
	for i, o := range overlapping {
		if skip[i] {
			continue;
		}
		if o != -1 {
			newClusters[j]  = append(clusters[i], clusters[o]...)
			skip[o] = true
		} else {
			newClusters[j] = clusters[i]
		}
		j += 1
	}
	
	return newClusters
}


// Clusters alignments into groups that overlap
func getClusters(alignments map[string]common.Alignment, alsToCluster []string) [][]string{
	currentClusters := make([][]string, len(alsToCluster))
	for i, alId := range alsToCluster {
		currentClusters[i] = []string{alId}
	}
	
	overlaps := true
	// Loop until no more overlaps occur
	for overlaps == true {
		overlaps = false
		overlapping := make([]int, len(currentClusters))
		for i := range overlapping {
			overlapping[i] = -1
		}
		for i, c1 := range currentClusters {
			for j := i +1; j < len(currentClusters); j++ {
				if overlapping[j] == -1 && anyOverlap(alignments, c1, currentClusters[j]) {
					overlaps = true
					overlapping[i] = j
					overlapping[j] = i
					break
				}
			}
		}
		currentClusters = merge(overlapping, currentClusters)
	}
	return currentClusters
}

func GetSimilarAlignments(alignments map[string]common.Alignment, alignmentDocIdMap map[string][]string) [][]string{
    clusterList := make([][]string, 0)
    for _, alignmentIds := range alignmentDocIdMap {
        clusters := getClusters(alignments, alignmentIds)
		
		for _,c := range clusters {
			fmt.Println("cluster2")
			for _,alId := range c {
				fmt.Println(alignments[alId].PrimaryDocumentID)
			}
		}
        clusterList = append(clusterList, clusters...)
    }    
    
    return clusterList
}
