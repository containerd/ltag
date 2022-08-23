package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBashApplier_ApplyHeader(t *testing.T) {
	tc := newBashTagContext(t)
	defer func() { _ = tc.templateFiles.dTemplateFile.Close() }()

	tmpDir := t.TempDir()

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
		t.Run(fileName, func(t *testing.T) {
			testfile := filepath.Join(tmpDir, fileName)
			copyFile(t, "./testdata/"+fileName, testfile)

			err := d.ApplyHeader(testfile, &tc)
			if err != nil {
				t.Fatalf("failed to apply header to %s: %+v", testfile, err)
			}
			compareFiles(t, testfile, "./testdata/"+fileName+".golden")
		})
	}
}

func TestBashFileApplier_CheckHeader(t *testing.T) {
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
	defer func() { _ = tc.templateFiles.dTemplateFile.Close() }()
	for _, f := range files {
		fileName := f
		t.Run(fileName, func(t *testing.T) {
			f, err := os.Open("./testdata/" + fileName)
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

func newBashTagContext(t *testing.T) TagContext {
	t.Helper()
	templateFile, err := loadTemplate("./testdata/templates/", "bash.txt")
	if err != nil {
		t.Fatalf("failed to load bash template")
	}
	return TagContext{
		templateFiles: TemplateFiles{shTemplateFile: templateFile},
		templatePath:  "./testdata/templates/",
	}
}
