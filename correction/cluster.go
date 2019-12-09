package correction

import (
    "postCorr/common"
    "postCorr/readWrite"
    "postCorr/queries"
    
    "fmt"
)

type cluster struct {
    // The set of alignments we include
    PrimaryAlignment string
    AlignmentSet map[string]bool
    DocumentIDSet map[string]bool
    Mappings map[string]map[int]int
    // Map alignment id to doc id
    DocIDOfMapping map[string]string
} 


func NewCluster(key string) cluster {
    cl := cluster{ 
        PrimaryAlignment: key,
        AlignmentSet: map[string]bool{},
        DocumentIDSet: map[string]bool{},
        Mappings: map[string]map[int]int{},
        DocIDOfMapping: map[string]string{},
    }
    
    keyAlignment,_ := queries.GetAlignmentByID(common.AlignmentIndex,key)

    cl.Mappings[key] = alignmentMap(keyAlignment.PrimaryAl, keyAlignment.SecondaryAl)
    cl.DocIDOfMapping[key] = keyAlignment.SecondaryDocumentID
    cl.AlignmentSet[key] = true
    cl.DocumentIDSet[keyAlignment.PrimaryDocumentID] = true
    cl.DocumentIDSet[keyAlignment.SecondaryDocumentID] = true
    return cl
}
/**
* There needs to be a function here that takes in the alignment graph and produces clusters
* We can ideally produce 1 cluster per alignment, if it's too small, we can stop
* The max distance level is how far we want to traverse the neighbours of the master's neighbours
* High max distances can lead to worse time complexity
**/

func ClusterAndCorrectAlignments (alignmentAdjacencyList map[string][]string, maxDistance int) {
    
    closeKeySet := map[string]bool{}
    // Loop through the adjancency list
    for key := range alignmentAdjacencyList{
        // Our key alignment is the 'master' alignment, we produce a cluster centred around it
        // Attempt to correct the primary alignment in the master
        if closeKeySet[key]{
            continue;
        }
        
        closeKeySet[key] = true
        cl := NewCluster(key)
        cl.recBuildCluster(alignmentAdjacencyList, maxDistance, closeKeySet, key, cl.Mappings[key])
        docToCorrect, correctedDocText := MajorityVote(cl)
        correctedDoc := common.Document{
            ID: docToCorrect.ID,
            Text: correctedDocText,
            ComponentLengths: docToCorrect.ComponentLengths,
        }
        fmt.Println(correctedDoc)
        fmt.Println("Plaintext write")
        readWrite.PlaintextWrite(correctedDoc.ID, correctedDoc)
    }
}


/**
* Recursively builds up our cluster
**/
func (cl cluster) recBuildCluster(alignmentAdjacencyList map[string][]string, maxDistance int, 
                                  closeKeySet map[string]bool, key string, mappings map[int]int){
    if maxDistance == 0 {
        return
    }
    
    keyAlignment,_ := queries.GetAlignmentByID(common.AlignmentIndex,key)
    cl.DocumentIDSet[keyAlignment.PrimaryDocumentID] = true
    for _, id := range alignmentAdjacencyList[key] {
        // Add something to the cluster here
        if cl.AlignmentSet[id] {
            continue;
        }
        cl.AlignmentSet[id] = true
        connectedAlignment, _ := queries.GetAlignmentByID(common.AlignmentIndex,id)
        cl.DocumentIDSet[connectedAlignment.PrimaryDocumentID] = true
        cl.DocumentIDSet[connectedAlignment.SecondaryDocumentID] = true
        if keyAlignment.PrimaryDocumentID == connectedAlignment.PrimaryDocumentID {
            newMappings := alignmentMap(connectedAlignment.PrimaryAl, connectedAlignment.SecondaryAl)
            cl.Mappings[connectedAlignment.ID] = newMappings
            cl.DocIDOfMapping[connectedAlignment.ID] = connectedAlignment.SecondaryDocumentID
            cl.recBuildCluster(alignmentAdjacencyList, maxDistance, closeKeySet, id, newMappings)
        } else {
            var newMappings map[int]int
            for i, ind := range connectedAlignment.PrimaryAl {
                if _, exist := mappings[ind]; exist {
                    newMappings[connectedAlignment.SecondaryAl[i]] = mappings[ind]
                }
            }
            cl.Mappings[connectedAlignment.ID] = newMappings
            cl.DocIDOfMapping[connectedAlignment.ID] = connectedAlignment.SecondaryDocumentID
            cl.recBuildCluster(alignmentAdjacencyList, maxDistance - 1, closeKeySet, id, newMappings)
        }
    } 
}

func alignmentMap(al1 []int, al2 []int) map[int]int {
    m := map[int]int{}
    for i, ind := range(al2) {
        m[ind] = al1[i]
    }
    return m
}
