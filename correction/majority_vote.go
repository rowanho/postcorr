package correction

import (
	"postcorr/common"
	"postcorr/flags"

	"bytes"
	"encoding/json"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
	"unicode"
	//"fmt"
)

var (
	words = []string{}
	n = 1
	reuseGraph = make(map[string][]map[string]string)
	reuseStartEndGraph = make(map[string][]map[string]int)
	oldStartEndGraph = make(map[string][]map[string]int)
	prevCount = 0
	substitutionGraph = make(map[string]map[int]string)
	deletionGraph = make(map[string]map[int]string)
	insertionGraph = make(map[string]map[int]string)

// Marks the indices for removal
	removeIndices = make(map[string]map[int]bool)
	editIndices = make(map[string]map[int]rune)
	additionIndices = make(map[string]map[int][]rune)
	deletionsAt = make(map[string]map[int]int)
	insertionsAt = make(map[string]map[int]int)
	mVoteLogs = make(map[string]map[int]common.Vote)
	newVoteLogs = make(map[string]map[int]common.Vote)

)


func isNewLineSpace(r rune) bool {
	return unicode.IsSpace(r) || r == '\t' || r == '\n'
}

func getDifferenceSoFar(primaryDocumentID string, start int) int {
	d := 0
	for key, entry := range deletionsAt[primaryDocumentID] {
		if key <= start {
			d += entry
		}
	}
	for key, entry := range insertionsAt[primaryDocumentID] {
		if key <= start {
			d -= entry
		}
	}
	return d
}

// Finds potential gaps for deletion, <= length of threshold
func applyDeletions(primaryDocumentID string, alignmentMaps []alignMap, documents []common.Document, docMap map[string]int) int {
	threshold := -math.Log2(flags.LmThreshold)
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
	deletions := 0

	primText := documents[docMap[primaryDocumentID]].Text
	if flags.UseLM {
		words = wordsBeforePoint(primText, minStart, n)
	}
	var currentWord string
	if len(words) > 0 {
		currentWord = words[len(words)-1]
	}
	requiresNewWord := false

	pairStart := 0
	gapSection := false
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
		count := 0
		notAlignedInPrim := 0
		for _, alMap := range alignmentMaps {
			if alMap.Start <= ind && alMap.End >= ind {
				if _, exists := alMap.Mapping[ind]; !exists {
					notAlignedInPrim += 1
				}
			}
			count += 1
		}

		if count > 1 && notAlignedInPrim > ( (count + 1) / 2) {
			if !gapSection {
				pairStart = ind
				gapSection = true
			}

		} else if gapSection {
			if ind-pairStart <= flags.LInsert {
				for j := pairStart; j < ind; j++ {
					if isNewLineSpace(primText[j]) {
						continue;
					}
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
								removeIndices[primaryDocumentID][j] = true
								deletions += 1
							} else {
								prevCount += 1
							}
						} else {
							removeIndices[primaryDocumentID][j] = true
							deletions += 1
						}
					} else {
						removeIndices[primaryDocumentID][j] = true
						deletions += 1
					}
				}
			}
			gapSection = false
		}
	}
	if _, exists := deletionsAt[primaryDocumentID]; !exists {
		deletionsAt[primaryDocumentID] = make(map[int]int)
	}
	deletionsAt[primaryDocumentID][maxEnd] = deletions
	return deletions
}

// Finds potential gaps for insertion, <= length of threshold
func applyInsertions(primaryDocumentID string, alignmentMaps []alignMap, documents []common.Document, docMap map[string]int) int {
	threshold := -math.Log2(flags.LmThreshold)
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
	primText := documents[docMap[primaryDocumentID]].Text
	if flags.UseLM {
		words = wordsBeforePoint(primText, minStart, n)
	}
	var currentWord string
	if len(words) > 0 {
		currentWord = words[len(words)-1]
	}
	requiresNewWord := false

	insertions := 0
	for ind := minStart; ind < maxEnd-1; ind++ {
		if flags.UseLM {
			if unicode.IsSpace(primText[ind]) {
				requiresNewWord = true
			} else if requiresNewWord {
				currentWord = getCurrentWord(primText, ind)
				words = append(words, currentWord)
				requiresNewWord = false
			}
		}
		commonStrings := make(map[string]int)
		count := 0
		for _, alMap := range alignmentMaps {
			start := -1
			end := -1
			if _, exists := alMap.Mapping[ind]; exists {
				start = alMap.Mapping[ind]
			}

			if _, exists := alMap.Mapping[ind+1]; exists {
				end = alMap.Mapping[ind+1]
			}

			if start > -1 && end-start > 1 {
				if end-start-1 <= flags.LDelete {
					count += 1
					s := documents[docMap[alMap.SecondaryDocumentID]].Text[start+1 : end]
					commonStrings[string(s)] += 1
				}
			}

		}
		prevInsertions := insertions
		mfreq := 0
		str := ""
		for s, freq := range commonStrings {
			if freq >  mfreq {
				str = s
				mfreq = freq
			}
		}

		if mfreq >= count/2 && count > 2 {
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
						additionIndices[primaryDocumentID][ind] = []rune(str)
						insertions += len(str)
					} else {
						prevCount += len(str)
					}
				} else {
					additionIndices[primaryDocumentID][ind] = []rune(str)
					insertions += len(str)
				}
			} else {
				additionIndices[primaryDocumentID][ind] = []rune(str)
				insertions += len(str)
			}
		}

		if prevInsertions < insertions {
			if _, exists := mVoteLogs[primaryDocumentID][ind]; !exists {
				mVoteLogs[primaryDocumentID][ind] = common.Vote{
					EditDict:   map[string]int{},
					InsertDict: map[string]int{},
				}
			}
			for key, val := range commonStrings {
				mVoteLogs[primaryDocumentID][ind].InsertDict[key] = val
			}
		}
	}
	if _, exists := insertionsAt[primaryDocumentID]; !exists {
		insertionsAt[primaryDocumentID] = make(map[int]int)
	}
	insertionsAt[primaryDocumentID][maxEnd] = insertions

	return insertions
}

//   Performs a majority vote across all parts of the alignment
//   If indices were counted as aligning, they are used in the vote
//  The relationship between alignments in a cluster is such that
//  the primary alignment region is very similar in both
//  Also eturns an integer representing the number of corrections made
func MajorityVote(primaryDocumentID string, alignmentMaps []alignMap, documents []common.Document, docMap map[string]int) int {

	if _, exists := removeIndices[primaryDocumentID]; !exists {
		removeIndices[primaryDocumentID] = make(map[int]bool)
	}

	if _, exists := additionIndices[primaryDocumentID]; !exists {
		additionIndices[primaryDocumentID] = make(map[int][]rune)
	}

	if _, exists := editIndices[primaryDocumentID]; !exists {
		editIndices[primaryDocumentID] = make(map[int]rune)
	}

	if _, exists := mVoteLogs[primaryDocumentID]; !exists {
		mVoteLogs[primaryDocumentID] = make(map[int]common.Vote)
		newVoteLogs[primaryDocumentID] = make(map[int]common.Vote)
	}

	threshold := -math.Log2(flags.LmThreshold)

	noDeletions := 0
	noInsertions := 0
	if flags.HandleInsertionDeletion {
		noDeletions = applyDeletions(primaryDocumentID, alignmentMaps, documents, docMap)
		noInsertions = applyInsertions(primaryDocumentID, alignmentMaps, documents, docMap)
	}
	noCorrections := noDeletions + noInsertions
	maxEnd := 0
	minStart := 100000000
	primText := documents[docMap[primaryDocumentID]].Text

	for _, alMap := range alignmentMaps {
		if alMap.Start < minStart {
			minStart = alMap.Start
		}
		if alMap.End > maxEnd {
			maxEnd = alMap.End
		}
	}

	if flags.UseLM {
		words = wordsBeforePoint(primText, minStart, n)
	}
	var currentWord string
	if len(words) > 0 {
		currentWord = words[len(words)-1]
	}
	requiresNewWord := false

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

			}
		}
		for r := range counts {
			if counts[r] > max {
				max = counts[r]
				maxRune = r
			}
		}
		prevCorrections := noCorrections
		//fmt.Println(counts)
		//fmt.Println(primText[ind])
		if primText[ind] != maxRune && max > numVotes/2 && numVotes > 2 {
			if flags.UseLM  && !isNewLineSpace(primText[ind]) {
				if len(words) > 0 {
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
							editIndices[primaryDocumentID][ind] = maxRune
							noCorrections += 1
						} else {
							prevCount += 1
						}
					} else {
						editIndices[primaryDocumentID][ind] = maxRune
						noCorrections += 1
					}
				} else {
						editIndices[primaryDocumentID][ind] = maxRune
						noCorrections += 1
				}
			} else if !isNewLineSpace(primText[ind]) {
				editIndices[primaryDocumentID][ind] = maxRune
				noCorrections += 1
			}
		}
		if prevCorrections < noCorrections {
			if _, exists := mVoteLogs[primaryDocumentID][ind]; !exists {
				mVoteLogs[primaryDocumentID][ind] = common.Vote{
					EditDict:   map[string]int{},
					InsertDict: map[string]int{},
				}
			}
			for key, val := range counts {
				mVoteLogs[primaryDocumentID][ind].EditDict[string([]rune{key})] = val
			}
		}
	}
	//fmt.Println(string(primText))
	if flags.Logging && noCorrections > 0 {
		reuseCluster := make(map[string]string)
		p := []rune(strings.Repeat("_", maxEnd+1-minStart))
		for _, m := range alignmentMaps {
			for i := m.Start; i <= m.End; i++ {
				if _, exists := m.Mapping[i]; exists {
					p[i-minStart] = primText[i]
				}
			}
		}
		reuseCluster[primaryDocumentID] = string(p)
		for _, m := range alignmentMaps {
			s := m.Mapping[m.Start]
			e := m.Mapping[m.End]
			r := strings.Repeat("_", m.Start-minStart)
			t := []rune(strings.Repeat("_", e+1-s))
			for secPos := s; secPos <= e; secPos++ {
				t[secPos-s] = documents[docMap[m.SecondaryDocumentID]].Text[secPos]
			}
			reuseCluster[m.SecondaryDocumentID] = r + string(t)

		}
		reuseGraph[primaryDocumentID] = append(reuseGraph[primaryDocumentID], reuseCluster)
		oldStartEndGraph[primaryDocumentID] = append(oldStartEndGraph[primaryDocumentID], map[string]int{"start": minStart, "end": maxEnd})
		reuseStartEndGraph[primaryDocumentID] = append(reuseStartEndGraph[primaryDocumentID],
			map[string]int{"start": minStart - getDifferenceSoFar(primaryDocumentID, minStart),
				"end": maxEnd - getDifferenceSoFar(primaryDocumentID, maxEnd)})
	}
	return noCorrections
}

func wordsBeforePoint(text []rune, pos int, n int) []string {
	words := make([]string, 0)
	wordStarts := make([]int, 0)
	hitChars := false
	for i := pos; i > -1; i-- {
		if unicode.IsSpace(text[i]) {
			if hitChars {
				wordStarts = append(wordStarts, i+1)
				hitChars = false
				if len(wordStarts) > n {
					break
				}
			}
		} else if !hitChars {
			hitChars = true
		}
	}

	for i := len(wordStarts) - 1; i > -1; i-- {
		words = append(words, getCurrentWord(text, wordStarts[i]))
	}
	return words
}

func getCurrentWord(text []rune, pos int) string {

	end := pos
	for i := pos; i < len(text); i++ {
		if unicode.IsSpace(text[i]) {
			break
		}
		end++
	}

	if end == len(text) {
		return string(text[pos:])
	}

	return string(text[pos : end+1])
}

func getLmScore(word string, context string) string {
	requestBody, _ := json.Marshal(map[string]string{
		"word":     word,
		"sentence": context,
	})

	resp, err := http.Post("http://localhost:5000/", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "inf"
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "inf"
	}
	s := string(body)
	return s
}
