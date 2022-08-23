package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGolangApplier_ApplyHeader(t *testing.T) {
	tc := newGolangTagContext(t)
	defer func() { _ = tc.templateFiles.dTemplateFile.Close() }()

	tmpDir := t.TempDir()

	files := []string{
		"go.basic",
		"go.generated",
		"go.single-buildtag",
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
	tests := []struct {
		fileName string
		expected bool
	}{
		{
			fileName: "go.basic",
			expected: false,
		},
		{
			fileName: "go.generated",
			expected: true, // Generated files are not checked, and always "ok"
		},
		{
			fileName: "go.single-buildtag",
			expected: false,
		},
		{
			fileName: "go.basic.golden",
			expected: true, // Generated files are not checked, and always "ok"
		},
		{
			fileName: "go.generated.golden",
			expected: true, // Generated files are not checked, and always "ok"
		},
		{
			fileName: "go.single-buildtag.golden",
			expected: true,
		},
	}
	g := golangApplier{}
	tagContext := newGolangTagContext(t)
	defer func() { _ = tagContext.templateFiles.dTemplateFile.Close() }()
	for _, tc := range tests {
		tc := tc
		t.Run(tc.fileName, func(t *testing.T) {
			f, err := os.Open("./testdata/" + tc.fileName)
			if err != nil {
				t.Fatalf("failed to open file %s: %+v", tc.fileName, err)
			}
			defer func() { _ = f.Close() }()
			found, err := g.CheckHeader(f, &tagContext)

			if err != nil {
				t.Fatalf("failed to check header: %+v", err)
			}
			if found != tc.expected {
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
