package ytsearch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearch(t *testing.T) {
	c := Client{}
	res, err := c.Search("nocopyrightsounds")

	// no errors
	assert.NoError(t, err)

	// response is not empty
	assert.NotEmpty(t, res)
}
