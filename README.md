### ltag

Prepends project files with given template.

- Can be used for adding licence or copyright information on src files of project.
- Skip file, if template (as provided) already present 
- Take cares of Golang compiler flags.

#### Install

```
go get github.com/kunalkushwaha/ltag
```


#### Usage
```
$ ltag
template path missing
Usage of ./ltag:
  -d    dry run
  -exclude string
        exclude folder (default "vendor")
  -ext string
        file extention for tagging (default ".go")
  -path string
        project path (default ".")
  -t string
        template file path
```

### Example

```
$ ltag -t=template.txt -path=tempProj
Following files are updated
tempProj/abc.go
tempProj/src/lvl1/temp.go
```
