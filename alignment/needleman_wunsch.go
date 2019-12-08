package alignment

import (
	"math"
)

/**
* Updates the maximum value of the alignment found, and adds to array a accordingly
**/
func updateMax(i int, j int, d [][]float64, r [][]int) {

	m := math.Max(d[i-1][j-1], math.Max(d[i][j-1], d[i-1][j]))
	if d[i-1][j-1] == m {
			r[i][j] = 3
		}
	} else if d[i][j-1] == m {
		r[i][j] = 2
	} else {
		r[i][j] = 1
	}
}

/**
* The needleman wunsch local alignment algorithm
**/
func NeedlemanWunsch(matchReward float64, gapCost float64, a []rune, b []rune) (float64, []int, []int) {
	lenA := len(a)
	lenB := len(b)

	d := make([][]float64, lenA+1)
	r := make([][]int, lenA+1)

	for i := range d {
		d[i] = make([]float64, lenB+1)
		r[i] = make([]int, lenB+1)
	}

	for i := 1; i < lenA+1; i++ {
		d[i][0] = d[i-1][0] - gapCost
	}

	for j := 1; j < lenB+1; j++ {
		d[0][j] = d[0][j-1] - gapCost
	}

	for i := 1; i < lenA+1; i++ {
		for j := 1; j < lenB+1; j++ {
			var match float64
			if a[i-1] == b[j-1] {
				match = d[i-1][j-1] + matchReward
			} else {
				match = d[i-1][j-1] - matchReward
			}
			del := d[i-1][j] - gapCost
			ins := d[i][j-1] - gapCost
			d[i][j] = math.Max(match, math.Max(del, ins))
			updateMax(i, j, d, r)
		}
	}

	indicesA := []int{}
	indicesB := []int{}
	i := lenA
	j := lenB

	for i > 0 && j > 0 {
		if i > 0 && r[i][j] == 1 {
			i -= 1
		} else if j > 0 && r[i][j] == 2 {
			j -= 1
		} else {
			indicesA = append([]int{i - 1}, indicesA...)
			indicesB = append([]int{j - 1}, indicesB...)
			i -= 1
			j -= 1
		}
	}
	return d[lenA][lenB], indicesA, indicesB
}
