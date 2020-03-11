package main

import (
    "postCorr/alignment"
    "postCorr/flags"
    "postCorr/fingerprinting"
    "postCorr/readWrite"
    "postCorr/common"
    
    "fmt"
    "testing"
)

func setup(length int) {
    flags.ShingleSize = 5
    flags.SimilarityProportion = 1.0
    flags.NumAligns = 1
    flags.FpType = common.ModFP
    flags.DirName = fmt.Sprintf("synthetic_data/benchmark_align/%d_chars/err", length)
    flags.P = 1
    fingerprinting.ResetRuntime()
}


// Benchmarking heuristic affine alignment
func benchmarkAlignment(b *testing.B, length int, affine bool, heuristic bool) {
    setup(length)
    flags.Affine = affine
    flags.FastAlign = heuristic
    docList, _ := readWrite.TraverseDocs()
    documentAdjacencyList := fingerprinting.GetSimilarDocuments(docList)
    b.ResetTimer()
    for n := 0; n < b.N; n++ {
        alignment.AlignSerial(documentAdjacencyList, docList)
    }
}
/**

func BenchmarkRegular(b *testing.B) {
    benchmarkAlignment(b, 3000, false, false)
}
**/
func BenchmarkHeuristic(b *testing.B) {
    benchmarkAlignment(b, 2500, false, true)
}

/**
func BenchmarkAffine(b *testing.B) {
    benchmarkAlignment(b, 3000, true, false)
}
**/
func BenchmarkAffineHeuristic(b *testing.B) {
    benchmarkAlignment(b, 2500, true, true)
}


