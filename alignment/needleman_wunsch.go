package alignment


// The needleman wunsch local alignment algorithm,  simplified
func NeedlemanWunsch(matchReward int, gapCost int, a []rune, b []rune) (int, []int, []int) {
	l1 := len(a)
	l2 := len(b)
	cost := 0
	max_i := 0
	matched := false
	if l1 == 1 {
		for i, c := range b {
			if c == a[0] {
				max_i = i
				matched = true
			}
			cost -= gapCost
		}
		if matched == true {
			return matchReward + cost + gapCost, []int{0}, []int{max_i}
		}
		return -matchReward + cost + gapCost, []int{0}, []int{0}
	} else if l2 == 1 {
		for i, c := range a {
			if c == b[0] {
				max_i = i
				matched = true
			}
			cost -= gapCost
		}
		if matched == true {
			return matchReward + cost + gapCost, []int{max_i}, []int{0}
		}
		return -matchReward + cost + gapCost, []int{0}, []int{0}
	}
	return 0, []int{}, []int{}
}
