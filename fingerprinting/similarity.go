package fingerprinting

// Computes the jaccardindex of the two sets of fingerprints
func FpJaccardScore(fp1 map[uint64]int,  fp2 map[uint64]int) float64 {
    
    intersection := 0
    union := 0
    
    // Iterate over the hashes
    for hash, count := range fp1 {
        if count2, ok := fp2[hash]; ok {
            if count > count2 {
                union += count - count2
                intersection += count2
            } else {
                union += count2 - count
                intersection += count
            }
        } else {
            union += count;
        }
    }
    
    for hash2, count2 := range fp2 {
        if _, ok := fp1[hash2]; !ok {
            union += count2;
        }
    }
    
    return float64(intersection) / float64(union)
}