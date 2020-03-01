#!/usr/bin/env bash

set -euo pipefail

PATH="${GOPATH}/bin"

function test {
 escape="${1}"
 cols="${2}"
 data="${3}"
 expect="${4}"
 result=$(fields ${escape} "${cols}" <<< "${data}")
 if [ "${expect}" != "${result}" ]; then
	echo "expected: ${expect}, got: ${result}"
	exit 1
 fi
}

function testn {
 test "--" "${1}" "${2}" "${3}"
}

function testp {
 test "" "${1}" "${2}" "${3}"
}

# single
testp 1 "a b c d e f g" "a"
testp 4 "a b c d e f g" "d"
testp 7 "a b c d e f g" "g"
testn -1 "a b c d e f g" "g"
testn -3 "a b c d e f g" "e"
testn -7 "a b c d e f g" "a"

# range bounded
testp 1:3 "a b c d e f g" "a b c"
testp 4:7 "a b c d e f g" "d e f g"
testp 5:3 "a b c d e f g" "e d c"

# range unbounded
testp 5: "a b c d e f g" "e f g"
testp 7: "a b c d e f g" "g"
testp :3 "a b c d e f g" "a b c"
testp :1 "a b c d e f g" "a"

# mix
testp 2,5,6 "a b c d e f g" "b e f"
testp 2:4,6: "a b c d e f g" "b c d f g"
testn :-5,-1 "a b c d e f g" "a b c g"


# examples
testp 3 "a b c d e f g" "c"
testn -3 "a b c d e f g" "e"
testp 1,-1,2,-2 "a b c d e f g" "a g b f"
testp 4: "a b c d e f g" "d e f g"
testn -2: "a b c d e f g" "f g"
testp :2 "a b c d e f g" "a b"
testp :-2 "a b c d e f g" "a b c d e f"
testp 2:5 "a b c d e f g" "b c d e"
testp 1,2,-2,3:5,2:,:-3 "a b c d e f g" "a b f c d e b c d e f g a b c d e"

