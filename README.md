### ltag

Prepends project files with given template.

- Can be used for adding licence or copyright information on src files of project.
- Skip file, if template (as provided) already present 
- Supports Golang source files, Dockerfile, Makefiles and bash scripts
   - Take cares of compiler flags for golang source files and shebang of bash scripts.
   - Take cares of Golang Package comments too.


#### Install

```
go install github.com/kunalkushwaha/ltag@latest
```

> [!NOTE]
>
> The module name is planned to be renamed to `github.com/containerd/ltag`.

#### Usage

``` console
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

To Apply header from `./template` folder

``` console
$ ltag  -path=temp -v
Files modified :  11
temp/Dockerfile/Dockerfile
...
```

To Check if files missing header

``` console
$ ltag  -path=temp --check -v
temp/Dockerfile/Dockerfile
temp/Dockerfile/abc.dockerfile
temp/src/lvl1/doc.go
temp/src/lvl1/temp.go
```
