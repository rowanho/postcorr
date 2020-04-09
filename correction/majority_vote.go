package correction

import (
	"postCorr/common"
	"postCorr/flags"
	"postCorr/readWrite"

	"strings"
	"unicode"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"math"
	"strconv"
    "path"
    "fmt"

	"github.com/rowanho/levenshtein"
)

var words = []string{}
var n = 1
var reuseGraph = make(map[string][]map[string]string)
var prevCount = 0
var correctionGraph = make(map[string][]map[int]string)

/**
*   Performs a majority vote across all parts of the alignment
*   If indices were counted as aligning, they are used in the vote
*   The relationship between alignments in a cluster is such that
*   the primary alignment region is very similar in both
*   Also eturns an integer representing the number of corrections made
**/

func MajorityVote(primaryDocumentID string, alignmentMaps []alignMap, documents []common.Document, docMap map[string]int) ([]rune, int) {
	threshold := -math.Log2(flags.LmThreshold)
	noCorrections := 0
	maxEnd := 0
	minStart := 100000000

	for _, alMap := range alignmentMaps {
		if alMap.Start < minStart {
			minStart = alMap.Start
		}
		if alMap.End > maxEnd {
			maxEnd = alMap.End
		}
	}
	var groundText []rune
	if flags.LogLevel > 1 && flags.Groundtruth != "" {
    	groundText, _ = readWrite.ReadRunes(path.Join(flags.Groundtruth, primaryDocumentID))
	}


	primText := documents[docMap[primaryDocumentID]].Text
	if flags.UseLM {
		words = wordsBeforePoint(primText, minStart, n)
	}
	var currentWord string
	if len(words) > 0 {
		currentWord = words[len(words) -1]
	}
	requiresNewWord := false
	reuseEdits := make(map[int]string)

	for ind := minStart; ind < maxEnd; ind++ {
		if flags.UseLM {
			if unicode.IsSpace(primText[ind]) {
				requiresNewWord = true
			} else if requiresNewWord {
				currentWord = getCurrentWord(primText, ind)
				words = append(words, currentWord)
				requiresNewWord = false
			}
		}
		
		numVotes := 1
		counts := map[rune]int{}
		max := 1
		maxRune := primText[ind]
		counts[primText[ind]] = 1
		for _, alMap := range alignmentMaps {
			if val, exists := alMap.Mapping[ind]; exists {
				numVotes += 1
				r := documents[docMap[alMap.SecondaryDocumentID]].Text[val]
				_, ok := counts[r]
				if ok == true {
					counts[r] += 1
				} else {
					counts[r] = 1
				}

				if counts[r] > max {
					max = counts[r]
					maxRune = r
				}
			}
		}
		//fmt.Println(counts)
		//fmt.Println(primText[ind])
		var prevText []rune
		if flags.LogLevel > 1 && flags.Groundtruth != "" {
			prevText = make([]rune, len(primText))
			copy(prevText, primText)
		}

		prevNoCorrections := noCorrections
		if primText[ind] != maxRune && max > numVotes / 2 {
			if flags.UseLM && len(words) > 0 {
				end := len(words) - 1
				start := end - n
				if start < 0 {
					start = 0
				}
				joined := strings.Join(words[start:end], " ")
				score := getLmScore(currentWord, joined)
				if score != "inf" {
					f, _ := strconv.ParseFloat(score, 64)
					if f > threshold {
						primText[ind] = maxRune
						noCorrections += 1						
					} else {
						prevCount += 1
					}
				} else {
					primText[ind] = maxRune
					noCorrections += 1					
				}
			} else {
				primText[ind] = maxRune
				noCorrections += 1
			}

		}

		if prevNoCorrections < noCorrections && flags.LogLevel > 1 && flags.Groundtruth != "" {
			before := levenshtein.ComputeDistance(groundText, prevText)
			after := levenshtein.ComputeDistance(groundText, primText)
			fmt.Println(before, after)
			if before < after{
				reuseEdits[ind] = "worse"
			} else if before == after{
				reuseEdits[ind] = "same"
			} else {
				reuseEdits[ind] = "better"
			}
		}


	}
	//fmt.Println(string(primText))
	if flags.LogLevel > 0 && noCorrections > 0 {
		reuseCluster := make(map[string]string)
		p := []rune(strings.Repeat("_", maxEnd + 1 - minStart))
		for _, m := range(alignmentMaps) {
			for i := m.Start; i <= m.End; i++ {
				if _, exists := m.Mapping[i]; exists {
					p[i - minStart] = primText[i]
				}
			}
		}
		reuseCluster[primaryDocumentID] = string(p)
		for _, m := range(alignmentMaps) {
			s := m.Mapping[m.Start]
			e := m.Mapping[m.End]
			r := strings.Repeat("_", m.Start - minStart)
			t := []rune(strings.Repeat("_", e + 1 - s))
			for _, secPos := range m.Mapping {
				t[secPos - s] = documents[docMap[m.SecondaryDocumentID]].Text[secPos]
			}
			reuseCluster[m.SecondaryDocumentID] = r + string(t)
		}
		scaledReuseEdits := make(map[int]string)

		for key, val := range reuseEdits { 
			scaledReuseEdits[key - minStart] = val
		}

		reuseGraph[primaryDocumentID] = append(reuseGraph[primaryDocumentID], reuseCluster)
		correctionGraph[primaryDocumentID] = append(correctionGraph[primaryDocumentID], scaledReuseEdits)
	}
		
	return primText, noCorrections
}

func wordsBeforePoint(text []rune, pos int, n int) []string {
	words := make([]string, 0)
	wordStarts := make([]int, 0)
	hitChars := false
	for i := pos; i > -1; i -- {
		if unicode.IsSpace(text[i]) {
			if hitChars {
				wordStarts = append(wordStarts, i + 1)
				hitChars = false
				if len(wordStarts) > n {
					break;
				}
			}
		} else if !hitChars {
			hitChars = true
		}
	}
	
	for i := len(wordStarts) -1; i  > -1; i -- {
		words = append(words, getCurrentWord(text, wordStarts[i]))
	}
	return words
}

func getCurrentWord(text []rune, pos int) string {
	
	end := pos
	for i := pos; i < len(text); i++ {
		if unicode.IsSpace(text[i]) {
			break;
		}
		end ++
	}
	
	if end == len(text) {
		return string(text[pos:])
	}
	
	return string(text[pos:end + 1])
}


func getLmScore(word string, context string) string {
	requestBody, _ := json.Marshal(map[string]string {
		"word" : word,
		"sentence" : context,
	})
	
	resp, err := http.Post("http://localhost:5000/", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "0.0"
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "0.0"
	}
	s := string(body)
	return s
}