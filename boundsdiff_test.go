package diffpatch_test

import (
	"reflect"
	"testing"

	"github.com/manterfield/diffpatch"
)

func TestBoundsDiff(t *testing.T) {
	tests := []struct {
		name   string
		oldArr []string
		newArr []string
	}{
		{
			name:   "empty arrays",
			oldArr: []string{},
			newArr: []string{},
		},
		{
			name:   "insert into empty",
			oldArr: []string{},
			newArr: []string{"a", "b"},
		},
		{
			name:   "delete all",
			oldArr: []string{"a", "b"},
			newArr: []string{},
		},
		{
			name:   "example from requirements",
			oldArr: []string{"a", "b", "c", "d", "e"},
			newArr: []string{"a", "x", "c", "y", "e"},
		},
		{
			name:   "no changes",
			oldArr: []string{"a", "b", "c"},
			newArr: []string{"a", "b", "c"},
		},
		{
			name:   "simple replacement",
			oldArr: []string{"a", "b", "c"},
			newArr: []string{"a", "x", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := diffpatch.BoundsDiff(tt.oldArr, tt.newArr)

			// Verify the patch can be applied correctly
			result := diffpatch.ApplyPatch(tt.oldArr, got)
			if !reflect.DeepEqual(result, tt.newArr) {
				t.Errorf("ApplyPatch() = %v, want %v", result, tt.newArr)
			}
		})
	}
}
