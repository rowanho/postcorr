package alignment

import (
	"testing"
	//"fmt"
)

var tests = []struct {
	s1    string
	s2    string
	score int
}{
	// insertion
	{"car", "carssss", 3.0},
	// substitution
	{"", "hello", 0.0},
	{"hello", "", 0.0},
	// two empties
	{"", "", 0.0},
	// unicode stuff!
	{"你好 再见", "你好 再见", 5.0},
	{"gatcatc", "atcgatc", 5.5},
}

// Smith-Waterman
func TestSmithWaterman(t *testing.T) {
	matchReward := 1.0
	gapCost := 0.5
	for _, tt := range tests {
		score, _, _ := SmithWaterman(matchReward, gapCost, []rune(tt.s1), []rune(tt.s2))
		if score != tt.score {
			t.Errorf("SmithWaterman('%s', '%s') = %v, want %v", tt.s1, tt.s2, score, tt.score)
		}
	}
}
