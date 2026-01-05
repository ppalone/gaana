package gaana

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

const (
	apiURL = "https://gaana.com/apiv2"
)

// Gaana Client.
type Client struct {
	httpClient *http.Client
}

// NewClient returns a new Gaana client
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	return &Client{httpClient}
}

// SearchSongs returns search songs results based on the search query.
// Even if we pass option with page (for example 2) it still returns the same results as before.
// This might be a bug from Gaana's end.
func (c *Client) SearchSongs(ctx context.Context, q string, opts ...SearchOption) (SearchSongsResult, error) {
	options := defaultSearchOptions()
	for _, opt := range opts {
		opt(options)
	}
	options.keyword = strings.TrimSpace(q)

	params := options.build()
	params["secType"] = "track" // search for tracks/songs

	req, err := makeRequest(ctx, params)
	if err != nil {
		return SearchSongsResult{}, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return SearchSongsResult{}, err
	}
	defer res.Body.Close()

	apiResponse := new(searchSongsAPIResponse)
	err = json.NewDecoder(res.Body).Decode(&apiResponse)
	if err != nil {
		return SearchSongsResult{}, err
	}

	return apiResponse.toResult(options)
}

func makeRequest(ctx context.Context, params map[string]string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, nil)
	if err != nil {
		return nil, err
	}

	// query params
	q := &url.Values{}
	for k, v := range params {
		q.Set(k, v)
	}
	req.URL.RawQuery = q.Encode()

	return req, nil
}
