// Reusable empty array to avoid allocations
const EMPTY_ARRAY = [];

function diff(oldArr, newArr) {
  const ops = [];
  let oldPos = 0;
  let newPos = 0;
  const oldLen = oldArr.length;
  const newLen = newArr.length;

  // Compute common prefix
  let prefix = 0;
  while (
    prefix < oldLen &&
    prefix < newLen &&
    oldArr[prefix] === newArr[prefix]
  ) {
    prefix++;
  }
  // Compute common suffix
  let suffix = 0;
  while (
    oldLen - suffix - 1 >= prefix &&
    newLen - suffix - 1 >= prefix &&
    oldArr[oldLen - suffix - 1] === newArr[newLen - suffix - 1]
  ) {
    suffix++;
  }
  // If we have trimmed prefix/suffix, diff the middle segments
  if (prefix > 0 || suffix > 0) {
    const midOldLen = oldLen - prefix - suffix;
    const midNewLen = newLen - prefix - suffix;
    if (midOldLen === 0 && midNewLen === 0) {
      return [];
    }
    if (midOldLen === 0) {
      return [[prefix, 0, newArr.slice(prefix, newLen - suffix)]];
    }
    if (midNewLen === 0) {
      return [[prefix, midOldLen, EMPTY_ARRAY]];
    }
    const oldMid = oldArr.slice(prefix, oldLen - suffix);
    const newMid = newArr.slice(prefix, newLen - suffix);
    const midOps = diff(oldMid, newMid);
    return midOps.map((op) => [op[0] + prefix, op[1], op[2]]);
  }
  // Build index of newArr positions for fast lookups
  const newIndex = new Map();
  for (let i = 0; i < newLen; i++) {
    const val = newArr[i];
    const arr = newIndex.get(val);
    if (arr) arr.push(i);
    else newIndex.set(val, [i]);
  }

  while (oldPos < oldLen || newPos < newLen) {
    // Skip matching elements
    while (
      oldPos < oldLen &&
      newPos < newLen &&
      oldArr[oldPos] === newArr[newPos]
    ) {
      oldPos++;
      newPos++;
    }

    if (oldPos >= oldLen && newPos >= newLen) break;

    const opStart = oldPos;
    const addStart = newPos;
    let syncOld = oldLen;
    let syncNew = newLen;

    // Optimized lookahead - avoid Math.max when possible
    const remainingOld = oldLen - oldPos;
    const remainingNew = newLen - newPos;
    const lookAhead =
      remainingOld > remainingNew
        ? remainingOld > 100
          ? 100
          : remainingOld
        : remainingNew > 100
        ? 100
        : remainingNew;

    // Find sync point using indexed positions for performance
    let outerFound = false;
    const searchEnd = newPos + lookAhead < newLen ? newPos + lookAhead : newLen;
    for (
      let ahead = 0;
      ahead <= lookAhead && oldPos + ahead < oldLen && !outerFound;
      ahead++
    ) {
      const oldIdx = oldPos + ahead;
      const element = oldArr[oldIdx];
      const positions = newIndex.get(element);
      if (!positions) continue;
      for (let k = 0; k < positions.length; k++) {
        const pos = positions[k];
        if (pos < newPos) continue;
        if (pos > searchEnd) break;

        // Count matching sequence
        let matchLen = 0;
        const maxOld = oldLen - oldIdx;
        const maxNew = newLen - pos;
        const maxLen = maxOld < maxNew ? maxOld : maxNew;
        while (
          matchLen < maxLen &&
          oldArr[oldIdx + matchLen] === newArr[pos + matchLen]
        ) {
          matchLen++;
        }

        // Valid match if significant or end of both arrays
        if (
          matchLen >= 2 ||
          (oldIdx + matchLen >= oldLen && pos + matchLen >= newLen)
        ) {
          syncOld = oldIdx;
          syncNew = pos;
          outerFound = true;
          break;
        }
      }
    }

    // Create operation
    const deleteCount = syncOld - opStart;
    const addCount = syncNew - addStart;

    if (deleteCount > 0 || addCount > 0) {
      let additions = EMPTY_ARRAY;
      if (addCount === 1) {
        additions = [newArr[addStart]];
      } else if (addCount > 1) {
        additions = new Array(addCount);
        for (let i = 0; i < addCount; i++) {
          additions[i] = newArr[addStart + i];
        }
      }

      ops[ops.length] = [opStart, deleteCount, additions];
    }

    oldPos = syncOld;
    newPos = syncNew;
  }

  return ops;
}

function applyPatch(arr, ops) {
  const opsLen = ops.length;
  if (opsLen === 0) return arr.slice();

  // Single operation optimization - specialized path
  if (opsLen === 1) {
    const op = ops[0];
    const arrLen = arr.length;
    const result = new Array(arrLen - op[1] + op[2].length);

    let pos = 0;
    const opIndex = op[0];
    const addLen = op[2].length;

    // Copy prefix - avoid indirection
    for (let i = 0; i < opIndex; i++) {
      result[pos++] = arr[i];
    }
    // Add new elements - avoid property lookup
    for (let i = 0; i < addLen; i++) {
      result[pos++] = op[2][i];
    }
    // Copy suffix - cache calculation
    const suffixStart = opIndex + op[1];
    for (let i = suffixStart; i < arrLen; i++) {
      result[pos++] = arr[i];
    }

    return result;
  }

  // Multiple operations: copy and apply in reverse - optimized
  const arrLen = arr.length;
  const result = new Array(arrLen);
  for (let i = 0; i < arrLen; i++) {
    result[i] = arr[i];
  }

  // Reverse iteration - avoid length property lookup
  for (let i = opsLen - 1; i >= 0; i--) {
    const op = ops[i];
    const addLen = op[2].length;

    if (addLen === 0) {
      result.splice(op[0], op[1]);
    } else if (op[1] === 0) {
      if (addLen === 1) {
        result.splice(op[0], 0, op[2][0]);
      } else {
        result.splice(op[0], 0, ...op[2]);
      }
    } else {
      if (addLen === 1) {
        result.splice(op[0], op[1], op[2][0]);
      } else {
        result.splice(op[0], op[1], ...op[2]);
      }
    }
  }

  return result;
}

// Export functions for interoperability testing
export { diff, applyPatch };
