package alignment

import (
	"math"
)

func Gotoh(matchReward int, gapOpen int, gapEx int, a []rune, b []rune) (int, []int, []int) {
	lenA := len(a)
	lenB := len(b)
	gapOpen = -gapOpen
	gapEx = -gapEx
	d := make([][]int, lenA+1)
	p := make([][]int, lenA+1)
	q := make([][]int, lenA+1)

	d_bc := make([][]int, lenA+1)
	p_bc := make([][]int, lenA+1)
	q_bc := make([][]int, lenA+1)
	for i := 0; i < lenA+1; i++ {
		d[i] = make([]int, lenB+1)
		p[i] = make([]int, lenB+1)
		q[i] = make([]int, lenB+1)
		d_bc[i] = make([]int, lenB+1)
		p_bc[i] = make([]int, lenB+1)
		q_bc[i] = make([]int, lenB+1)
	}

	for j := 1; j < lenB+1; j++ {
		q[0][j] = math.MinInt32
		p[0][j] = math.MinInt32
		d_bc[0][j] = 100
		q_bc[0][j] = 1
	}

	for i := 1; i < lenA+1; i++ {
		p[i][0] = math.MinInt32
		q[i][0] = math.MinInt32
		d_bc[i][0] = 10
		p_bc[i][0] = 1
	}
	var matchScore int
	maxScore := 0
	maxI := 0
	maxJ := 0
	for i := 1; i < lenA+1; i++ {
		for j := 1; j < lenB+1; j++ {
			p[i][j] = Max(d[i-1][j]+gapOpen, p[i-1][j]+gapEx)
			q[i][j] = Max(d[i][j-1]+gapOpen, q[i][j-1]+gapEx)
			if a[i-1] == b[j-1] {
				matchScore = matchReward
			} else {
				matchScore = -matchReward
			}
			d[i][j] = Max(0, Max(d[i-1][j-1]+matchScore, Max(p[i][j], q[i][j])))

			if d[i-1][j]+gapOpen < p[i-1][j]+gapEx {
				p_bc[i][j] = 1
			}

			if d[i][j-1]+gapOpen < q[i][j-1]+gapEx {
				q_bc[i][j] = 1
			}

			if d[i][j] == d[i-1][j-1]+matchScore {
				d_bc[i][j] = 1
			}
			if d[i][j] == p[i][j] {
				d_bc[i][j] += 10
			}
			if d[i][j] == q[i][j] {
				d_bc[i][j] += 100
			}

			if d[i][j] > maxScore {
				maxI = i
				maxJ = j
				maxScore = d[i][j]
			}
		}

	}
	i, j := maxI, maxJ
	level := 0
	aIndices, bIndices := make([]int, 0), make([]int, 0)
	for d[i][j] > 0 {
		if level == 0 {
			if (d_bc[i][j]/10)%10 == 1 {
				level = 1
			} else if (d_bc[i][j]/100)%10 == 1 {
				level = 2
			} else if d_bc[i][j]%10 == 1 {
				aIndices = append(aIndices, i-1)
				bIndices = append(bIndices, j-1)
				j -= 1
				i -= 1
			}
		} else if level == 1 {
			if p_bc[i][j] == 0 {
				level = 0
			}
			i -= 1
		} else {
			if q_bc[i][j] == 0 {
				level = 0
			}
			j -= 1
		}
	}
	d = nil
	p = nil
	q = nil

	d_bc = nil
	p_bc = nil
	q_bc = nil
	return maxScore, reverseInt(aIndices), reverseInt(bIndices)
}
