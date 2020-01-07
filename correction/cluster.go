package correction

import (
    "postCorr/common"
    "postCorr/readWrite"
    "postCorr/queries"
)

type cluster struct {
    // The set of alignments we include
    PrimaryAlignment string
    DocumentIDSet map[string]bool
    Mappings map[string]map[int]int
    // Map alignment id to doc id
    DocIDOfMapping map[string]string
    PrimaryDocId string
} 


func NewCluster(key string) cluster {
    cl := cluster{ 
        PrimaryAlignment: key,
        DocumentIDSet: map[string]bool{},
        Mappings: map[string]map[int]int{},
        DocIDOfMapping: map[string]string{},
        PrimaryDocId: "",
    }
    
    keyAlignment,_ := queries.GetAlignmentByID(common.AlignmentIndex,key)

    cl.Mappings[key] = alignmentMap(keyAlignment.PrimaryAl, keyAlignment.SecondaryAl)
    cl.DocIDOfMapping[key] = keyAlignment.SecondaryDocumentID
    cl.DocumentIDSet[keyAlignment.PrimaryDocumentID] = true
    cl.DocumentIDSet[keyAlignment.SecondaryDocumentID] = true
    cl.PrimaryDocId = keyAlignment.PrimaryDocumentID
    return cl
}

/**
* There needs to be a function here that takes in the alignment graph and produces clusters
* We can ideally produce 1 cluster per alignment, if it's too small, we can stop
* The max distance level is how far we want to traverse the neighbours of the master's neighbours
* High max distances can lead to worse time complexity
**/

func ClusterAndCorrectAlignments (alignmentAdjacencyList map[string][]string, maxDistance int) int {
    
    totalCorrections := 0
    alreadyCorrected := map[string]bool{}
    correctedDocs := map[string]bool{}
    // Loop through the adjancency list
    for key := range alignmentAdjacencyList{
        if  _, exists := alreadyCorrected[key]; exists{
            continue;
        }
        // Our key alignment is the 'master' alignment, we produce a cluster centred around it
        // Attempt to correct the primary alignment in the master
        if  _,exists := alreadyCorrected[key]; !exists && len(alignmentAdjacencyList[key]) > 2 {
          cl := NewCluster(key)
          for _, alignmentId := range alignmentAdjacencyList[key] {
            alreadyCorrected[alignmentId] = true
            alignment,_ := queries.GetAlignmentByID(common.AlignmentIndex, alignmentId)
            cl.DocumentIDSet[alignment.SecondaryDocumentID] = true
            cl.Mappings[alignmentId] = alignmentMap(alignment.PrimaryAl, alignment.SecondaryAl)
            cl.DocIDOfMapping[alignmentId] = alignment.SecondaryDocumentID
          }
          docs, noCorrections := MajorityVote(cl)
          for docId := range cl.DocumentIDSet { 
            correctedDocs[docId] = true
          }
          queries.BulkUpdateDocuments(common.DocumentIndex, docs)
          totalCorrections += noCorrections
        }
        
        alreadyCorrected[key] = true        
    }
    for docId := range correctedDocs {
      readWrite.PlaintextWrite(docId)
    }
    return totalCorrections
}


func alignmentMap(al1 []int, al2 []int) map[int]int {
    m := map[int]int{}
    for i, ind := range(al1) {
        m[ind] = al2[i]
    }
    return m
}
