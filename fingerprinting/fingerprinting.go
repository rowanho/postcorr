package fingerprinting

/**
* Function ModP - Simple overlap fingerprinting with downsampling
* parameter text - The string to turn into fingerprints
* parameter windowSize - The size of the sliding window to use
* parameter p - The mod divisor
* returns array of byte arrays - the array of fingerprints
 */

func ModP(text string, windowSize int, p uint64) map[uint64]int {

	fpCounts := make(map[uint64]int)
	for i := 0; i+windowSize < len(text); i++ {
		fp := ComputeFNV64(text[i : i+windowSize])
		// Apply mod, check if 0
		if fp%p == 0 {
			fpCounts[fp] += 1
		}
	}
	return fpCounts
}
