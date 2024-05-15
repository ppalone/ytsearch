package ytsearch

import "context"

const (
	ytSearchInternalBaseURL = "https://www.youtube.com/youtubei/v1/search"
)

// YouTube Search Client
type Client struct{}

func (c *Client) Search(query string) (SearchResponse, error) {
	return SearchResponse{}, nil
}

func (c *Client) SearchWithContext(ctx context.Context, query string) (SearchResponse, error) {
	return SearchResponse{}, nil
}

func (c *Client) Next(key string) (SearchResponse, error) {
	return SearchResponse{}, nil
}

func (c *Client) NextWithContext(ctx context.Context, key string) (SearchResponse, error) {
	return SearchResponse{}, nil
}
