package fingerprinting


import (
    "crypto/md5"
)



/**
* Function ComputeHashes
* @parameter text - the string to hash
* @parameter hashType - string specifying what hash type to use
* @returns string - the hashed result
*/

func ComputeHash(text string, hashType string) []byte {
    var res []byte
    switch hashType {
    case "md5":
        res = computeMd5(text)
    }
    return res
}

/**
* Function ComputeMd5
* @parameter text - The string to hash
* @returns string - The md5 hash of the string
*/
func computeMd5(text string) []byte {
    h := md5.New()
    h.Write([]byte(text))
    return h.Sum(nil)
}
