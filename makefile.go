package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type makefileApplier struct {
}

func (g *makefileApplier) CheckHeader(target *os.File, t *TagContext) (bool, error) {
	tbuf, err := ioutil.ReadFile(filepath.Join(t.templatePath, "makefile.txt"))
	if err != nil {
		return false, err
	}

	templateBuf := string(tbuf)

	targetBuf := make([]byte, len(templateBuf))

	n, err := target.Read(targetBuf)
	if err != nil {
		return false, err
	}

	if n == len(templateBuf) {
		if strings.Compare(string(templateBuf), string(targetBuf)) == 0 {
			return true, nil
		}
	}
	return false, nil
}

func (g *makefileApplier) ApplyHeader(path string, t *TagContext) error {
	file, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	t.templateFiles.mTemplateFile.Seek(0, 0)

	headerExist, err := g.CheckHeader(file, t)
	if err != nil {
		return err
	}

	if headerExist {
		return nil
	}

	//Reset the read pointers to begining of file.
	t.templateFiles.mTemplateFile.Seek(0, 0)
	file.Seek(0, 0)

	tempFile := path + ".tmp"
	tFile, err := os.OpenFile(tempFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer tFile.Close()

	reader := bufio.NewReader(file)
	t.templateFiles.mTemplateFile.Seek(0, 0)
	_, err = io.Copy(tFile, t.templateFiles.mTemplateFile)
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
	return nil
}
