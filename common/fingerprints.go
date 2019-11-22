package common

type Fingerprints struct {
    DocumentID string         `json:"documentID"`
    FpCounts   map[uint64]int `json:"fpCounts"`
}