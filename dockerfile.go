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

type dockerfileApplier struct {
}

func (g *dockerfileApplier) CheckHeader(target *os.File, t *TagContext) (bool, error) {
	dFlag, dBuf, err := g.checkParserDirectives(target)
	if err != nil {
		return false, err
	}
	target.Seek(0, 0)

	tbuf, err := os.ReadFile(filepath.Join(t.templatePath, "dockerfile.txt"))
	if err != nil {
		return false, err
	}

	var templateBuf string
	if dFlag {
		templateBuf = fmt.Sprintf("%s%s%s", dBuf, "\n\n", tbuf)
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

func (g *dockerfileApplier) ApplyHeader(path string, t *TagContext) error {
	file, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	t.templateFiles.dTemplateFile.Seek(0, 0)

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

	dFlag, dBuf, err := g.checkParserDirectives(file)
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
	if dFlag {
		tFile.Write(dBuf)
		tFile.Write([]byte("\n\n"))
		_, err = reader.Discard(len(dBuf))
	}

	t.templateFiles.dTemplateFile.Seek(0, 0)
	_, err = io.Copy(tFile, t.templateFiles.dTemplateFile)
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

var parserDirective = regexp.MustCompile(`^#\s*\w+=.*`)

// checkParserDirectives captures Dockerfile parser directives.
//
// Parser directives are optional, and affect the way in which subsequent lines
// in a Dockerfile are handled. Parser directives are written as a special type
// of comment in the form # directive=value. A single directive may only be used
// once.
//
// Once a comment, empty line or builder instruction has been processed, Docker
// no longer looks for parser directives. Instead it treats anything formatted
// as a parser directive as a comment and does not attempt to validate if it
// might be a parser directive. Therefore, all parser directives must be at the
// very top of a Dockerfile.
//
// Parser directives are not case-sensitive. However, convention is for them to
// be lowercase. Convention is also to include a blank line following any parser
// directives. Line continuation characters are not supported in parser directives.
func (g *dockerfileApplier) checkParserDirectives(target *os.File) (bool, []byte, error) {
	reader := bufio.NewReader(target)
	line, err := reader.ReadBytes('\n')
	if err != nil {
		return false, nil, err
	}

	if !parserDirective.Match(line) {
		return false, nil, nil
	}
	buf := append([]byte{}, line...)

	for {
		line, err = reader.ReadBytes('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return false, nil, err
		}
		if !parserDirective.Match(line) {
			// found a non-directive: stop looking for more directives,
			// and return those that were found
			return true, buf, nil
		}
		buf = append(buf, line...)
	}
	return true, buf, nil
}
