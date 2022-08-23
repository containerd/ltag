package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type bashApplier struct {
}

func (g *bashApplier) CheckHeader(target *os.File, t *TagContext) (bool, error) {

	// Check compiler flags.
	sbFlag, sbBuf, err := g.checkSheBang(target)
	if err != nil {
		return false, err
	}
	target.Seek(0, 0)

	tbuf, err := os.ReadFile(filepath.Join(t.templatePath, "bash.txt"))
	if err != nil {
		return false, err
	}

	var templateBuf string
	if sbFlag {
		templateBuf = fmt.Sprintf("%s%s%s", sbBuf, "\n\n", tbuf)
	} else {
		templateBuf = string(tbuf)
	}

	targetBuf := make([]byte, len(templateBuf))

	n, err := target.Read(targetBuf)
	if err != nil {
		return false, err
	}

	if n == len(templateBuf) {
		if strings.Compare(templateBuf, string(targetBuf)) == 0 {
			return true, nil
		}
	}

	return false, nil
}

func (g *bashApplier) ApplyHeader(path string, t *TagContext) error {

	file, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	t.templateFiles.goTemplateFile.Seek(0, 0)

	headerExist, err := g.CheckHeader(file, t)
	if err != nil {
		return err
	}

	if headerExist {
		return nil
	}

	// Reset the read pointers to begining of file.
	t.templateFiles.goTemplateFile.Seek(0, 0)
	file.Seek(0, 0)

	sbFlag, sbBuf, err := g.checkSheBang(file)
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
	if sbFlag {
		tFile.Write(sbBuf)
		tFile.Write([]byte("\n\n"))
		_, err = reader.Discard(len(sbBuf))
	}

	t.templateFiles.shTemplateFile.Seek(0, 0)
	_, err = io.Copy(tFile, t.templateFiles.shTemplateFile)
	if err != nil {
		return err
	}

	_, err = io.Copy(tFile, reader)
	if err != nil {
		return err
	}

	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	//	info.Mode

	err = os.Rename(tempFile, path)
	if err != nil {
		return err
	}
	err = os.Chmod(path, info.Mode().Perm())
	if err != nil {
		return err
	}
	return nil
}

var parserSheBang = regexp.MustCompile(`^#!(.*)`)

func (g *bashApplier) checkSheBang(target *os.File) (bool, []byte, error) {
	reader := bufio.NewReader(target)
	buf, _, err := reader.ReadLine()
	if err != nil {
		return false, nil, err
	}

	if parserSheBang.Match(buf) {
		return true, buf, nil
	}

	return false, nil, nil
}
