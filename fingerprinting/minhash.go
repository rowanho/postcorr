package fingerprinting

import (
    "postCorr/common"
    
    minhashlsh "github.com/ekzhu/minhash-lsh"
)

var lsh *minhashlsh.MinhashLSH
var minhashSeed int64 = 2342342
var minhashSize = 100

func GetLSHObject(numHash int, threshold float64, count int) {
    lsh = minhashlsh.NewMinhashLSH(numHash, threshold, count)
}

func IndexMinHashObject(){
    lsh.Index()
}

func MinHash(key string, text string, windowSize int) common.LSH_fp {
	mh := minhashlsh.NewMinhash(minhashSeed, minhashSize)
    for i := 0; i+windowSize < len(text); i++ {
        mh.Push([]byte(text[i:i+windowSize]))
    }
    sigs := mh.Signature()
    lsh.Add(key, sigs)
    return common.LSH_fp{Signature: sigs}
}

func SameBucketIds(sigs []uint64) []string {
    similarIds :=  lsh.Query(sigs)
    returnIds := make([]string, len(similarIds))
    for i, id := range similarIds {
        returnIds[i] = id.(string)
    }
    return returnIds
}