package gaana

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Eyevinn/hls-m3u8/m3u8"
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

func (c *Client) GetStream(ctx context.Context, stream Stream) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, stream.URL, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// the response returned is a master playlist
	// and the master playlist has the relative URL to media playlist
	masterPlaylist, playlistType, err := m3u8.DecodeFrom(res.Body, true)
	if err != nil {
		return nil, err
	}

	if playlistType != m3u8.MASTER {
		return nil, fmt.Errorf("playlist is not master")
	}

	master := masterPlaylist.(*m3u8.MasterPlaylist)
	if len(master.Variants) == 0 {
		return nil, fmt.Errorf("no media found in master playlist")
	}

	base, err := url.Parse(stream.URL)
	if err != nil {
		return nil, err
	}

	relative, err := url.Parse(master.Variants[0].URI)
	if err != nil {
		return nil, err
	}

	mediaPlaylistURL := base.ResolveReference(relative)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, mediaPlaylistURL.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// media playlist
	mediaPlaylist, playlistType, err := m3u8.DecodeFrom(resp.Body, true)
	if err != nil {
		return nil, err
	}

	if playlistType != m3u8.MEDIA {
		return nil, fmt.Errorf("playlist is not media")
	}

	// media playlist containing the segments
	media := mediaPlaylist.(*m3u8.MediaPlaylist)

	// reader, writer
	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()

		for _, seg := range media.Segments {
			// skip nil segments
			if seg == nil {
				continue
			}

			segmentURL, err := mediaPlaylistURL.Parse(seg.URI)
			if err != nil {
				pw.CloseWithError(err)
				return
			}

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, segmentURL.String(), nil)
			if err != nil {
				pw.CloseWithError(err)
				return
			}

			resp, err := c.httpClient.Do(req)
			if err != nil {
				pw.CloseWithError(err)
				return
			}

			if resp.StatusCode != http.StatusOK {
				resp.Body.Close()
				pw.CloseWithError(fmt.Errorf("unexpected status code: %d", resp.StatusCode))
				return
			}

			_, err = io.Copy(pw, resp.Body)
			if err != nil {
				resp.Body.Close()
				pw.CloseWithError(err)
				return
			}

			resp.Body.Close()
		}
	}()

	return pr, nil
}

func (c *Client) GetStreamByTrackId(ctx context.Context, id int, opts ...StreamOption) (io.ReadCloser, error) {
	options := defaultStreamOptions()
	for _, opt := range opts {
		opt(options)
	}

	detail, err := c.GetSongDetailByTrackId(ctx, id)
	if err != nil {
		return nil, err
	}

	if len(detail.StreamURLs) == 0 {
		return nil, fmt.Errorf("no stream URLs available")
	}

	var stream Stream
	if options.quality == "" {
		stream = detail.StreamURLs[0]
	} else {
		s, ok := findStream(detail.StreamURLs, options.quality)
		if !ok {
			return nil, fmt.Errorf("stream not found for quality: %s", options.quality)
		}
		stream = s
	}

	return c.GetStream(ctx, stream)
}

func findStream(streams []Stream, quality Quality) (Stream, bool) {
	for _, s := range streams {
		if s.Quality == string(quality) {
			return s, true
		}
	}
	return Stream{}, false
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
