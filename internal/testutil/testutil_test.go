package testutil

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestTempHelpersAndAssertions(t *testing.T) {
	file := TempFile(t, "alpha")
	body, err := os.ReadFile(file)
	if err != nil || string(body) != "alpha" {
		t.Fatalf("temp file %s %v", body, err)
	}

	dir := TempDirWithFiles(t, map[string]string{"nested/b.txt": "beta"})
	if _, err := os.Stat(filepath.Join(dir, "nested", "b.txt")); err != nil {
		t.Fatal(err)
	}

	AssertEqual(t, 1, 1)
	AssertContains(t, "hello", "ell")
	AssertNoError(t, nil)
	AssertError(t, errors.New("x"))
	if contains("ab", "abc") || contains("hello", "zzz") {
		t.Fatal("contains false cases")
	}
}
