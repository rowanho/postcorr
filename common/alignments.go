package common


/**
* An alignment lines up a primary string with multiple secondary strings
* A primary alignment has a common start and end, but the alignments with different regions
* may be slightly different, hence the alignment has slightly different
*
* For now, we are looking at alignments contained within single documents,
* but spanning multiple components
**/

type Alignment = struct {
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
