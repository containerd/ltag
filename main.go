package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//Return values
const (
	FileReadError = 0
	NormalFiles   = 1
	CompilerFlags = 2
	AutoGenerated = 3
)

//Applier interface required to implement for new file type
type Applier interface {
	CheckHeader(target *os.File, t *TagContext) (bool, error)
	ApplyHeader(path string, t *TagContext) error
}

//TagContext keeps context info for Applier
type TagContext struct {
	excludeList   []string
	templatePath  string
	templateFiles TemplateFiles
	dryRun        bool
	outfileList   []string
}

//TemplateFiles stores template of header
type TemplateFiles struct {
	goTemplateFile *os.File
	shTemplateFile *os.File
	dTemplateFile  *os.File
	mTemplateFile  *os.File
}

func main() {
	ppath := flag.String("path", ".", "project path")
	excludes := flag.String("excludes", "vendor", "exclude folders")
	tpath := flag.String("t", "./template", "template files path")
	dryRun := flag.Bool("check", false, "check files missing header")
	verbose := flag.Bool("v", false, "verbose output")
	flag.Parse()

	dTFile, err := loadTemplate(*tpath, "dockerfile.txt")
	if err != nil {
		fmt.Println("No template file for Dockerfile, shall skip all Dockerfile")
	}
	if dTFile != nil {
		defer dTFile.Close()
	}

	goTFile, err := loadTemplate(*tpath, "go.txt")
	if err != nil {
		fmt.Println("No template file for golang files, shall skip all golang files")
	}
	if goTFile != nil {
		defer goTFile.Close()
	}

	bashTFile, err := loadTemplate(*tpath, "bash.txt")
	if err != nil {
		fmt.Println("No template file for bash scripts, shall skip all bash scripts")
	}
	if bashTFile != nil {
		defer bashTFile.Close()
	}

	makeTFile, err := loadTemplate(*tpath, "makefile.txt")
	if err != nil {
		fmt.Println("No template file for Makefile, shall skip all makefiles")
	}
	if makeTFile != nil {
		defer makeTFile.Close()
	}

	excludeList := strings.Split(*excludes, " ")

	templateFiles := TemplateFiles{
		mTemplateFile:  makeTFile,
		shTemplateFile: bashTFile,
		goTemplateFile: goTFile,
		dTemplateFile:  dTFile}

	t := TagContext{
		excludeList:   excludeList,
		templateFiles: templateFiles,
		templatePath:  *tpath,
		dryRun:        *dryRun}

	//TODO:
	// itterate to all template file handlers, if not nill defer Close.

	err = filepath.Walk(*ppath, t.tagFiles)
	if err != nil {
		panic(err)
	}

	if !*dryRun {
		fmt.Println("Files modified : ", len(t.outfileList))
	} else {
		fmt.Println("Files missing header : ", len(t.outfileList))
	}
	if *verbose {
		for _, path := range t.outfileList {
			fmt.Println(path)
		}
	}
}

func (t *TagContext) tagFiles(path string, f os.FileInfo, err error) error {

	var applier Applier
	processed := false

	if (f.Name() == ".git" || f.Name() == ".svn" || f.Name() == "..") && f.IsDir() {
		return filepath.SkipDir
	}

	if f.IsDir() {
		for _, exclude := range t.excludeList {
			if f.Name() == exclude {
				return filepath.SkipDir
			}
		}
	}

	if !f.IsDir() && f.Size() > 0 {

		if f.Name() == "LICENSE" || f.Name() == "MAINTAINERS" {
			return nil
		}

		file, err := os.OpenFile(path, os.O_RDONLY, 0666)
		if err != nil {
			return err
		}
		defer file.Close()

		fname := strings.Split(f.Name(), ".")
		if len(fname) == 1 { //Without extension.
			if f.Mode()&0111 != 0 && t.templateFiles.shTemplateFile != nil {
				applier = &bashApplier{}
				processed = true
			} else if fname[0] == "Makefile" && t.templateFiles.mTemplateFile != nil {
				applier = &makefileApplier{}
				processed = true
			} else if strings.ToLower(fname[0]) == "dockerfile" && t.templateFiles.dTemplateFile != nil {
				applier = &dockerfileApplier{}
				processed = true
			}
		} else {
			if fname[1] == "go" && t.templateFiles.goTemplateFile != nil {
				applier = &golangApplier{}
				processed = true
			}
			if strings.ToLower(fname[1]) == "dockerfile" && t.templateFiles.dTemplateFile != nil {
				applier = &dockerfileApplier{}
				processed = true
			}
			if strings.ToLower(fname[0]) == "makefile" && t.templateFiles.mTemplateFile != nil {
				applier = &dockerfileApplier{}
				processed = true
			}
			if fname[1] == "sh" && t.templateFiles.shTemplateFile != nil {
				applier = &bashApplier{}
				processed = true
			}
		}
		if !processed {
			return nil
		}
		processed = false
		headerExist, err := applier.CheckHeader(file, t)
		if err != nil {
			return err
		}
		if headerExist {
			return nil
		}

		if !t.dryRun {
			err = applier.ApplyHeader(path, t)
			if err != nil {
				return err
			}
		}
		t.outfileList = append(t.outfileList, path)
	}
	return nil
}

func loadTemplate(path string, name string) (*os.File, error) {
	templateFile := filepath.Join(path, name)
	tFile, err := os.OpenFile(templateFile, os.O_RDONLY, 0666)
	return tFile, err
}
