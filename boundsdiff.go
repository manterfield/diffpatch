package diffpatch

// BoundsDiff implements an efficient bounds-based diffing algorithm
// The algorithm focuses on finding the 'bounds' of each change by reading forward
// until it finds a difference, then using a scanning approach to find sync points
func BoundsDiff(oldArr, newArr []string) []Operation {
	oldLen := len(oldArr)
	newLen := len(newArr)

	if oldLen == 0 {
		if newLen == 0 {
			return []Operation{}
		}
		return []Operation{{0, 0, newArr}}
	}

	if newLen == 0 {
		return []Operation{{0, oldLen, []string{}}}
	}

	var ops []Operation
	oldPos := 0
	newPos := 0

	for oldPos < oldLen || newPos < newLen {
		// Find start of change by scanning forward until difference found
		changeOldStart := oldPos
		changeNewStart := newPos

		// Skip matching prefix
		for oldPos < oldLen && newPos < newLen && oldArr[oldPos] == newArr[newPos] {
			oldPos++
			newPos++
		}

		if oldPos >= oldLen && newPos >= newLen {
			break
		}

		// We found a difference, now find the bounds of this change
		changeOldStart = oldPos
		changeNewStart = newPos

		// Find the next synchronization point using a forward scan
		syncFound := false
		syncOldPos := oldLen
		syncNewPos := newLen

		// Look for the next matching sequence
		maxScanDistance := 50 // Prevent excessive scanning

		for scanDist := 1; scanDist <= maxScanDistance && !syncFound; scanDist++ {
			// Try matching elements at various distances
			for oldOffset := 0; oldOffset <= scanDist && changeOldStart+oldOffset < oldLen; oldOffset++ {
				newOffset := scanDist - oldOffset
				if changeNewStart+newOffset >= newLen {
					continue
				}

				oldIdx := changeOldStart + oldOffset
				newIdx := changeNewStart + newOffset

				if oldArr[oldIdx] == newArr[newIdx] && isGoodSyncPoint(oldArr, newArr, oldIdx, newIdx) {
					syncOldPos = oldIdx
					syncNewPos = newIdx
					syncFound = true
					break
				}
			}
		}

		// Create operation for this change
		deleteCount := syncOldPos - changeOldStart
		addCount := syncNewPos - changeNewStart

		if deleteCount > 0 || addCount > 0 {
			additions := make([]string, addCount)
			if addCount > 0 {
				copy(additions, newArr[changeNewStart:changeNewStart+addCount])
			}

			ops = append(ops, Operation{changeOldStart, deleteCount, additions})
		}

		oldPos = syncOldPos
		newPos = syncNewPos
	}

	return ops
}

// isGoodSyncPoint checks if two positions represent a good synchronization point
func isGoodSyncPoint(oldArr, newArr []string, oldIdx, newIdx int) bool {
	if oldIdx >= len(oldArr) || newIdx >= len(newArr) {
		return oldIdx >= len(oldArr) && newIdx >= len(newArr)
	}

	// Check for at least 2 consecutive matches or end of arrays
	matches := 0
	maxCheck := 2

	for i := 0; i < maxCheck && oldIdx+i < len(oldArr) && newIdx+i < len(newArr); i++ {
		if oldArr[oldIdx+i] == newArr[newIdx+i] {
			matches++
		} else {
			break
		}
	}

	// Good sync point if we have 2+ matches or we've reached the end of both arrays
	return matches >= 2 || (oldIdx+matches >= len(oldArr) && newIdx+matches >= len(newArr))
}
