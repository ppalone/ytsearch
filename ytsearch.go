package ytsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

const (
	ytSearchInternalBaseURL = "https://www.youtube.com/youtubei/v1/search"
	ytDefaultVideoParams    = "EgIQAQ%3D%3D"
)

// YouTube Search Client
type Client struct {
	httpClient *http.Client
}

// NewClient returns a new YtSearch client.
func NewClient(c *http.Client) *Client {
	if c == nil {
		c = &http.Client{}
	}

	return &Client{c}
}

// Search
func (c *Client) Search(ctx context.Context, query string) (SearchResponse, error) {
	return c.searchQuery(ctx, query)
}

// SearchNext
func (c *Client) SearchNext(ctx context.Context, key string) (SearchResponse, error) {
	return c.searchNext(ctx, key)
}

func makeRequest(ctx context.Context, method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func extractVideos(items *itemSectionRenderer) []VideoInfo {
	videos := []VideoInfo{}

	for _, item := range items.Contents {
		v := item.VideoRenderer

		// Check if VideoID is present
		// we don't want other stuff
		if len(v.VideoID) == 0 {
			continue
		}

		video := VideoInfo{
			VideoID:    v.VideoID,
			Thumbnails: v.Thumbnail.Thumbnails,
			Views:      v.ViewCountText.SimpleText,
			Duration:   v.LengthText.SimpleText,
		}

		if len(v.Title.Runs) > 0 {
			video.Title = v.Title.Runs[0].Text
		}

		if len(v.OwnerText.Runs) > 0 {
			video.Channel = v.OwnerText.Runs[0].Text
		}

		videos = append(videos, video)
	}

	return videos
}

func extractContinuationToken(item *continuationItemRenderer) string {
	return item.ContinuationEndpoint.ContinuationCommand.Token
}

func (c *Client) searchQuery(ctx context.Context, query string) (SearchResponse, error) {
	// prepare innertube request data
	d := prepareInnertubeRequestForSearch(query)

	reqData, err := json.Marshal(d)
	if err != nil {
		return SearchResponse{}, err
	}

	req, err := makeRequest(ctx, http.MethodPost, ytSearchInternalBaseURL, bytes.NewReader(reqData))
	if err != nil {
		return SearchResponse{}, err
	}

	// make http call
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return SearchResponse{}, err
	}
	defer resp.Body.Close()

	rawResponse := new(innertubeRawResponse)

	err = json.NewDecoder(resp.Body).Decode(rawResponse)
	if err != nil {
		return SearchResponse{}, err
	}

	// extract contents
	contents := rawResponse.Contents.TwoColumnSearchResultsRenderer.PrimaryContents.SectionListRenderer.Contents

	searchResponse := SearchResponse{}

	// iterate over contents
	for _, content := range contents {
		if content.ItemSectionRenderer != nil {
			searchResponse.Results = extractVideos(content.ItemSectionRenderer)
		} else if content.ContinuationItemRenderer != nil {
			searchResponse.Continuation = extractContinuationToken(content.ContinuationItemRenderer)
		}
	}

	return searchResponse, nil
}

func (c *Client) searchNext(ctx context.Context, key string) (SearchResponse, error) {
	// prepare innertube request data
	d := prepareInnertubeRequestForNext(key)

	reqData, err := json.Marshal(d)
	if err != nil {
		return SearchResponse{}, err
	}

	req, err := makeRequest(ctx, http.MethodPost, ytSearchInternalBaseURL, bytes.NewReader(reqData))
	if err != nil {
		return SearchResponse{}, err
	}

	// make http call
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return SearchResponse{}, err
	}
	defer resp.Body.Close()

	rawResponse := new(innertubeRawResponse)

	err = json.NewDecoder(resp.Body).Decode(rawResponse)
	if err != nil {
		return SearchResponse{}, err
	}

	// extract contents
	t := rawResponse.OnResponseReceivedCommands

	searchResponse := SearchResponse{}

	if len(t) > 0 {
		contents := t[0].AppendContinuationItemsAction.ContinuationItems

		// iterate over contents
		for _, content := range contents {
			if content.ItemSectionRenderer != nil {
				searchResponse.Results = extractVideos(content.ItemSectionRenderer)
			} else if content.ContinuationItemRenderer != nil {
				searchResponse.Continuation = extractContinuationToken(content.ContinuationItemRenderer)
			}
		}
	}

	return searchResponse, nil
}
