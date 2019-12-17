package fingerprinting

import (
    "strings"
    
    "github.com/DearMadMan/minhash"
)

var m = minhash.New(128)

func MinHash(text string) minhash.Set {
    words := strings.Fields(text)
    return m.NewSet(words)
}