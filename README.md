# fields

`fields` is a new and enchanced replacement for the coreutils `cut` command

[![Go Report Card](https://goreportcard.com/badge/github.com/shoenig/fields)](https://goreportcard.com/report/github.com/shoenig/fields) [![Build Status](https://travis-ci.org/shoenig/fields.svg?branch=master)](https://travis-ci.org/shoenig/fields) [![GoDoc](https://godoc.org/github.com/shoenig/fields?status.svg)](https://godoc.org/github.com/shoenig/fields) [![License](https://img.shields.io/github/license/shoenig/fields.svg?style=flat-square)](LICENSE)

### Install
Like any typical Go executable, just use `go get` to install if you have a Go work
environment setup. Otherwise, ask a friend.

`go get github.com/shoenig/fields`

### Usage

The `fields` command reads input from stdin, so you can use the usual bash
file redirection or piped input.

`fields [flags] <columns/spans>`

#### Flags

`--help` : print a help message
`--no-newline` : prevents a newline from being printed at the end of output

### Examples

Selecting columns 2 and 3 (remember, fields is zero based)

```bash
$ ./fields 2,3 <<< "A man a plan a canal panama"
a plan
```

Selecting columns 3-5,1 (support for spans)

```bash
$ ./fields 3-5,1 <<< "A man a plan a canal panama"
plan a canal man
```

Selecting columns 6,4-0 (support for reverse spans)

```bash
$ ./fields 6,4-0 <<< "A man a plan a canal panama"
panama a plan a man A
```

Selecting column 3 with no newline
```bash
$ ./fields --no-newline 3 <<< "A man a plan a canal panama"
plan$
```