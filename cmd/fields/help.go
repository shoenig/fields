package main

const helpText = `fields - CLI columnar text processor [v1]
	fields is a tool for selecting columns of text. One or more individual columns
	and/or ranges of columns can be selected by index. Indexes can be positive
	(counting from the left) or negative (counting from the right). Ranges can
	bounded or unbounded (e.g. select all input to the left or right of N).

	Input is expected on STDIN. Typically input is provided to fields through
	pipes or redirection.

	Usage)
	fields [options] <columns>

	Options)
	- newline: print a trailing newline character (default true)

	Syntax)
	- Column indexes are counted starting from 1
	- A positive index means start counting from the left
	- A negative index means start counting from the right
	- Closed ranges are denoted by two indexes separated by ':'
	- Right open ranges are denoted by an index and a ':'
	- Left open ranges are denoted by a ':' and an index
	- High:Low ranges can reverse the order of columns
	- Combine selections using comma separated list

	Examples)
	select a single column (from left)
	  $ fields 3 <<< "a b c d e f g"			=> c

	select a single column (from right)
	  $ fields -- -3 <<< "a b c d e f g"			=> e

	select columns to the right of N (from left)
	  $ fields 4: <<< "a b c d e f g"			=> d e f h

	select columns to the right of N (from right)
	  $ fields -- -2: <<< "a b c d e f g"			=> f g

	select range of columns
	  $ fields 2:5 <<< "a b c d e f g"			=> b c d e

	select range of columns (reverse)
	  $ fields 5:2 <<< "a b c d e f g"			=> e d c b

	select multiple combinations
	  $ fields -- -2,3:5,6: <<< "a b c d e f g"		=> f c d e f g

	Bugs)
	Report bugs & feature requests @ https://github.com/shoenig/fields/issues
`
