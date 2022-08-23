package main

import (
	"bytes"
	"os"
	"os/exec"
	"testing"
)

func copyFile(t *testing.T, src, dst string) {
	t.Helper()
	input, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("failed to read %s: %+v", src, err)
	}
	err = os.WriteFile(dst, input, 0777)
	if err != nil {
		t.Fatalf("failed to write %s: %+v", dst, err)
	}
}

func compareFiles(t *testing.T, a, b string) {
	contentA, err := os.ReadFile(a)
	if err != nil {
		t.Fatalf("failed to read file %s: %+v", a, err)
	}
	contentB, err := os.ReadFile(b)
	if err != nil {
		t.Fatalf("failed to read golden file %s: %+v", b, err)
	}
	if !bytes.Equal(contentA, contentB) {
		if _, err := exec.LookPath("git"); err == nil {
			diff := exec.Command("git", "diff", "--no-index", a, b)
			stdoutStderr, _ := diff.CombinedOutput()
			t.Fatalf("%s\n", stdoutStderr)
		} else {
			t.Fatalf("content of %s did not match %s.\nContent A:\n%s\n\nContent B:\n%s\n", a, b, contentA, contentB)
		}
	}
}
