package getClusters



/**
* Function Kgram - the simplest overlap methods
* @parameter text - The string to turn into fingerprints
* @parameter windowSize - the size of the sliding window to use
* @parameter hashType - the type of hashing function to use
* @returns array of byte arrays - the array of fingerprints
*/

func Kgram(text string, windowSize int, hashType string) [][]byte{

    fingerprints := make([][]byte, 0)

    for i := 0; i + windowSize < len(text); i++ {
        fp := ComputeHash(text[i:i + windowSize], hashType)
        fingerprints = append(fingerprints, fp)
    }

    return fingerprints
}
