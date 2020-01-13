package common

/**
* An alignment lines up a primary string a secondary string
* Type Alignment has indexes for the full string of the document
**/

type Alignment = struct {
	ID string `json:"id"`

	Score float64 `json:"score"`

	PrimaryAl         []int  `json:"primaryAl"`
	PrimaryDocumentID string `json:"primaryDocumentID"`
	PrimaryStartIndex int    `json:"primaryStartIndex"`
	PrimaryEndIndex   int    `json:"primaryEndIndex"`

	SecondaryAl         []int  `json:"secondaryAl"`
	SecondaryDocumentID string `json:"secondaryDocumentID"`

	SecondaryStartIndex int `json:"secondaryStartIndex"`
	SecondaryEndIndex   int `json:"secondaryEndIndex"`
}
