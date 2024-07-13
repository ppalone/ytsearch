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
		c := Client{}
		res, err := c.Search("nocopyrightsounds")

		// no errors
		assert.NoError(t, err)

		// response is not empty
		assert.NotEmpty(t, res.Results, "results are empty")

		// continuation key is not empty
		assert.NotEmpty(t, res.Continuation, "continuation token is empty")
	})

	t.Run("with custom http client", func(t *testing.T) {
		httpclient := &http.Client{
			Timeout: time.Nanosecond * 1,
		}

		// client
		c := Client{
			HTTPClient: httpclient,
		}
		res, err := c.Search("nocopyrightsounds")

		assert.ErrorContains(t, err, context.DeadlineExceeded.Error())
		assert.Empty(t, res)
	})
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

func Test_Next(t *testing.T) {
	t.Run("valid continuation token", func(t *testing.T) {
		c := Client{}
		res, err := c.Search("hacker")

		assert.NoError(t, err)
		assert.NotEmpty(t, res.Results, "results")
		assert.NotEmpty(t, res.Continuation, "continuation token")

		// debug
		t.Log("results length:", len(res.Results))
		t.Log("continuation token:", res.Continuation)

		res, err = c.Next(res.Continuation)

		assert.NoError(t, err)
		assert.NotEmpty(t, res.Results, "next results")
		assert.NotEmpty(t, res.Continuation, "next continuation token")

		// debug
		t.Log("next results length:", len(res.Results))
		t.Log("next continuation token:", res.Continuation)
	})
}
