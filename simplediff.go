package diffpatch

// SimpleDiff implements a fast O(n) linear scan algorithm
// Inspired by streaming algorithms and one-pass processing
func SimpleDiff(oldArr, newArr []string) []Operation {
	if len(oldArr) == 0 {
		if len(newArr) == 0 {
			return []Operation{}
		}
		return []Operation{{0, 0, newArr}}
	}

	if len(newArr) == 0 {
		return []Operation{{0, len(oldArr), []string{}}}
	}

	// Fast path: check if arrays are identical
	if arraysEqual(oldArr, newArr) {
		return []Operation{}
	}

	// Single pass algorithm - find common prefix and suffix
	prefixLen := 0
	for prefixLen < len(oldArr) && prefixLen < len(newArr) && oldArr[prefixLen] == newArr[prefixLen] {
		prefixLen++
	}

	// Find common suffix
	suffixLen := 0
	oldEnd := len(oldArr) - 1
	newEnd := len(newArr) - 1

	for suffixLen < (len(oldArr)-prefixLen) && suffixLen < (len(newArr)-prefixLen) &&
		oldArr[oldEnd-suffixLen] == newArr[newEnd-suffixLen] {
		suffixLen++
	}

	// Calculate the middle section that needs to be replaced
	oldStart := prefixLen
	oldLen := len(oldArr) - prefixLen - suffixLen
	newMiddle := newArr[prefixLen : len(newArr)-suffixLen]

	if oldLen == 0 && len(newMiddle) == 0 {
		return []Operation{} // No changes needed
	}

	return []Operation{{oldStart, oldLen, newMiddle}}
}

// arraysEqual checks if two string arrays are identical
func arraysEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
