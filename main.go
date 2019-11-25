package main

import (
	"postCorr/reader"
	"postCorr/common"
	
	"fmt"
	"flag"
)

func main() {
	dirName := flag.String("dir","test_dataset","path to dataset")
	formatType := flag.String("format", common.Plaintext, "the dataset file format")
	
	flag.Parse()
	
	execute(*dirName, *formatType)
}


func execute(dirName string, formatType string) {
	
	docIDList, docsErr := reader.TraverseAndIndexDocs(dirName, formatType)
	
	if docsErr != nil {
		fmt.Println("Error indexing documents %s", docsErr)
		return
	}
	
	fmt.Println(docIDList)
	
}


