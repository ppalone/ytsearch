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
	HTTPClient *http.Client
}

// Innertube Client
type innertubeClient struct {
	Hl            string `json:"hl"`
	Gl            string `json:"gl"`
	UserAgent     string `json:"userAgent"`
	ClientName    string `json:"clientName"`
	ClientVersion string `json:"clientVersion"`
}

// Innertube request context
type innertubeRequestContext struct {
	Client innertubeClient `json:"client"`
}

// Innertube Request
type innertubeRequest struct {
	Context      innertubeRequestContext `json:"context"`
	Query        string                  `json:"query,omitempty"`
	Params       string                  `json:"params,omitempty"`
	Continuation string                  `json:"continuation,omitempty"`
}

// Innertube web client
var innertubeWebClient innertubeClient = innertubeClient{
	Hl:            "en",
	Gl:            "US",
	UserAgent:     "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36,gzip(gfe)",
	ClientName:    "WEB",
	ClientVersion: "2.20240514.03.00",
}

func makeRequest(ctx context.Context, method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func prepareInnertubeRequestForSearch(query string) innertubeRequest {
	return innertubeRequest{
		Context: innertubeRequestContext{
			Client: innertubeWebClient,
		},
		Query:  query,
		Params: ytDefaultVideoParams,
	}
}

func prepareInnertubeRequestForNext(key string) innertubeRequest {
	return innertubeRequest{
		Context: innertubeRequestContext{
			Client: innertubeWebClient,
		},
		Continuation: key,
		Params:       ytDefaultVideoParams,
	}
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
	if c.HTTPClient == nil {
		c.HTTPClient = http.DefaultClient
	}

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
	resp, err := c.HTTPClient.Do(req)
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
	if c.HTTPClient == nil {
		c.HTTPClient = http.DefaultClient
	}

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
	resp, err := c.HTTPClient.Do(req)
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

func (c *Client) Search(query string) (SearchResponse, error) {
	return c.SearchWithContext(context.Background(), query)
}

func (c *Client) SearchWithContext(ctx context.Context, query string) (SearchResponse, error) {
	return c.searchQuery(ctx, query)
}

func (c *Client) Next(key string) (SearchResponse, error) {
	return c.NextWithContext(context.Background(), key)
}

func (c *Client) NextWithContext(ctx context.Context, key string) (SearchResponse, error) {
	return c.searchNext(ctx, key)
}
