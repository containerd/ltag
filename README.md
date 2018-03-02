### ltag

Prepends project files with given template.

- Can be used for adding licence or copyright information on src files of project.
- Skip file, if template (as provided) already present 
- Supports Golang, Dockerfile, Makefiles and bash scripts
      - Take cares of Golang compiler flags.
      - Take cares of Golang Package comments too.


#### Install

```
go get github.com/kunalkushwaha/ltag
```


#### Usage
```
$ ltag
$ ltag --help
Usage of ltag:
  -check
        check files missing header
  -excludes string
        exclude folders (default "vendor")
  -path string
        project path (default ".")
  -t string
        template files path (default "./template")
  -v    verbose output

```

### Example

```
$ ltag -t=template.txt -path=tempProj
Following files are updated
tempProj/abc.go
tempProj/src/lvl1/temp.go
```
