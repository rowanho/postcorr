package common

// Represents an entire document, in such a way we can reconstruct the original OCR representation
type Document struct {
	ID               string `json:"id"`
	Text             []rune `json:"text"`
	ComponentLengths []int  `json:"componentLengths"`
}
