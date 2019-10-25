package fingerprinting



/**
* Function ModP - Simple overlap fingerprinting with downsampling
* parameter text - The string to turn into fingerprints
* parameter windowSize - the size of the sliding window to use
* parameter 
* parameter hashType - the type of hashing function to use
* returns array of byte arrays - the array of fingerprints
*/


func ModP(text string, windowSize int, p uint64) []uint64{
    fingerprints := make([]uint64, 0)
    for i := 0; i + windowSize < len(text); i++ {
        fp := ComputeFNV64(text[i:i + windowSize])
        // Apply mod, check if 0
        if fp % p == 0 {
            fingerprints = append(fingerprints, fp)
        }
    }
    return fingerprints
}
