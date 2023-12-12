# go-makefile-gen

![logo](go-makefile-gen.png)

A simple utility tool for generating a [Makefile](https://en.wikipedia.org/wiki/Make_(software)#Makefiles) to your [Go](https://go.dev/) project. It also offers the ability of adding a new target to a given [Makefile](https://en.wikipedia.org/wiki/Make_(software)#Makefiles).

## installation

```
go install github.com/tiagomelo/go-makefile-gen
```

## usage

### creating a `Makefile`

```
gomakefile generate 
```

It generates this `Makefile` in the current directory:

```
.PHONY: help
## help: shows this help message
help:
	@ echo "Usage: make [target]\n"
	@ sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: test
## test: run unit tests
test:
	@ go test -v ./... -count=1

.PHONY: coverage
## coverage: run unit tests and generate coverage report in html format
coverage:
	@ go test -coverprofile=coverage.out ./...  && go tool cover -html=coverage.out
```

You can also specify the desired path for it:

```
gomakefile generate -p <path/to/Makefile>
```

### adding a new target to a `Makefile`

```
gomakefile addtarget -t "my-new-target"
```

It adds the desired target to the `Makefile` that is present in the currenty directory:

```
.PHONY: my-new-target
## my-new-target: explain what my-new-target does
my-new-target:
```

You can also specify the path for the existing `Makefile` you want to add the target to:

```
gomakefile addtarget -t "my-new-target" -p <path/to/Makefile>
```

## unit tests

```
make test
```

## unit tests coverage

```
make coverage
```