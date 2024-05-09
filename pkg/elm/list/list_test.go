package list

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCons(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Cons single element", func(t *testing.T) {
		l := []int{2, 3, 4}

		SUT := Cons(1, l)
		asserts.Equal([]int{1, 2, 3, 4}, SUT)
	})

	t.Run("Cons with nil slice", func(t *testing.T) {
		var l []int

		SUT := Cons(1, l)
		asserts.Equal([]int{1}, SUT)
	})
}
