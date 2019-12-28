package alignment

import (
	"postCorr/common"

	"fmt"
	"testing"

	"github.com/google/uuid"
)

var testPrimaryDocs = []common.Document{
	{
		ID: uuid.New().String(),
		TextComponents: map[string][]rune{
			"comp1": []rune("agtctagctacatcgactacgtcgatgagcatcgctagcatcgat"),
			"comp2": []rune("tctctatatatatgggg"),
			"comp3": []rune("gagagagagatctctctc"),
		},
		ComponentOrder: []string{"comp1", "comp2", "comp3"},
	},
	{
		ID: uuid.New().String(),
		TextComponents: map[string][]rune{
			"comp1": []rune("test test test test"),
			"comp2": []rune("this this this this"),
			"comp3": []rune("test this test this"),
		},
		ComponentOrder: []string{"comp1", "comp2", "comp3"},
	},
}

var testSecondaryDocs = []common.Document{
	{
		ID: uuid.New().String(),
		TextComponents: map[string][]rune{
			"comp1": []rune("gctagcatcgactagctactacaagtcatcatctaaaa"),
			"comp2": []rune("cgatcgatcacgatcgactagctagcatgtagtatatatattt"),
			"comp3": []rune("cgatcgactagcatcgatcatcagctacgtaccatgctatacgatcgat"),
		},
		ComponentOrder: []string{"comp1", "comp2", "comp3"},
	},
	{
		ID: uuid.New().String(),
		TextComponents: map[string][]rune{
			"comp1": []rune("test test test test"),
			"comp2": []rune("this this this this"),
			"comp3": []rune("test this test this"),
		},
		ComponentOrder: []string{"comp1", "comp2", "comp3"},
	},
}

func TestSuboptimalAlign(t *testing.T) {

	for _, tp := range testPrimaryDocs {
		for _, ts := range testSecondaryDocs {
			alignments := GetAlignments(1.0, 1.5, tp, ts, 3, 0.0)
			fmt.Println(alignments)
		}
	}
}
