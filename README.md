go-env
======

[![Latest release][latest-release-img]][latest-release-url]
[![Build status][build-status-img]][build-status-url]
[![Go Report Card][report-img]][report-url]
[![Documentation][doc-img]][doc-url]

[latest-release-img]: https://img.shields.io/github/release/roeldev/go-env.svg?label=latest
[latest-release-url]: https://github.com/roeldev/go-env/releases
[build-status-img]: https://github.com/roeldev/go-env/workflows/Go/badge.svg
[build-status-url]: https://github.com/roeldev/go-env/actions?query=workflow%3AGo
[report-img]: https://goreportcard.com/badge/github.com/roeldev/go-env
[report-url]: https://goreportcard.com/report/github.com/roeldev/go-env
[doc-img]: https://godoc.org/github.com/roeldev/go-env?status.svg
[doc-url]: https://pkg.go.dev/github.com/roeldev/go-env

Env is a Go package for parsing environment variables. It can be used to read from any file or `io.Reader` or to parse the values of `os.Environ()` and `os.Args`.


```sh
go get github.com/roeldev/go-env
```
```go
import "github.com/roeldev/go-env"
```

## Using files
It is possible to read from any `io.Reader` or open a file and parse its contents. Supported formats:
- `key=value`; value without quotes
- `key="value"`; value with double quotes
- `key='value'`; value with single quotes

```go
func main() {
	// create Map to store results
	envs := make(env.Map)
	// open file, parse it and store results
	_, err := env.Open(".env", envs)
	if err != nil {
		// handle error
	}

	// use envs any way you like
}
```

## Using `io.Reader`
It is possible to read from any `io.Reader`. Supported formats are the same as when reading from files.

```go
func main() {
	// somehow get a reader
	// reader := strings.NewReader(str)

	// create Map to store results
	envs := make(env.Map)
	// read + parse from reader and store results
	_, err := env.Read(reader, envs)
	if err != nil {
		// handle error
	}

	// use envs any way you like
}
```

## Using CLI arguments
It is possible to extract environment variables from `os.Args`. Supported input formats from CLI are:
- `-e=key=val`; flag with single dash
- `-e key=val`
- `--e=key=val`; flag with double dash
- `--e key=val`

```go
func main() {
	// create Map to store results
	envs := make(env.Map)
	// parse flags and store them in the provided Map
	n := env.ParseFlagArgs("e", os.Args[1:], envs)
	
	// use envs any way you like
}
```
The args slice may have multiple entries for the same flag provided by CLI arguments, eg:
```sh
yourcliprogram -e key=val -someFlag -e="another=env var"
```
Flags may also contain multiple values, thus above example can also be written like:
```sh
yourcliprogram -e key=val another="env var" -someFlag
```
When parsed, both example result in:
```go
env.Map{"key": "val", "another": "env var"}
```


## Examples
Some example programs are provided in the [examples](examples) folder. They further illustrate some very basic ways to use this package.


## Documentation
Additional detailed documentation is available at [go.dev][doc-url]


### Created with
<a href="https://www.jetbrains.com/?from=roeldev/go-env" target="_blank"><img src="https://pbs.twimg.com/profile_images/1206615658638856192/eiS7UWLo_400x400.jpg" width="35" /></a>


## License
[GPL-3.0+](LICENSE) © 2020 [Roel Schut](https://roelschut.nl)
