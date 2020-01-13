package alignment

/**
* Given a list of indices, that may span multiple components,
* generates a list of per-component indices instead
* We can then use this in our alignment representation
**/

func PerComponentIndices(indices []int, componentLengths []int, startIndex int, endIndex int,
	startCompIndex int) ([][]int, int, int) {
	newIndices := make([][]int, 0)

	lastEnd := 0

	currentIndices := make([]int, 0)
	c := 0
	cLen := componentLengths[c]

	firstAffected := -1
	lastAffected := 0
	for i, indice := range indices {
		// Now on the next component
		for indice >= (lastEnd + cLen) {
			c += 1
			if i == 0 {
				firstAffected = c
			}
			lastAffected = c
			if len(currentIndices) > 0 {
				newIndices = append(newIndices, currentIndices)
				currentIndices = []int{}
			}
			lastEnd += cLen
			cLen = componentLengths[c]
		}
		currentIndices = append(currentIndices, indice-lastEnd)
	}
	if firstAffected == -1 {
		firstAffected = 0
	}
	newIndices = append(newIndices, currentIndices)
	return newIndices, startCompIndex + firstAffected, startCompIndex + lastAffected // TODO: Fix broken login with first ans last affected
}
