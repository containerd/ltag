package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestBashApplier_ApplyHeader(t *testing.T){
	tc := newBashTagContext(t)
	defer func() { _ = tc.templateFiles.dTemplateFile.Close()}()


	tmpDir, err := ioutil.TempDir("", t.Name())
	if err != nil {
		t.Fatalf("failed to create temp directory")
	}
	defer func() { _ = os.RemoveAll(tmpDir)}()

	files := []string{
		"bash.simple",
		"bash.comments",
		"bash.shebang",
		"bash.env.shebang",
		"bash.sh.shebang",
		"bash.make.shebang",
	}

	d := bashApplier{}
	for _, f := range files {
		fileName := f
		t.Run(fileName, func(t *testing.T){
			testfile := filepath.Join(tmpDir, fileName)
			copyBashFile(t, "./testdata/"+fileName, testfile)

			err = d.ApplyHeader(testfile, &tc)
			if err != nil {
				t.Fatalf("failed to apply header to %s: %+v", testfile, err)
			}
			compareBashFiles(t, testfile, "./testdata/"+fileName+".golden")
		})
	}
}

func TestBashFileApplier_CheckHeader(t *testing.T){
	files := []string{
		// The non-golden files don't have a header present
		"bash.simple",
		"bash.comments",
		"bash.shebang",
		"bash.env.shebang",
		"bash.sh.shebang",
		"bash.make.shebang",

		// The golden files should have the header present
		"bash.simple.golden",
		"bash.comments.golden",
		"bash.shebang.golden",
		"bash.env.shebang.golden",
		"bash.sh.shebang.golden",
		"bash.make.shebang.golden",
	}

	d := bashApplier{}
	tc := newBashTagContext(t)
	defer func() { _ = tc.templateFiles.dTemplateFile.Close()}()
	for _, f := range files {
		fileName := f
		t.Run(fileName, func(t *testing.T){
			f, err := os.Open("./testdata/"+ fileName)
			if err != nil {
				t.Fatalf("failed to optn file %s: %+v", fileName, err)
			}
			defer func() { _ = f.Close() }()
			found, err := d.CheckHeader(f, &tc)
			if err != nil {
				t.Fatalf("failed to check header: %+v", err)
			}
			expected := strings.HasSuffix(fileName, ".golden")
			if found != expected {
				t.Fail()
			}
		})
	}
}

func newBashTagContext(t *testing.T) TagContext{
	t.Helper()
	templateFile, err := loadTemplate( "./testdata/templates/", "bash.txt")
	if err != nil {
		t.Fatalf("failed to load bash template")
	}
	return TagContext{
		templateFiles: TemplateFiles{shTemplateFile: templateFile},
		templatePath:   "./testdata/templates/",
	}
}

func copyBashFile(t *testing.T, src, dst string) {
	t.Helper()
	input, err := ioutil.ReadFile(src)
	if err != nil {
		t.Fatalf("failed to read %s: %+v", src, err)
	}
	err = ioutil.WriteFile(dst, input, 0777)
	if err != nil {
		t.Fatalf("failed to write %s: %+v", dst, err)
	}
}

func compareBashFiles(t *testing.T, a, b string) {
	contentA, err := ioutil.ReadFile(a)
	if err != nil {
		t.Fatalf("failed to read file %s: %+v", a, err)
	}
	contentB, err := ioutil.ReadFile(b)
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
