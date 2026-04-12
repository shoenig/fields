fields
======

Use the `fields` CLI command as a modern replacement for `awk` + `cut`.
With fields you specify which columns of text you want, in a flexible format.

# Getting Started

#### Download from releases

Available to download pre-compiled binaries on [releases](https://github.com/shoenig/fields/releases)

#### Build from source

The `fields` command can be installed via Go by running
```bash
$ go install github.com/shoenig/fields@latest
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

# License

The `github.com/shoenig/fields` module is open source under the [MIT](LICENSE) license.
