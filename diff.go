package diffpatch

import (
	"encoding/json"
)

// Operation represents a diff operation as [index, deleteCount, additions]
type Operation []interface{}

// Diff computes the difference between oldArr and newArr
func Diff(oldArr, newArr []string) []Operation {
	var ops []Operation
	oldPos := 0
	newPos := 0
	oldLen := len(oldArr)
	newLen := len(newArr)

	// Compute common prefix
	prefix := 0
	for prefix < oldLen && prefix < newLen && oldArr[prefix] == newArr[prefix] {
		prefix++
	}

	// Compute common suffix
	suffix := 0
	for oldLen-suffix-1 >= prefix && newLen-suffix-1 >= prefix &&
		oldArr[oldLen-suffix-1] == newArr[newLen-suffix-1] {
		suffix++
	}

	// If we have trimmed prefix/suffix, diff the middle segments
	if prefix > 0 || suffix > 0 {
		midOldLen := oldLen - prefix - suffix
		midNewLen := newLen - prefix - suffix

		if midOldLen == 0 && midNewLen == 0 {
			return []Operation{}
		}

		if midOldLen == 0 {
			additions := newArr[prefix : newLen-suffix]
			return []Operation{{prefix, 0, additions}}
		}

		if midNewLen == 0 {
			return []Operation{{prefix, midOldLen, []string{}}}
		}

		oldMid := oldArr[prefix : oldLen-suffix]
		newMid := newArr[prefix : newLen-suffix]
		midOps := Diff(oldMid, newMid)

		// Adjust indices for the middle operations
		for i := range midOps {
			midOps[i][0] = midOps[i][0].(int) + prefix
		}
		return midOps
	}

	// Build index of newArr positions for fast lookups
	newIndex := make(map[string][]int)
	for i, val := range newArr {
		newIndex[val] = append(newIndex[val], i)
	}

	for oldPos < oldLen || newPos < newLen {
		// Skip matching elements
		for oldPos < oldLen && newPos < newLen && oldArr[oldPos] == newArr[newPos] {
			oldPos++
			newPos++
		}

		if oldPos >= oldLen && newPos >= newLen {
			break
		}

		opStart := oldPos
		addStart := newPos
		syncOld := oldLen
		syncNew := newLen

		// Optimized lookahead
		remainingOld := oldLen - oldPos
		remainingNew := newLen - newPos
		lookAhead := remainingOld
		if remainingNew > remainingOld {
			lookAhead = remainingNew
		}
		if lookAhead > 100 {
			lookAhead = 100
		}

		// Find sync point using indexed positions for performance
		outerFound := false
		searchEnd := newPos + lookAhead
		if searchEnd > newLen {
			searchEnd = newLen
		}

		for ahead := 0; ahead <= lookAhead && oldPos+ahead < oldLen && !outerFound; ahead++ {
			oldIdx := oldPos + ahead
			element := oldArr[oldIdx]
			positions, exists := newIndex[element]
			if !exists {
				continue
			}

			for _, pos := range positions {
				if pos < newPos || pos > searchEnd {
					continue
				}

				// Count matching sequence
				matchLen := 0
				maxOld := oldLen - oldIdx
				maxNew := newLen - pos
				maxLen := maxOld
				if maxNew < maxOld {
					maxLen = maxNew
				}

				for matchLen < maxLen && oldArr[oldIdx+matchLen] == newArr[pos+matchLen] {
					matchLen++
				}

				// Valid match if significant or end of both arrays
				if matchLen >= 2 || (oldIdx+matchLen >= oldLen && pos+matchLen >= newLen) {
					syncOld = oldIdx
					syncNew = pos
					outerFound = true
					break
				}
			}
		}

		// Create operation
		deleteCount := syncOld - opStart
		addCount := syncNew - addStart

		if deleteCount > 0 || addCount > 0 {
			var additions []string
			if addCount == 0 {
				additions = []string{}
			} else {
				additions = make([]string, addCount)
				copy(additions, newArr[addStart:addStart+addCount])
			}

			ops = append(ops, Operation{opStart, deleteCount, additions})
		}

		oldPos = syncOld
		newPos = syncNew
	}

	return ops
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
