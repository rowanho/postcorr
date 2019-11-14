package common


// Represents an entire document, in such a way we can reconstruct the original OCR representation
type Document struct {
    ID string  
    // textComponents maps a component id to a string
    TextComponents map[string][]rune  
    // The component order gives a 'reading order' to the component IDs
    // In some datasets we may not be able to fully complete component orders
    ComponentOrder []string
}


// Member functions 

func (doc Document) AllStrings() []rune {
    
    all := make([]rune, 0)    
    for _, s := range doc.TextComponents {
        all = append(all, s...)
    }
    
    return all
}