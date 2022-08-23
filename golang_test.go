package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGolangApplier_ApplyHeader(t *testing.T) {
	tc := newGolangTagContext(t)
	defer func() { _ = tc.templateFiles.dTemplateFile.Close() }()

	tmpDir := t.TempDir()

	files := []string{
		"go.basic",
	}

	g := golangApplier{}
	for _, f := range files {
		fileName := f
		t.Run(fileName, func(t *testing.T) {
			testfile := filepath.Join(tmpDir, fileName)
			copyFile(t, "./testdata/"+fileName, testfile)

			err := g.ApplyHeader(testfile, &tc)
			if err != nil {
				t.Fatalf("failed to apply header to %s: %+v", testfile, err)
			}
			compareFiles(t, testfile, "./testdata/"+fileName+".golden")
		})
	}
}

func TestGolangApplier_CheckHeader(t *testing.T) {
	files := []string{
		// The non-golden files don't have a header present
		"go.basic",

		// The golden files should have the header present
		"go.basic.golden",
	}

	g := golangApplier{}
	tc := newGolangTagContext(t)
	defer func() { _ = tc.templateFiles.dTemplateFile.Close() }()
	for _, f := range files {
		fileName := f
		t.Run(fileName, func(t *testing.T) {
			f, err := os.Open("./testdata/" + fileName)
			if err != nil {
				t.Fatalf("failed to optn file %s: %+v", fileName, err)
			}
			defer func() { _ = f.Close() }()
			found, err := g.CheckHeader(f, &tc)
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

func newGolangTagContext(t *testing.T) TagContext {
	t.Helper()
	templateFile, err := loadTemplate("./testdata/templates/", "go.txt")
	if err != nil {
		t.Fatalf("failed to load Go template")
	}
	return TagContext{
		templateFiles: TemplateFiles{goTemplateFile: templateFile},
		templatePath:  "./testdata/templates/",
	}
}
