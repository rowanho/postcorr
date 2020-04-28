# OCR Post Correction
Post OCR correction through local text reuse detection (work in progress!)


## Structure of the Project

The code is structured into folders which are all parts of the golang module "postCorr".
Each folder is a package of the folder's name.

### Alignment 

Contains functions to perform local alignment(smith waterman algorithm + a faster heuristic algorithm) and suboptimal alignments between two regions of text. Also includes the code that 'clusters' alignments together to produce multiple witnesses to correct a passage.

### Fingerprinting
 
 Contains code for making representation of documents/strings into fingerprints - sets of hashed numerical values. 
 
### Readwrite

Functions that help convert from OCR format to our common document format and back again. Plans to support some commonly used formats like 
ALTO and PAGE.

### Correction

Contains implementation of a consensus method to produce a common output from an alignment cluster.

## python-utils

Contains python code for implementing a language model, and 


## Running the Code

## Basic Setup With Golang
All instructions are for mac/linux. I can't guarantee there won't be issues on windows, but nothing is hardcoded which prevents windows use.
### Installing and Building
* Install go: https://golang.org/doc/install
* Build the code with `go build`, this will create an executable called **postCorr**
### Command Line Parameters
* Run the executable with the relevant command line flags.
* To view flags and their descriptions, run `./postCorr -h`
## Using Additional Python Based Parts Of the Codebase

### Installing dependencies
* Install a new python 3.7 environment
* Within the environment, install dependencies with `pip install -r requirements.txt` 


### Running the Language Model Code 
nltk is the relevant dependency here, and requires some extra steps to get data for training.
* run the following inside a shell of the python environment 
```
>>> import nltk
>>> nltk.download('reuters')
>>> nltk.download('brown')
```
* Next, in the base directory of the code, run the model with `python python-utils/language_model/language_model.py`, and wait until the code finishes training and displays a server ready message.
* We can now run the **postCorr** executable with the flag `-useLM=true`

### Running the Browser Based Visualisation Tool
The browser based tool relies on the a run of the main program being completed with the *logging* flag being set to true. After this, the usage steps are as follows:

* Run the python server with `python python-utils/visualiser/server.py *input-directory*`, where \*input\* directory is the name of the directory containing the original OCR data.
* To view a file in the browser tool, navigate to `localhost:3000?filename=*file*`, where \*file\* is the path to the file within the OCR data directory.

