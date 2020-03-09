package main

import (
    "postCorr/alignment"
    "postCorr/flags"
    "postCorr/fingerprinting"
    "postCorr/readWrite"
    "postCorr/common"
    
    "testing"
)

func setup() {
    flags.ShingleSize = 5
    flags.SimilarityProportion = 1.0
    flags.NumAligns = 1
    flags.FpType = common.ModFP
    flags.DirName = "synthetic_data/benchmark_align/500_chars/err"
    flags.P = 1
    fingerprinting.ResetRuntime()
}

// Benchmarking alignment
func BenchmarkAlignment(b *testing.B) {
    setup()
    flags.Affine = false
    flags.FastAlign = false
    docList, _ := readWrite.TraverseDocs()
    documentAdjacencyList := fingerprinting.GetSimilarDocuments(docList)
    b.ResetTimer()
    for n := 0; n < b.N; n++ {
        alignment.AlignSerial(documentAdjacencyList, docList)
    }
}

// Benchmarking heuristic affine alignment
func BenchmarkHeuristic(b *testing.B) {
    setup()
    flags.Affine = false
    flags.FastAlign = true
    docList, _ := readWrite.TraverseDocs()
    documentAdjacencyList := fingerprinting.GetSimilarDocuments(docList)
    b.ResetTimer()
    for n := 0; n < b.N; n++ {
        alignment.AlignSerial(documentAdjacencyList, docList)
    }
}

// Benchmarking alignment
func BenchmarkAffineAlignment(b *testing.B) {
    setup()
    flags.Affine = true
    docList, _ := readWrite.TraverseDocs()
    documentAdjacencyList := fingerprinting.GetSimilarDocuments(docList)
    b.ResetTimer()
    for n := 0; n < b.N; n++ {
        alignment.AlignSerial(documentAdjacencyList, docList)
    }
}

// Benchmarking heuristic affine alignment
func BenchmarkAffineHeuristic(b *testing.B) {
    setup()
    flags.Affine = true
    flags.FastAlign = true

    docList, _ := readWrite.TraverseDocs()
    documentAdjacencyList := fingerprinting.GetSimilarDocuments(docList)
    b.ResetTimer()
    for n := 0; n < b.N; n++ {
        alignment.AlignSerial(documentAdjacencyList, docList)    
    }
}



