package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDockerfileApplier_ApplyHeader(t *testing.T) {
	tc := newTagContext(t)
	defer func() { _ = tc.templateFiles.dTemplateFile.Close() }()

	tmpDir := t.TempDir()

	files := []string{
		"Dockerfile.nodirectives",
		"Dockerfile.comments",
		"Dockerfile.single-directive",
		"Dockerfile.multiple-directives",
	}

	d := dockerfileApplier{}
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

func TestDockerfileApplier_CheckHeader(t *testing.T) {
	files := []string{
		// The non-golden files don't have a header present
		"Dockerfile.nodirectives",
		"Dockerfile.comments",
		"Dockerfile.single-directive",
		"Dockerfile.multiple-directives",

		// The golden files should have the header present
		"Dockerfile.nodirectives.golden",
		"Dockerfile.comments.golden",
		"Dockerfile.single-directive.golden",
		"Dockerfile.multiple-directives.golden",
	}

	d := dockerfileApplier{}
	tc := newTagContext(t)
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

func newTagContext(t *testing.T) TagContext {
	t.Helper()
	templateFile, err := loadTemplate("./testdata/templates/", "dockerfile.txt")
	if err != nil {
		t.Fatalf("failed to load dockerfile template")
	}
	return TagContext{
		templateFiles: TemplateFiles{dTemplateFile: templateFile},
		templatePath:  "./testdata/templates/",
	}
}
