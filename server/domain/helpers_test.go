package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testObject struct {
	ID string
}

func TestContain_ReturnsTrueOnContains(t *testing.T) {
	t.Parallel()
	// Arrange
	is := testObject{ID: "abc"}
	in := []testObject{{ID: "a"}, {ID: "b"}, {ID: "abc"}}

	// Act
	_, result := containsWithKey(is, in, func(t testObject) string {
		return t.ID
	})

	// Assert
	assert.True(t, result)
}

func TestContain_ReturnsFalseOnNotContains(t *testing.T) {
	t.Parallel()
	// Arrange
	is := testObject{ID: "abc"}
	in := []testObject{{ID: "a"}, {ID: "b"}, {ID: "c"}}

	// Act
	_, result := containsWithKey(is, in, func(t testObject) string {
		return t.ID
	})

	// Assert
	assert.False(t, result)
}
