# OCR Post Correction
Post OCR correction through local text reuse detection (work in progress!)


## Structure of the Project

The code is structured into folders which are all parts of the golang module "postCorr".
Each folder is a package of the folder's name.

### Alignment 

Contains functions to perform local alignment(smith waterman algorithm + a faster heuristic algorithm) and suboptimal alignments between two regions of text. Also includes the code that 'clusters' alignments together to produce multiple witnesses to correct a passage.

### Common

Contains common objects used throughout the other parts of the codebase, and constants.

### Fingerprinting
 
 Contains code for making representation of documents/strings into fingerprints - sets of hashed numerical values. 
 
### Readwrite

Functions that help convert from OCR format to our common document format and back again. Plans to support some commonly used formats like 
ALTO and PAGE.

### Correction

Contains implementation of a consensus method to produce a common output from an alignment cluster.
