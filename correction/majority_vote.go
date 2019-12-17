package correction

import (
    "postCorr/common"
    "postCorr/queries"
)


/**
*   Performs a majority vote across all parts of the alignment
*   If indices were counted as aligning, they are used in the vote
*   The relationship between alignments in a cluster is such that
*   the primary alignment region is very similar in both
*   Also eturns an integer representing the number of corrections made
**/

func MajorityVote(cluster cluster) (common.Document, []rune, int){
    noCorrections := 0;
    primAlign, _ := queries.GetAlignmentByID(common.AlignmentIndex, cluster.PrimaryAlignment)
    docToCorrect, _ := queries.GetDocByID(common.DocumentIndex, primAlign.PrimaryDocumentID)
    
    correctedDocText := make([]rune, len(docToCorrect.Text))
    copy(correctedDocText, docToCorrect.Text)
    docs := map[string][]rune{}
    for docID,_ := range cluster.DocumentIDSet {
        doc,_ := queries.GetDocByID(common.DocumentIndex, docID)
        docs[docID] = doc.Text
    }
    for id, mapping := range cluster.Mappings{
        cluster.Mappings[id] = invertMap(mapping)
    }
    for _, ind := range primAlign.PrimaryAl{
        counts := map[rune]int{} 
        counts[docToCorrect.Text[ind]] = 1
        max := 1
        maxRune := docToCorrect.Text[ind]
        comparisons := 1
        for id, mapping := range cluster.Mappings{
            if val, exists := mapping[ind]; exists{
                comparisons += 1
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
        if (maxRune != docToCorrect.Text[ind]){
            noCorrections += 1
        }
        
        correctedDocText[ind] = maxRune
    }
    return docToCorrect, correctedDocText, noCorrections
}

func invertMap(m map[int]int) map[int]int{
    invertedMap := map[int]int{}
    
    for k, v := range m {
        invertedMap[v] = k
    }
    return invertedMap
}
