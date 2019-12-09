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
**/

func MajorityVote(cluster cluster) (common.Document, []rune){
    primAlign, _ := queries.GetAlignmentByID(common.AlignmentIndex, cluster.PrimaryAlignment)
    docToCorrect, _ := queries.GetDocByID(common.DocumentIndex, primAlign.PrimaryDocumentID)
    
    correctedDocText := make([]rune, len(docToCorrect.Text))
    var docs map[string][]rune
    for docID,_ := range cluster.DocumentIDSet {
        doc,_ := queries.GetDocByID(common.DocumentIndex, docID)
        docs[docID] = doc.Text
    }
    
    for i, ind := range primAlign.PrimaryAl{
        var counts map[rune]int
        max := 0
        var maxRune rune
        for id, mapping := range cluster.Mappings{
            if val, exists := mapping[ind]; exists{
                r := docs[id][val]
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
        correctedDocText[i] = maxRune
    }
    return docToCorrect, correctedDocText
}                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                               

                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      