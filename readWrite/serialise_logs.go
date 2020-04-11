package readWrite

import (
    "postCorr/common"
    
    "encoding/json"
    "io/ioutil"
    "path"
    "os"
    "fmt"
)

type Edge = struct {
    DocumentID string `json:"docId"`
    Score int `json:"score"`
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
    fn := fmt.Sprintf("%s_graph.json",common.LogDir)
    ioutil.WriteFile(path.Join(common.LogDir, fn), bytes, 0644)
}



func SerialiseJaccards(scores []float64) {
    os.Mkdir(common.LogDir, os.ModePerm)
    fn := fmt.Sprintf("%s_jaccard_indexes.txt", common.LogDir)
    f, _ := os.Create(path.Join(common.LogDir, fn))
	defer f.Close()
	
	for _, j := range scores {
		f.WriteString(fmt.Sprintf("%f", j) + "\n")
	}    
}


func SerialiseVote(r map[string][]map[string]string) {
    bytes, _ := json.Marshal(r)
    fn := fmt.Sprintf("%s_vote_graph.json",common.LogDir)
    ioutil.WriteFile(path.Join(common.LogDir, fn), bytes, 0644)
}

func SerialiseStartEnds(r map[string][]map[string]int) {
    bytes, _ := json.Marshal(r)
    fn := fmt.Sprintf("%s_vote_start_ends.json",common.LogDir)
    ioutil.WriteFile(path.Join(common.LogDir, fn), bytes, 0644)    
}
func SerialiseEdits(e map[string][]map[int]string) {
    bytes, _ := json.Marshal(e)
    fn := fmt.Sprintf("%s_edit_graph.json",common.LogDir)
    ioutil.WriteFile(path.Join(common.LogDir, fn), bytes, 0644)
}