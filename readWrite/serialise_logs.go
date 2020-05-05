package readWrite

import (
	"postCorr/common"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

type Edge = struct {
	DocumentID string `json:"docId"`
	Score      int    `json:"score"`
}

func SerialiseGraph(alignments map[string]common.Alignment, alignmentsPerDocument map[string][]string) {

	graphMap := make(map[string][]Edge)

	for docId, alIds := range alignmentsPerDocument {
		for _, alId := range alIds {
			e := Edge{
				DocumentID: alignments[alId].SecondaryDocumentID,
				Score:      alignments[alId].Score,
			}

			if _, exists := graphMap[docId]; !exists {
				graphMap[docId] = []Edge{e}
			} else {
				graphMap[docId] = append(graphMap[docId], e)
			}
		}
	}

	bytes, _ := json.Marshal(graphMap)
	fn := "graph.json"
	ioutil.WriteFile(path.Join(common.LogDir, fn), bytes, 0644)
}

func SerialiseJaccards(scores []float64) {
	os.Mkdir(common.LogDir, os.ModePerm)
	fn := "jaccard_indexes.txt"
	f, _ := os.Create(path.Join(common.LogDir, fn))
	defer f.Close()

	for _, j := range scores {
		f.WriteString(fmt.Sprintf("%f", j) + "\n")
	}
}

func SerialiseVote(r map[string][]map[string]string) {
	bytes, _ := json.Marshal(r)
	fn := "vote_graph.json"
	ioutil.WriteFile(path.Join(common.LogDir, fn), bytes, 0644)
}

func SerialiseStartEnds(r map[string][]map[string]int, suffix string) {
	bytes, _ := json.Marshal(r)
	fn := "vote_start_ends_" + suffix + ".json"
	ioutil.WriteFile(path.Join(common.LogDir, fn), bytes, 0644)
}

func SerialiseEdits(e map[string]map[int]string, suffix string) {
	bytes, _ := json.Marshal(e)
	fn := "edit_graph_" + suffix + ".json"
	ioutil.WriteFile(path.Join(common.LogDir, fn), bytes, 0644)
}

func SerialiseMVote(e map[string]map[int]common.Vote) {
	bytes, _ := json.Marshal(e)
	fn := "vote_details.json"
	ioutil.WriteFile(path.Join(common.LogDir, fn), bytes, 0644)
}

func SerialiseDirname() {
	fn := "dirname.txt"
	bytes := []byte(outdir)
	ioutil.WriteFile(path.Join(common.LogDir, fn), bytes, 0644)
}
