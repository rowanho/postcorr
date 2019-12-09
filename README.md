# OCR Post Correction
Post OCR correction through local text reuse detection (work in progress!)

## Requirements 

### Standing up elasticsearch

The program requires an instance of elasticseach 7.4 running on port 9200
https://www.elastic.co/guide/en/elasticsearch/reference/7.4/install-elasticsearch.html


## Structure of the Project

The code is structured into folders which are all parts of the golang module "postCorr".
Each folder is a package of the folder's name.

### Alignment 

Contains functions to perform local alignment and suboptimal alignments between
two regions of text. These are implemented from the ground up.

### Common

Contains common objects (structs) used throughout the other parts of the codebase.
This includes efficient ways to represent documents once they are converted out of various
OCR formats, and ways to represent clusters of alignments.

### Queries

Functions that connect to elasticsearch and performing indexing, aggregations and search.
Makes use of Olivere's excellent client https://github.com/olivere/elastic
 
### Fingerprinting
 
 Contains code for making representation of documents/strings into hashed numerical values. Plans to include more advanced
 techniques like locality sensitive hashing later into the project
 
### Readwrite

Functions that help convert from OCR format to our common document format and back again. Plans to support some commonly used formats like 
ALTO and PAGE.

### Correction

Plans to contain implementation of different consensus methods to produce a common output from an alignment cluster.
