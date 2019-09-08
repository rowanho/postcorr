package getClusters


import (
    "crypto/md5"
    "encoding/hex"
)



/**
* Function ComputeHashes
* @parameter text - the string to hash
* @parameter hashType - string specifying what hash type to use
* @returns string - the hashed result
*/

func ComputeHash(text string, hashType string) string {
    switch hashType {
    case "md5":
        return computeMd5(text)
    }
}

/**
* Function ComputeMd5
* @parameter text - The string to hash
* @returns string - The md5 hash of the string
*/
func computeMd5(text string) string {
    h := md5.New()
    h.Write([]byte(text))
    return hex.EncodeToString(h.Sum(nil))
}
