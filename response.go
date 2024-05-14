package ytsearch

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
	URL    string
	Width  uint
	Height uint
}

// Search Response
type SearchResponse struct {
	Results      []VideoInfo
	Continuation string
}
