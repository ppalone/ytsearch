package ytsearch

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Search(t *testing.T) {
	c := Client{}
	res, err := c.Search("nocopyrightsounds")

	// no errors
	assert.NoError(t, err)

	// response is not empty
	assert.NotEmpty(t, res.Results, "results are empty")

	// continuation key is not empty
	assert.NotEmpty(t, res.Continuation, "continuation token is empty")
}

func Test_SearchWithContext(t *testing.T) {
	t.Run("context with insufficent timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*1)

		// client
		c := Client{}
		res, err := c.SearchWithContext(ctx, "nocopyrightsounds")

		assert.ErrorIs(t, err, context.DeadlineExceeded)
		assert.Empty(t, res.Results)

		// clean up
		t.Cleanup(func() {
			cancel()
		})
	})

	t.Run("context with sufficent timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

		// client
		c := Client{}
		res, err := c.SearchWithContext(ctx, "nocopyrightsounds")

		assert.NoError(t, err)
		assert.NotEmpty(t, res.Results)

		// clean up
		t.Cleanup(func() {
			cancel()
		})
	})
}
