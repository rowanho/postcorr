package main

import (
	"postcorr/alignment"
	"postcorr/common"
	"postcorr/fingerprinting"
	"postcorr/flags"
	"postcorr/iohandler"

	"fmt"
	"testing"
)

func setup(length int) {
	flags.K = 5
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
	docList, _ := iohandler.TraverseDocs()
	documentAdjacencyList := fingerprinting.GetSimilarDocuments(docList)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		alignment.AlignSerial(documentAdjacencyList, docList)
	}
}

func BenchmarkHeuristic(b *testing.B) {
	benchmarkAlignment(b, 2500, false, true)
}

func BenchmarkAffineHeuristic(b *testing.B) {
	benchmarkAlignment(b, 2500, true, true)
}
