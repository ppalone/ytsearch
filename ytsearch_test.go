package ytsearch

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Search(t *testing.T) {
	t.Run("with default http client", func(t *testing.T) {
		c := NewClient(nil)
		res, err := c.Search(context.Background(), "nocopyrightsounds")

		// no errors
		assert.NoError(t, err)

		// response is not empty
		assert.NotEmpty(t, res.Results, "results are empty")

		// continuation key is not empty
		assert.NotEmpty(t, res.Continuation, "continuation token is empty")
	})

	t.Run("with custom http client", func(t *testing.T) {
		httpClient := &http.Client{
			Timeout: time.Nanosecond * 1,
		}

		// client
		c := NewClient(httpClient)
		res, err := c.Search(context.Background(), "nocopyrightsounds")

		assert.ErrorContains(t, err, context.DeadlineExceeded.Error())
		assert.Empty(t, res)
	})

	t.Run("context with insufficent timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*1)

		// client
		c := NewClient(nil)
		res, err := c.Search(ctx, "nocopyrightsounds")

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
		c := NewClient(nil)
		res, err := c.Search(ctx, "nocopyrightsounds")

		assert.NoError(t, err)
		assert.NotEmpty(t, res.Results)

		// clean up
		t.Cleanup(func() {
			cancel()
		})
	})
}

func Test_SearchNext(t *testing.T) {
	t.Run("valid continuation token", func(t *testing.T) {
		c := NewClient(nil)
		res, err := c.Search(context.Background(), "proximity music")

		assert.NoError(t, err)
		assert.NotEmpty(t, res.Results, "results")
		assert.NotEmpty(t, res.Continuation, "continuation token")

		res, err = c.SearchNext(context.Background(), res.Continuation)

		assert.NoError(t, err)
		assert.NotEmpty(t, res.Results, "next results")
		assert.NotEmpty(t, res.Continuation, "next continuation token")
	})
}
