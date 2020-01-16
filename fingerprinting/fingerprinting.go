package fingerprinting

import (
	"strings"
)

func preProcess(text string) string {
	return strings.ToLower(text)
}

/**
* Function ModP - Simple overlap fingerprinting with downsampling
* parameter text - The string to turn into fingerprints
* parameter windowSize - The size of the sliding window to use
* parameter p - The mod divisor
* returns array of byte arrays - the array of fingerprints
 */

func ModP(text string, windowSize int, p int) map[uint64]bool {
	pU := uint64(p)
	fpCounts := make(map[uint64]bool)
	for i := 0; i+windowSize < len(text); i++ {
		// Apply mod, check if 0
		fp := ComputeFNV64(text[i : i+windowSize])
		if fp%pU == 0 {
			if _, exists := fpCounts[fp]; !exists {
				fpCounts[fp] = true
			}
		}
	}
	return fpCounts
}
