package common

const (
/**
* Define the index names here
**/

	MinHashIndex = "test_minhash_fingerprints"
	FpIndex = "test_fingerprints"
	DocumentIndex = "test_documents"
	AlignmentIndex = "test_alignments"

/**
* Define the types of files here, thus avoiding magic strings
**/

	Plaintext = "plaintext"

/**
* Define different hashing/reuse detection algorithms here
**/

	MinhashFP = "minhash"
	ModFP = "modp"
	Winnowing = "winnowing"

/**
* Define different alignment algorithms here
**/

	HeuristicAlignment = "blast"
	SwAlignment = "smith_waterman"

/**
* Define similarity algorithms here
**/

	Jaccard = "regular"
	WeightedJaccard = "weighted"

/**
* Directory names
**/
	LogDir = "logs"
	OutDir = "corrected"
)