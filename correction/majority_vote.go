package correction

import (
    "postCorr/common"
    "postCorr/queries"
    
    "fmt"
)


/**
*   Performs a majority vote across all parts of the alignment
*   If indices were counted as aligning, they are used in the vote
*   The relationship between alignments in a cluster is such that
*   the primary alignment region is very similar in both
*   Also eturns an integer representing the number of corrections made
**/

func MajorityVote(cluster cluster) (map[string][]rune, int){
    noCorrections := 0;
    primAlign, _ := queries.GetAlignmentByID(common.AlignmentIndex, cluster.PrimaryAlignment)
    
    docs := map[string][]rune{}
    for docID,_ := range cluster.DocumentIDSet {
        doc,_ := queries.GetDocByID(common.DocumentIndex, docID)
        docs[docID] = doc.Text
    }

    for _, ind := range primAlign.PrimaryAl {
        counts := map[rune]int{} 
        counts[docs[cluster.PrimaryDocId][ind]] = 1
        max := 1
        maxRune := docs[cluster.PrimaryDocId][ind]
        for id, mapping := range cluster.Mappings{
            if val, exists := mapping[ind]; exists{
                r := docs[cluster.DocIDOfMapping[id]][val]
                _, ok := counts[r]
                if ok == true{
                    counts[r] += 1
                } else{
                    counts[r] = 1
                }
                
                if counts[r] > max {
                    max = counts[r]
                    maxRune = r
                }
            }
        }
        
        for id, mapping := range cluster.Mappings {
            if val, exists := mapping[ind]; exists {
                r := docs[cluster.DocIDOfMapping[id]][val]
                if maxRune != r {
                  noCorrections += 1
                  docs[cluster.DocIDOfMapping[id]][val] = maxRune
                }
            }
        }
        
        if docs[cluster.PrimaryDocId][ind] != maxRune {
          fmt.Println("changed")
          docs[cluster.PrimaryDocId][ind] = maxRune
          noCorrections += 1
        }
        fmt.Println(counts)
    }
    return docs, noCorrections
}
