package flags

// Values that are changed via command line flags go here
var (
	DirName string
	Logging bool
	WriteOutput bool
	FpType string
	SimilarityProportion float64
	JaccardType string
	K int
	WinnowingT int
	P int
	Groundtruth string
	Affine bool
	FastAlign bool
	NumAligns int
	AlignThreshold int
	UseLM bool
	LmThreshold float64
	HandleInsertionDeletion bool
	LDelete int
	LInsert int
	BandWidth int
)