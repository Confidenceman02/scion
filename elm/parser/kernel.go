package parser

import "unicode/utf8"

/*
Validates a sub string against a source string returning an incremented offset, row and column.

ss : small string value

The offset will get incremented by the byte length of each rune,
which allows us to support utf-8 encoding.

The offset is handy for when you want to work with the raw bytes of the string.

offset : Current byte offset

		      The following shows a 3 byte rune at the start of a string and the respective offset values.
	        e2 8c 98 bd b2 3d bc 20
	         ^        ^  ^  ^  ^  ^
	         0        3  4  5  6  7

row = Current line row

The row increments when the parser sees newline unicode values.

You would use the col to surface a location on a row in a file like an editor.

col = Current character column index starting at 1.
*/
func IsSubString(
	ss string,
	offset int,
	row int,
	col int,
	bs string) (int, int, int) {

	smallLength := len(ss)
	isGood := offset+smallLength <= len(bs)

	for i, ssRune := range ss {
		if !isGood {
			break
		}
		// Extract rune from source string
		bsRune, width := utf8.DecodeRuneInString(bs[i:])
		if ssRune == bsRune {
			col++
			offset += width

			// \n
			if ssRune == 0x000A {
				row++
				col = 1
			}
		} else {
			isGood = false
		}
	}
	if !isGood {
		offset = -1
	}
	return offset, row, col
}
