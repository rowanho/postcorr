package fingerprinting

// Computes the jaccardindex of the two sets of fingerprints
func FpJaccardScore(fp1 map[uint64]int,  fp2 map[uint64]int) float64 {
    
    intersection := 0
    union := 0
    
    // Iterate over the hashes
    for hash, _ := range fp1 {
        if _, ok := fp2[hash]; ok {
            union += 1
            intersection += 1
        } else {
            union += 1
        }
    }
    
    for hash2, _ := range fp2 {
        if _, ok := fp1[hash2]; !ok {
            union += 1;
        }
    }
    
    return float64(intersection) / float64(union)
}