package diffpatch

import (
	"encoding/json"
)

// Operation represents a diff operation as [index, deleteCount, additions]
type Operation []interface{}

// Diff computes the difference between oldArr and newArr
func Diff(oldArr, newArr []string) []Operation {
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

	ops := make([]Operation, 0, 4) // Pre-allocate with capacity to reduce allocations
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

// ApplyPatch applies the diff operations to the original array
func ApplyPatch(arr []string, ops []Operation) []string {
	opsLen := len(ops)
	if opsLen == 0 {
		result := make([]string, len(arr))
		copy(result, arr)
		return result
	}

	// Single operation optimization
	if opsLen == 1 {
		op := ops[0]
		opIndex := op[0].(int)
		deleteCount := op[1].(int)
		additions := op[2].([]string)

		arrLen := len(arr)
		result := make([]string, arrLen-deleteCount+len(additions))

		pos := 0
		// Copy prefix
		for i := 0; i < opIndex; i++ {
			result[pos] = arr[i]
			pos++
		}
		// Add new elements
		for _, add := range additions {
			result[pos] = add
			pos++
		}
		// Copy suffix
		suffixStart := opIndex + deleteCount
		for i := suffixStart; i < arrLen; i++ {
			result[pos] = arr[i]
			pos++
		}

		return result
	}

	// Multiple operations: copy and apply in reverse
	result := make([]string, len(arr))
	copy(result, arr)

	// Apply operations in reverse order
	for i := opsLen - 1; i >= 0; i-- {
		op := ops[i]
		opIndex := op[0].(int)
		deleteCount := op[1].(int)
		additions := op[2].([]string)
		addLen := len(additions)

		if addLen == 0 {
			// Delete only
			result = append(result[:opIndex], result[opIndex+deleteCount:]...)
		} else if deleteCount == 0 {
			// Insert only
			temp := make([]string, len(result)+addLen)
			copy(temp[:opIndex], result[:opIndex])
			copy(temp[opIndex:opIndex+addLen], additions)
			copy(temp[opIndex+addLen:], result[opIndex:])
			result = temp
		} else {
			// Replace
			temp := make([]string, len(result)-deleteCount+addLen)
			copy(temp[:opIndex], result[:opIndex])
			copy(temp[opIndex:opIndex+addLen], additions)
			copy(temp[opIndex+addLen:], result[opIndex+deleteCount:])
			result = temp
		}
	}

	return result
}

// Helper function to convert Go operations to JSON format compatible with JS
func OperationsToJSON(ops []Operation) (string, error) {
	jsonOps := make([][]interface{}, len(ops))
	for i, op := range ops {
		jsonOps[i] = []interface{}{op[0], op[1], op[2]}
	}
	jsonBytes, err := json.Marshal(jsonOps)
	return string(jsonBytes), err
}

// Helper function to parse JSON operations from JS
func OperationsFromJSON(jsonStr string) ([]Operation, error) {
	var jsonOps [][]interface{}
	err := json.Unmarshal([]byte(jsonStr), &jsonOps)
	if err != nil {
		return nil, err
	}

	ops := make([]Operation, len(jsonOps))
	for i, jsonOp := range jsonOps {
		// Convert additions from []interface{} to []string
		additionsInterface := jsonOp[2].([]interface{})
		additions := make([]string, len(additionsInterface))
		for j, add := range additionsInterface {
			additions[j] = add.(string)
		}
		ops[i] = Operation{int(jsonOp[0].(float64)), int(jsonOp[1].(float64)), additions}
	}
	return ops, nil
}
