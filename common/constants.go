package common

/**
* Define the index names here
**/

const MinHashIndex = "test_minhash_fingerprints"
const FpIndex = "test_fingerprints"
const DocumentIndex = "test_documents"
const AlignmentIndex = "test_alignments"

/**
* Define the types of files here, thus avoiding magic strings
**/

const Plaintext = "plaintext"

/**
* Define different hashing/reuse detection algorithms here
**/

const MinhashFP = "minhash"
const ModFP = "modp"

/**
* Define different alignment algorithms here
**/

const HeuristicAlignment = "blast"
const SwAlignment = "smith_waterman"
