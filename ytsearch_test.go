package ytsearch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Search(t *testing.T) {
	c := Client{}
	res, err := c.Search("nocopyrightsounds")

	// no errors
	assert.NoError(t, err)

	// response is not empty
	assert.NotEmpty(t, res.Results)

	// continuation key is not empty
	assert.NotEmpty(t, res.Continuation)
}
