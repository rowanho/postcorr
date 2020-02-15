package readWrite

import (
    "postCorr/common"
    "postCorr/flags"
    
    "encoding/json"
    "io/ioutil"
    "path"
    "os"
    "fmt"
)

type Edge = struct {
    DocumentID string
    Score int
}

func SerialiseGraph(alignments map[string]common.Alignment, alignmentsPerDocument map[string][]string) {
    
    graphMap := make(map[string][]Edge)
    
    for docId, alIds := range alignmentsPerDocument {
        for _, alId := range alIds {
            e := Edge {
                DocumentID : alignments[alId].SecondaryDocumentID,
                Score : alignments[alId].Score,
            }
            
            if _, exists := graphMap[docId]; !exists {
                graphMap[docId] = []Edge{e}    
            } else {
                
                graphMap[docId] = append(graphMap[docId], e)
            }
        }
    }
    
    bytes, _ := json.Marshal(graphMap)
    fn := fmt.Sprintf("%s_graph%d.json",flags.DirName, flags.ShingleSize)
    ioutil.WriteFile(path.Join(flags.LogDir, fn), bytes, 0644)
}



func SerialiseJaccards(scores []float64) {
    os.Mkdir(flags.LogDir, os.ModePerm)
    fn := fmt.Sprintf("%s_jaccard_indexes%d.txt", flags.DirName, flags.ShingleSize)
    f, _ := os.Create(path.Join(flags.LogDir, fn))
	defer f.Close()
	
	for _, j := range scores {
		f.WriteString(fmt.Sprintf("%f", j) + "\n")
	}    
}