package gaana

import "strconv"

// Search Options.
type searchOptions struct {
	// search query
	keyword string

	// page will always start from 0
	page int

	// default to "IN" for now
	country string
}

type SearchOption func(opts *searchOptions)

func defaultSearchOptions() *searchOptions {
	return &searchOptions{
		keyword: "",
		page:    0,
		country: "IN",
	}
}

func WithSearchPage(page int) SearchOption {
	return func(opts *searchOptions) {
		opts.page = page
	}
}

func (opts *searchOptions) build() map[string]string {
	m := make(map[string]string)

	m["keyword"] = opts.keyword
	m["page"] = strconv.Itoa(opts.page)
	m["country"] = opts.country

	// "type" wil be same for song/artist/album search
	m["type"] = "search"

	return m
}
