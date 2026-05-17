package gaana

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	baseAPI   = "https://api.gaana.com/"
	searchAPI = "https://gsearch.gaana.com/vichitih/go/v2/"
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
	options.query = strings.TrimSpace(q)

	params := options.build()
	params["include"] = "track" // search for tracks/songs

	req, err := makeRequest(ctx, http.MethodGet, searchAPI, params)
	if err != nil {
		return SearchSongsResult{}, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return SearchSongsResult{}, err
	}
	defer res.Body.Close()

	apiResponse := new(searchSongsAPIResponse)
	err = json.NewDecoder(res.Body).Decode(apiResponse)
	if err != nil {
		return SearchSongsResult{}, err
	}

	return apiResponse.toResult()
}

func (c *Client) GetSongDetailByTrackId(ctx context.Context, id int) (SongDetail, error) {
	params := map[string]string{
		"track_id": strconv.Itoa(id),
	}
	return c.getSong(ctx, params)
}

func (c *Client) GetSongDetailBySeoKey(ctx context.Context, seoKey string) (SongDetail, error) {
	params := map[string]string{
		"seokey": seoKey,
	}
	return c.getSong(ctx, params)
}

func (c *Client) getSong(ctx context.Context, params map[string]string) (SongDetail, error) {
	// params
	params["request_type"] = "web"
	params["st"] = "hls"
	params["pkc"] = "true"
	params["type"] = "song"
	params["subtype"] = "song_detail"

	// make request
	req, err := makeRequest(ctx, http.MethodGet, baseAPI, params)
	if err != nil {
		return SongDetail{}, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return SongDetail{}, err
	}
	defer res.Body.Close()

	apiResponse := new(songDetailAPIResponse)
	err = json.NewDecoder(res.Body).Decode(apiResponse)
	if err != nil {
		return SongDetail{}, err
	}

	return apiResponse.toSongDetail()
}

func makeRequest(ctx context.Context, requestMethod string, requestURL string, params map[string]string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, requestMethod, requestURL, nil)
	if err != nil {
		return nil, err
	}

	// set headers
	req.Header.Set("devicetype", "GaanaWebsiteApp")
	req.Header.Set("gaanaappversion", "gaanaAndroid-8.60.1")
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.Header.Set("appversion", "V6")

	// query params
	q := &url.Values{}
	for k, v := range params {
		q.Set(k, v)
	}
	req.URL.RawQuery = q.Encode()

	return req, nil
}
