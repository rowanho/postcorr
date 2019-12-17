package fingerprinting

import (
    "github.com/DearMadMan/minhash"
)

var m = minhash.New(10)

func MinHash(text string, windowSize int, sampleRate int) minhash.Set {
    words := make([]string, 0)
    for i := 0; i+windowSize < len(text); i++ {
        if i % sampleRate == 0 {
            words = append(words, text[i : i+windowSize])
        }
    }
    return m.NewSet(words)
}