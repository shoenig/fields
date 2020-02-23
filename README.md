fields
======

Use `fields` to chop input into selectable columns.

[![Go Report Card](https://goreportcard.com/badge/gophers.dev/cmds/fields)](https://goreportcard.com/report/gophers.dev/cmds/fields)
[![Build Status](https://travis-ci.org/shoenig/fields.svg?branch=master)](https://travis-ci.org/shoenig/fields)
[![GoDoc](https://godoc.org/gophers.dev/cmds/fields?status.svg)](https://godoc.org/gophers.dev/cmds/fields)
![NetflixOSS Lifecycle](https://img.shields.io/osslifecycle/shoenig/fields.svg)
![GitHub](https://img.shields.io/github/license/shoenig/fields.svg)

# Project Overview

Module `gophers.dev/cmds/fields` provides a command-line utility for processing
columns of input text.

# Getting Started

#### Install from SnapCraft

The `fields` command can be installed as a snap
```bash
$ sudo snap install fields
```

#### Build from source

The `fields` command can be compiled by running
```bash
$ go get gophers.dev/cmds/fields/cmd/fields
```

# Example Usages

#### select a single column (from left)
```bash
$ fields 3 <<< "a b c d e f g"
c
```

#### select a single column (from right)
```bash
$ fields -- -3 <<< "a b c d e f g"
e
```

#### select multiple columns
```bash
$ fields 1,-1,2,-2 <<< "a b c d e f g"
a g b f
```

#### select columns to the right of N (from left)
```bash
$ fields 4: <<< "a b c d e f g"
d e f g
```

#### select columns to the right of N (from right)
```bash
$ fields -- -2: <<< "a b c d e f g"
f g
```

#### select columns to the left of N (from left)
```bash
$ fields :2 <<< "a b c d e f g"
a b
```

#### select columns to the left of N (from right)
```bash
$ fields :-2 <<< "a b c d e f g"
a b c d e f
```

#### select range of columns
```bash
$ fields 2:5 <<< "a b c d e f g"
b c d e
```

#### any combination of the above, all together
```bash
$ fields 1,2,-2,3:5,2:,:-3 <<< "a b c d e f g"
a b f c d e b c d e f g a b c d e
```

# Contributing

The `gophers.dev/cmds/fields` module is always improving with new features
and error corrections. For contributing bug fixes and new features please file an issue.

# License

The `gophers.dev/cmds/fields` module is open source under the [MIT](LICENSE) license.
