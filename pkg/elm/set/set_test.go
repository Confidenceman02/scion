package set

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestListMethods(t *testing.T) {
	asserts := assert.New(t)

	t.Run("FromList", func(t *testing.T) {
		list := []string{"a", "b", "c"}

		SUT := FromList(list)

		asserts.Equal([]string{"a", "b", "c"}, ToList(SUT))
	})
	t.Run("Creates Set with no dupes", func(t *testing.T) {
		list := []string{"a", "b", "c", "c"}

		SUT := FromList(list)

		asserts.Equal([]string{"a", "b", "c"}, ToList(SUT))
	})
}
