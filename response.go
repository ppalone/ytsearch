package ytsearch

type itemSectionRenderer struct {
	Contents []struct {
		VideoRenderer struct {
			VideoID   string `json:"videoId"`
			Thumbnail struct {
				Thumbnails []Thumbnail `json:"thumbnails"`
			} `json:"thumbnail"`
			Title struct {
				Runs []struct {
					Text string `json:"text"`
				} `json:"runs"`
			} `json:"title"`
			ViewCountText struct {
				SimpleText string `json:"simpleText"`
			} `json:"viewCountText"`
			LengthText struct {
				SimpleText string `json:"simpleText"`
			} `json:"lengthText"`
			OwnerText struct {
				Runs []struct {
					Text string `json:"text"`
				} `json:"runs"`
			} `json:"ownerText"`
		} `json:"videoRenderer"`
	} `json:"contents"`
}

type continuationItemRenderer struct {
	ContinuationEndpoint struct {
		ContinuationCommand struct {
			Token string `json:"token"`
		} `json:"continuationCommand"`
	} `json:"continuationEndpoint"`
}

// Innertube's raw response
type innertubeRawResponse struct {
	Contents struct {
		TwoColumnSearchResultsRenderer struct {
			PrimaryContents struct {
				SectionListRenderer struct {
					Contents []struct {
						ItemSectionRenderer      *itemSectionRenderer      `json:"itemSectionRenderer,omitempty"`
						ContinuationItemRenderer *continuationItemRenderer `json:"continuationItemRenderer,omitempty"`
					} `json:"contents"`
				} `json:"sectionListRenderer"`
			} `json:"primaryContents"`
		} `json:"twoColumnSearchResultsRenderer"`
	} `json:"contents"`
}

// Video info
type VideoInfo struct {
	VideoID    string
	Title      string
	Channel    string
	Thumbnails []Thumbnail
	Views      string
	Duration   string
}

// Thumbnail
type Thumbnail struct {
	URL    string `json:"url"`
	Width  uint   `json:"width"`
	Height uint   `json:"height"`
}

// Search Response
type SearchResponse struct {
	Results      []VideoInfo
	Continuation string
}
