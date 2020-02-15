package alignment

import (
	"postCorr/common"
	
)

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
			if overlapping[i] == -1 {
				for j := i +1; j < len(currentClusters); j++ {
					if overlapping[j] == -1 && anyOverlap(alignments, c1, currentClusters[j]) {
						overlaps = true
						overlapping[i] = j
						overlapping[j] = i
						break
					}
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
		
        clusterList = append(clusterList, clusters...)
    }    
    
    return clusterList
}
