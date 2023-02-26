fields
======

Use the `fields` CLI command as a modern replacement for `awk` + `cut`.
With fields you specify which columns of text you want, in a flexible format.

![GitHub](https://img.shields.io/github/license/shoenig/fields.svg)
[![run-ci](https://github.com/shoenig/fields/actions/workflows/ci.yml/badge.svg)](https://github.com/shoenig/fields/actions/workflows/ci.yml)

# Getting Started

#### Build from source

The `fields` command can be installed via Go by running
```bash
$ go install github.com/shoenig/fields/cmd/fields@latest
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

The `github.com/shoenig/fields` module is always improving with new features
and error corrections. For contributing bug fixes and new features please file an issue.

# License

The `github.com/shoenig/fields` module is open source under the [MIT](LICENSE) license.
