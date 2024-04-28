package parser

import (
	"testing"
)

func TestIsSubString(t *testing.T) {
	t.Run("Offset", func(t *testing.T) {
		SUT, _, _ := IsSubString("hello", 0, 1, 1, "hello world")

		if SUT != 5 {
			t.Error("Incorrect offset. I was expecting 5 and got:", SUT)
		}
	})

	t.Run("Offset - equal length strings", func(t *testing.T) {
		SUT, _, _ := IsSubString("hello", 0, 1, 1, "hello")

		if SUT != 5 {
			t.Error("Incorrect offset. I was expecting 5 and got:", SUT)
		}
	})

	t.Run("Offset - big string smaller than small string", func(t *testing.T) {
		SUT, _, _ := IsSubString("hello", 0, 1, 1, "hell")

		if SUT != -1 {
			t.Error("Incorrect offset. I was expecting -1 and got:", SUT)
		}
	})

	t.Run("Row", func(t *testing.T) {
		_, SUT, _ := IsSubString("hello", 0, 1, 1, "hello world")

		if SUT != 1 {
			t.Error("Incorrect row. I was expecting 1 and got:", SUT)
		}
	})

	t.Run("Row - with new line", func(t *testing.T) {
		_, SUT, _ := IsSubString("he\nllo", 0, 1, 1, "he\nllo world")

		if SUT != 2 {
			t.Error("Incorrect row, should be 2 but got:", SUT)
		}
	})

	t.Run("Column - with new line", func(t *testing.T) {
		_, _, SUT := IsSubString("he\nllo", 0, 1, 1, "he\nllo world")

		if SUT != 4 {
			t.Error("Incorrect column, should be 1 but got:", SUT)
		}
	})

	t.Run("Column", func(t *testing.T) {
		_, _, SUT := IsSubString("hello", 0, 1, 1, "hello world")

		if SUT != 6 {
			t.Error("Incorrect column", SUT)
		}
	})

	t.Run("Next char", func(t *testing.T) {
		src := "hello world"
		offset, _, _ := IsSubString("hello", 0, 1, 1, src)

		SUT := src[offset]

		// space
		if SUT != 32 {
			t.Error("Incorrect offset character. Was expecting 32 and got:", SUT)
		}
	})

	t.Run("Offset with multi byte strings", func(t *testing.T) {
		const sample = "\xe2\x8c\x98\xbd\xb2\x3d\xbc\x20"
		SUT, _, _ := IsSubString("⌘", 0, 1, 1, sample)

		if SUT != 3 {
			t.Error("Incorrect offset character. Was expecting 3 and got:", SUT)
		}
	})
}
