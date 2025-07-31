// Reusable empty array to avoid allocations
const EMPTY_ARRAY = [];

function diff(oldArr, newArr) {
  const oldLen = oldArr.length;
  const newLen = newArr.length;

  if (oldLen === 0) {
    if (newLen === 0) {
      return [];
    }
    return [{ i: [0, 0], a: newArr.slice() }];
  }

  if (newLen === 0) {
    return [{ i: [0, oldLen], a: EMPTY_ARRAY }];
  }

  const ops = []; // Pre-allocate with capacity to reduce allocations
  let oldPos = 0;
  let newPos = 0;

  while (oldPos < oldLen || newPos < newLen) {
    // Skip matching prefix
    while (
      oldPos < oldLen &&
      newPos < newLen &&
      oldArr[oldPos] === newArr[newPos]
    ) {
      oldPos++;
      newPos++;
    }

    if (oldPos >= oldLen && newPos >= newLen) {
      break;
    }

    // We found a difference, now find the bounds of this change
    const changeOldStart = oldPos;
    const changeNewStart = newPos;

    // Find the next synchronization point using a forward scan
    let syncFound = false;
    let syncOldPos = oldLen;
    let syncNewPos = newLen;

    // Look for the next matching sequence
    const maxScanDistance = 50; // Prevent excessive scanning

    for (
      let scanDist = 1;
      scanDist <= maxScanDistance && !syncFound;
      scanDist++
    ) {
      // Try matching elements at various distances
      const maxOldOffset = Math.min(scanDist, oldLen - changeOldStart - 1);
      for (let oldOffset = 0; oldOffset <= maxOldOffset; oldOffset++) {
        const newOffset = scanDist - oldOffset;
        if (changeNewStart + newOffset >= newLen) {
          continue;
        }

        const oldIdx = changeOldStart + oldOffset;
        const newIdx = changeNewStart + newOffset;

        if (
          oldArr[oldIdx] === newArr[newIdx] &&
          isGoodSyncPoint(oldArr, newArr, oldIdx, newIdx)
        ) {
          syncOldPos = oldIdx;
          syncNewPos = newIdx;
          syncFound = true;
          break;
        }
      }
    }

    // Create operation for this change
    const deleteCount = syncOldPos - changeOldStart;
    const addCount = syncNewPos - changeNewStart;

    if (deleteCount > 0 || addCount > 0) {
      let additions;
      if (addCount === 0) {
        additions = EMPTY_ARRAY;
      } else if (addCount === 1) {
        additions = [newArr[changeNewStart]];
      } else {
        additions = newArr.slice(changeNewStart, changeNewStart + addCount);
      }

      ops.push({ i: [changeOldStart, deleteCount], a: additions });
    }

    oldPos = syncOldPos;
    newPos = syncNewPos;
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
    const result = new Array(arrLen - op.i[1] + op.a.length);

    let pos = 0;
    const opIndex = op.i[0];
    const addLen = op.a.length;

    // Copy prefix - avoid indirection
    for (let i = 0; i < opIndex; i++) {
      result[pos++] = arr[i];
    }
    // Add new elements - avoid property lookup
    for (let i = 0; i < addLen; i++) {
      result[pos++] = op.a[i];
    }
    // Copy suffix - cache calculation
    const suffixStart = opIndex + op.i[1];
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
    const addLen = op.a.length;

    if (addLen === 0) {
      result.splice(op.i[0], op.i[1]);
    } else if (op.i[1] === 0) {
      if (addLen === 1) {
        result.splice(op.i[0], 0, op.a[0]);
      } else {
        result.splice(op.i[0], 0, ...op.a);
      }
    } else {
      if (addLen === 1) {
        result.splice(op.i[0], op.i[1], op.a[0]);
      } else {
        result.splice(op.i[0], op.i[1], ...op.a);
      }
    }
  }

  return result;
}

// isGoodSyncPoint checks if two positions represent a good synchronization point
function isGoodSyncPoint(oldArr, newArr, oldIdx, newIdx) {
  const oldLen = oldArr.length;
  const newLen = newArr.length;

  if (oldIdx >= oldLen || newIdx >= newLen) {
    return oldIdx >= oldLen && newIdx >= newLen;
  }

  // Check for at least 2 consecutive matches or end of arrays
  let matches = 0;
  const maxCheck = 2;
  const maxOld = oldLen - oldIdx;
  const maxNew = newLen - newIdx;
  const checkLimit =
    maxCheck < maxOld
      ? maxCheck < maxNew
        ? maxCheck
        : maxNew
      : maxOld < maxNew
      ? maxOld
      : maxNew;

  for (let i = 0; i < checkLimit; i++) {
    if (oldArr[oldIdx + i] === newArr[newIdx + i]) {
      matches++;
    } else {
      break;
    }
  }

  // Good sync point if we have 2+ matches or we've reached the end of both arrays
  return (
    matches >= 2 || (oldIdx + matches >= oldLen && newIdx + matches >= newLen)
  );
}

// Export functions for interoperability testing
export { diff, applyPatch };
