package diffpatch_test

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/manterfield/diffpatch"
)

// Test case represents a test scenario
type testCase struct {
	name    string
	oldFile string
	newFile string
	splitBy string
}

// Load test file pairs automatically
func loadTestCases() ([]testCase, error) {
	var cases []testCase

	// Find all test pairs
	pattern := regexp.MustCompile(`^test(\d+)_old\.md$`)

	err := filepath.WalkDir("test_documents", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		matches := pattern.FindStringSubmatch(d.Name())
		if len(matches) == 2 {
			testNum := matches[1]
			oldFile := path
			newFile := filepath.Join("test_documents", fmt.Sprintf("test%s_new.md", testNum))

			// Check if corresponding new file exists
			if _, err := os.Stat(newFile); err == nil {
				// Add test cases for each split type
				cases = append(cases, testCase{
					name:    fmt.Sprintf("test%s_lines", testNum),
					oldFile: oldFile,
					newFile: newFile,
					splitBy: "\n",
				})
				cases = append(cases, testCase{
					name:    fmt.Sprintf("test%s_sentences", testNum),
					oldFile: oldFile,
					newFile: newFile,
					splitBy: ".",
				})
				cases = append(cases, testCase{
					name:    fmt.Sprintf("test%s_words", testNum),
					oldFile: oldFile,
					newFile: newFile,
					splitBy: " ",
				})
			}
		}
		return nil
	})

	return cases, err
}

// Load file and split by delimiter
func loadFileSplit(filename, delimiter string) ([]string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(string(content), delimiter)
	// Filter out empty parts
	var result []string
	for _, part := range parts {
		if strings.TrimSpace(part) != "" {
			result = append(result, part)
		}
	}
	return result, nil
}

func TestDiffPatch(t *testing.T) {
	cases, err := loadTestCases()
	if err != nil {
		t.Fatalf("Failed to load test cases: %v", err)
	}

	if len(cases) == 0 {
		t.Fatal("No test cases found")
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Load old and new arrays
			oldArr, err := loadFileSplit(tc.oldFile, tc.splitBy)
			if err != nil {
				t.Fatalf("Failed to load old file %s: %v", tc.oldFile, err)
			}

			newArr, err := loadFileSplit(tc.newFile, tc.splitBy)
			if err != nil {
				t.Fatalf("Failed to load new file %s: %v", tc.newFile, err)
			} // Generate patch
			patch := diffpatch.Diff(oldArr, newArr)

			// Apply patch
			result := diffpatch.ApplyPatch(oldArr, patch)

			// Verify result equals newArr
			if !reflect.DeepEqual(result, newArr) {
				t.Errorf("Patch application failed for %s", tc.name)
				t.Logf("Expected: %v", newArr)
				t.Logf("Got: %v", result)
			}
		})
	}
}

func BenchmarkDiffPatch(b *testing.B) {
	cases, err := loadTestCases()
	if err != nil {
		b.Fatalf("Failed to load test cases: %v", err)
	}

	for _, tc := range cases {
		// Load test data once
		oldArr, err := loadFileSplit(tc.oldFile, tc.splitBy)
		if err != nil {
			b.Fatalf("Failed to load old file %s: %v", tc.oldFile, err)
		}

		newArr, err := loadFileSplit(tc.newFile, tc.splitBy)
		if err != nil {
			b.Fatalf("Failed to load new file %s: %v", tc.newFile, err)
		}

		b.Run(tc.name+"_diff", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = diffpatch.Diff(oldArr, newArr)
			}
		})

		b.Run(tc.name+"_apply", func(b *testing.B) {
			patch := diffpatch.Diff(oldArr, newArr)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = diffpatch.ApplyPatch(oldArr, patch)
			}
		})

		b.Run(tc.name+"_combined", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				patch := diffpatch.Diff(oldArr, newArr)
				_ = diffpatch.ApplyPatch(oldArr, patch)
			}
		})
	}
}
