package common

// Represents an entire document, in such a way we can reconstruct the original OCR representation
type Document struct {
	ID string `json:"id"`
	// textComponents maps a component id to a string
	TextComponents map[string][]rune `json:"components"`
	// The component order gives a 'reading order' to the component IDs
	// In some datasets we may not be able to fully complete component orders
	ComponentOrder []string `json:"componentOrder"`
}

// For LSH
type DocString struct {
	ID string `json:"id"`
	Text string `json:"text"`
}


// Member functions

func (doc Document) AllStrings() []rune {

	all := make([]rune, 0)
	for _, s := range doc.TextComponents {
		all = append(all, s...)
	}

	return all
}



func (doc Document) ToDocString() DocString {
	return DocString{
		ID: doc.ID,
		Text: string(doc.AllStrings()),
	}
}

// Takes modified doc string and 'inserts' it back
func (doc Document) InsertDocString(text []rune) {
	c := 0
	for _, component := range doc.ComponentOrder {
		l := len(doc.TextComponents[component])
		doc.TextComponents[component] = text[c: c + l]
		c += l
	} 
}