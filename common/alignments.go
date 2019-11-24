package common


/**
* An alignment lines up a primary string a secondary string
* Type Alignment has indexes for the full string of the document
* Type total alignment has the indices rescored per component, so we
* can reconstruct this in the original documents.
*
* For now, we are looking at alignments contained within single documents,
* but spanning multiple components
**/

type Alignment = struct {
	Score float64 `json:"score"`

	PrimaryAl             []int  `json:"primaryAl"`
	PrimaryDocumentID     string   `json:"primaryDocumentID"`
	PrimaryStartIndex     int      `json:"primaryStartIndex"`
	PrimaryEndIndex       int      `json:"primaryEndIndex"`

	SecondaryAl           []int  `json:"secondaryAl"`
	SecondaryDocumentID   string   `json:"secondaryDocumentID"`
	SecondaryComponentIDs []string `json:"secondaryComponentIDs"`
}


type TotalAlignment = struct {
	Score float64 `json:"score"`

	PrimaryAl             [][]int  `json:"primaryAl"`
	PrimaryDocumentID     string   `json:"primaryDocumentID"`
	PrimaryComponentIDs   []string `json:"primaryComponentIDs"`
	PrimaryStartComponent string   `json:"primaryStartComponent"`
	PrimaryEndComponent   string   `json:"primaryEndComponent"`
	PrimaryStartIndex     int      `json:"primaryStartIndex"`
	PrimaryEndIndex       int      `json:"primaryEndIndex"`

	SecondaryAl           [][]int  `json:"secondaryAl"`
	SecondaryDocumentID   string   `json:"secondaryDocumentID"`
	SecondaryComponentIDs []string `json:"secondaryComponentIDs"`
	
}