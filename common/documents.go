package common


// Represents an entire document, in such a way we can reconstruct the original OCR representation
type Document struct {
    // textComponents maps a component id to a string
    textComponents map[string][]rune    
}


// Member functions 

func AllStrings(doc *Document) []rune {
    
    all := make([]rune, 0)    
    for _, s := range doc {
        all = append(all, s...)
    }
    
    return all
}