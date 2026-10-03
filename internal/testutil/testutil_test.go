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

	fake := &captureTester{}
	AssertEqual(fake, 1, 2)
	AssertContains(fake, "alpha", "zzz")
	AssertError(fake, nil)
	AssertNoError(fake, errors.New("bad"))
	if len(fake.errors) != 4 {
		t.Fatalf("recorded %d failures", len(fake.errors))
	}
}

type captureTester struct {
	errors []string
}

func (c *captureTester) Helper() {}
func (c *captureTester) Errorf(format string, args ...interface{}) {
	c.errors = append(c.errors, format)
}
func (c *captureTester) Error(args ...interface{}) {
	c.errors = append(c.errors, "error")
}
