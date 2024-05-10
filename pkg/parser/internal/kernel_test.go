package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsSubString(t *testing.T) {
	asserts := assert.New(t)
	t.Run("Offset", func(t *testing.T) {
		SUT, _, _ := IsSubString("hello", 0, 1, 1, "hello world")

		asserts.Equal(5, SUT)
	})

	t.Run("Offset - equal length strings", func(t *testing.T) {
		SUT, _, _ := IsSubString("hello", 0, 1, 1, "hello")

		asserts.Equal(5, SUT)
	})

	t.Run("Offset - big string smaller than small string", func(t *testing.T) {
		SUT, _, _ := IsSubString("hello", 0, 1, 1, "hell")

		asserts.Equal(-1, SUT)
	})

	t.Run("Row", func(t *testing.T) {
		_, SUT, _ := IsSubString("hello", 0, 1, 1, "hello world")

		asserts.Equal(1, SUT)
	})

	t.Run("Row - with new line", func(t *testing.T) {
		_, SUT, _ := IsSubString("he\nllo", 0, 1, 1, "he\nllo world")

		asserts.Equal(2, SUT)
	})

	t.Run("Column - with new line", func(t *testing.T) {
		_, _, SUT := IsSubString("he\nllo", 0, 1, 1, "he\nllo world")

		asserts.Equal(4, SUT)
	})

	t.Run("Column", func(t *testing.T) {
		_, _, SUT := IsSubString("hello", 0, 1, 1, "hello world")

		asserts.Equal(6, SUT)
	})

	t.Run("Next char", func(t *testing.T) {
		src := "hello world"
		offset, _, _ := IsSubString("hello", 0, 1, 1, src)

		SUT := src[offset]

		asserts.Equal(" ", string(SUT))
	})

	t.Run("Offset with multi byte runes", func(t *testing.T) {
		const sample = "\xe2\x8c\x98\xbd\xb2\x3d\xbc\x20"
		SUT, _, _ := IsSubString("⌘", 0, 1, 1, sample)

		asserts.Equal(3, SUT)
	})

	t.Run("Column with multi byte runes", func(t *testing.T) {
		const sample = "\xe2\x8c\x98\xbd\xb2\x3d\xbc\x20"
		_, _, SUT := IsSubString("⌘", 0, 1, 1, sample)

		asserts.Equal(2, SUT)
	})
}

func TestIsSubChar(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Matches rune", func(t *testing.T) {
		SUT := IsSubChar(func(char int32) bool { return char == 'h' }, 0, "hello")

		asserts.Equal(1, SUT)
	})
}
