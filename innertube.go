package ytsearch

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

// Default Innertube web client
var innertubeWebClient innertubeClient = innertubeClient{
	Hl:            "en",
	Gl:            "US",
	UserAgent:     "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36,gzip(gfe)",
	ClientName:    "WEB",
	ClientVersion: "2.20240514.03.00",
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
