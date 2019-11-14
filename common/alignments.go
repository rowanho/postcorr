package common

/**
* An alignment lines up a primary string with multiple secondary strings
* A primary alignment has a common start and end, but the alignments with different regions 
* may be slightly different, hence the alignment has slightly different
*
* For now, we are looking at alignments contained within single documents, 
* but spanning multiple components
**/

type AlignmentCluster = struct {
    
    
    PrimaryAl [][]int
    PrimaryDocumentID string
    PrimaryComponentIDs []string 
    
    PrimaryStartComponent string
    PrimaryEndComponent string
    PrimaryStartIndex int
    PrimaryEndIndex int
    
    SecondaryAl [][]int
    SecondaryDocumentID string
    SecondaryComponentIDs []string
}