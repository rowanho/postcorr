package fingerprinting

import (
	"crypto/md5"
	"hash/fnv"
)

/**
* Function ComputeMd5
* parameter text - The string to hash
* returns string - The md5 hash of the string
 */
func ComputeMd5(text string) []byte {
	h := md5.New()
	h.Write([]byte(text))
	return h.Sum(nil)
}

func ComputeFNV64(text string) uint64 {
	h := fnv.New64()
	h.Write([]byte(text))
	return h.Sum64()
}
