package testutil

import (
	"os"
	"path/filepath"
	"testing"
)

type Tester interface {
	Helper()
	Errorf(format string, args ...interface{})
	Error(args ...interface{})
}

func TempFile(t *testing.T, content string) string {
	t.Helper()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "testfile.txt")

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	return filePath
}

func TempDirWithFiles(t *testing.T, files map[string]string) string {
	t.Helper()

	tmpDir := t.TempDir()

	for name, content := range files {
		filePath := filepath.Join(tmpDir, name)
		dir := filepath.Dir(filePath)

		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("failed to create dir: %v", err)
		}

		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatalf("failed to create file: %v", err)
		}
	}

	return tmpDir
}

func AssertEqual(t Tester, expected, actual interface{}) {
	t.Helper()

	if expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func AssertContains(t Tester, s, substr string) {
	t.Helper()

	if !contains(s, substr) {
		t.Errorf("expected '%s' to contain '%s'", s, substr)
	}
}

func AssertError(t Tester, err error) {
	t.Helper()

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func AssertNoError(t Tester, err error) {
	t.Helper()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func contains(s, substr string) bool {
	if len(substr) > len(s) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
