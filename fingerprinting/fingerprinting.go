package fingerprinting

import (
	"strings"
	"math"
)

func preProcess(text string) string {
	return strings.ToLower(text)
}



func Kgrams(text string, k int) []uint64 {
	fps := make([]uint64, len(text) + 1 - k)
	for i := 0; i+k < len(text) + 1; i++ {
		// Apply mod, check if 0
		fp := ComputeFNV64(text[i : i+k])
		fps[i] = fp
	}
	return fps
	
}
/**
* Function ModP - Simple overlap fingerprinting with downsampling
* parameter text - The string to turn into fingerprints
* parameter windowSize - The size of the sliding window to use
* parameter p - The mod divisor
* returns array of byte arrays - the array of fingerprints
 */

func ModP(text string, windowSize int, p int) map[uint64]int {
	pU := uint64(p)
	fps := make(map[uint64]int)
	for i := 0; i+windowSize < len(text) + 1; i++ {
		// Apply mod, check if 0
		fp := ComputeFNV64(text[i : i+windowSize])
		if fp%pU == 0 {
			if _, exists := fps[fp]; !exists {
				fps[fp] += 1
			}
		}
	}
	return fps
}

func min(hashes []uint64) uint64{
	currentMin := uint64(math.MaxUint64)
	for _, hash := range hashes {
		if hash < currentMin {
			currentMin = hash
		}
	}
	return currentMin
}


/**
Winnowing algorithm
*/

func Winnowing(text string, k int, t int) map[uint64]int {
	fps := make(map[uint64]int)
	kgrams := Kgrams(text, k)
	windowSize := t - k + 1
	for start := 0; start < len(kgrams) - windowSize; start ++ {
		fps[min(kgrams[start:start + windowSize])] += 1
	}
	return fps
}
