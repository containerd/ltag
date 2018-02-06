package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type TagContext struct {
	sourceExt    string
	excludeList  []string
	templatePath string
	templateFile *os.File
	dryRun       bool
}

func main() {
	ppath := flag.String("path", ".", "project path")
	srcExt := flag.String("ext", ".go", "file extention for tagging")
	excludes := flag.String("excludes", "vendor", "exclude folders")
	tpath := flag.String("t", "", "template file path")
	dryRun := flag.Bool("d", false, "dry run")
	flag.Parse()

	if *tpath == "" {
		fmt.Println("template path missing")
		flag.Usage()
		return
	}

	excludeList := strings.Split(*excludes, " ")

	tfile, err := os.OpenFile(*tpath, os.O_RDONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer tfile.Close()

	t := TagContext{sourceExt: *srcExt, excludeList: excludeList, templatePath: *tpath, templateFile: tfile, dryRun: *dryRun}

	if *dryRun {
		fmt.Println("Following files can be updated")
	} else {
		fmt.Println("Following files are updated")
	}

	err = filepath.Walk(*ppath, t.tagFiles)
	if err != nil {
		panic(err)
	}
}

func (t *TagContext) tagFiles(path string, f os.FileInfo, err error) error {
	if (f.Name() == ".git" || f.Name() == ".svn" || f.Name() == "..") && f.IsDir() {
		return filepath.SkipDir
	}

	for _, exclude := range t.excludeList {
		if f.Name() == exclude {
			return filepath.SkipDir
		}
	}

	if !f.IsDir() && filepath.Ext(f.Name()) == t.sourceExt && f.Size() > 0 {

		file, err := os.OpenFile(path, os.O_RDONLY, 0666)
		if err != nil {
			return err
		}
		defer file.Close()
		t.templateFile.Seek(0, 0)

		headerExist, err := t.checkTemplateHeader(file)
		if err != nil {
			return err
		}

		if headerExist {
			return nil
		}
		// Prints the file requires update
		fmt.Println(path)

		if t.dryRun {
			return nil
		}

		//Reset the read pointers to begining of file.
		t.templateFile.Seek(0, 0)
		file.Seek(0, 0)

		compilerFlagExist, flags, err := t.checkCompilerFlags(file)
		if err != nil {
			return err
		}
		file.Seek(0, 0)

		tempFile := path + ".tmp"
		tFile, err := os.OpenFile(tempFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return err
		}
		defer tFile.Close()

		reader := bufio.NewReader(file)
		if compilerFlagExist {
			tFile.Write(flags)
			tFile.Write([]byte("\n\n"))
			_, _, err = reader.ReadLine()
			_, _, err = reader.ReadLine()
		}

		_, err = io.Copy(tFile, t.templateFile)
		if err != nil {
			return err
		}

		_, err = io.Copy(tFile, reader)
		if err != nil {
			return err
		}

		err = os.Rename(tempFile, path)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *TagContext) checkTemplateHeader(target *os.File) (bool, error) {
	buf, err := ioutil.ReadFile(t.templatePath)
	if err != nil {
		return false, err
	}

	targetBuf := make([]byte, len(buf))

	n, err := target.Read(targetBuf)
	if err != nil {
		return false, err
	}

	if n == len(buf) {
		if strings.Compare(string(buf), string(targetBuf)) == 0 {
			return true, nil
		}
	}

	return false, nil
}

func (t *TagContext) checkCompilerFlags(target *os.File) (bool, []byte, error) {

	reader := bufio.NewReader(target)
	buf, _, err := reader.ReadLine()
	if err != nil {
		return false, nil, err
	}

	if strings.HasPrefix(string(buf), "//") && (strings.Contains(string(buf), "build") ||
		strings.Contains(string(buf), "unix") ||
		strings.Contains(string(buf), "linux") ||
		strings.Contains(string(buf), "windows") ||
		strings.Contains(string(buf), "darwin") ||
		strings.Contains(string(buf), "freebsd")) {
		return true, buf, nil
	}
	return false, nil, nil
}
