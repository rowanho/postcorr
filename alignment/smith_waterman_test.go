package alignment

import (
    "testing"
    //"fmt"
)
var tests = []struct {
    s1 string
    s2 string
    score float64
} {
        // insertion
    	{"car", "carssss", 3.0},
    	// substitution
    	{"library", "librari", 6.0},
    	// deletion
    	{"library", "librar", 6.0},
    	// transposition
    	{"library", "librayr", 5.5},
    	// one empty, left
    	{"", "library", 0.0},
    	// one empty, right
    	{"library", "", 0.0},
    	// two empties
    	{"", "", 0.0},
    	// unicode stuff!
    	{"Schüßler", "Schübler", 6.0},
    	{"Ant Zucaro", "Anthony Zucaro", 8.0},
    	{"Schüßler", "Schüßler", 8.0},
    	{"Schßüler", "Schüßler", 6.0},
    	{"Schüßler", "Schüler", 6.5},
    	{"Schüßler", "Schüßlers", 8.0},
        {"ggajppp","ggafffppp", 4.0},
    
    
}

// Smith-Waterman
func TestSmithWaterman(t *testing.T) {
    matchReward := 1.0
    gapCost := 0.5
	for _, tt := range tests {
		score, _ , _ := SmithWaterman(matchReward, gapCost,[]rune(tt.s1), []rune(tt.s2))
		if score != tt.score {
			t.Errorf("SmithWaterman('%s', '%s') = %v, want %v", tt.s1, tt.s2, score, tt.score)
		}
	}
}