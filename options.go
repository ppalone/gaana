package gaana

import "strconv"

// Search Options.
type searchOptions struct {
	// search query
	query string

	// start index
	// always a multiple of 50 (Example 0, 50, 100 ...)
	// doesn't work as expected, API returns same results
	startIndex int
}

type SearchOption func(opts *searchOptions)

func defaultSearchOptions() *searchOptions {
	return &searchOptions{
		query:      "",
		startIndex: 0,
	}
}

func WithSearchStartIndex(startIndex int) SearchOption {
	return func(opts *searchOptions) {
		opts.startIndex = startIndex
	}
}

func (opts *searchOptions) build() map[string]string {
	m := make(map[string]string)

	m["query"] = opts.query
	m["startIndex"] = strconv.Itoa(opts.startIndex)

	return m
}

type Quality string

const (
	QualityLow    Quality = "low"
	QualityMedium Quality = "medium"
	QualityHigh   Quality = "high"
	QualityAuto   Quality = "auto"
)

type streamOptions struct {
	quality Quality
}

type StreamOption func(opts *streamOptions)

func defaultStreamOptions() *streamOptions {
	return &streamOptions{
		quality: "", // if kept empty, whichever is available will be used
	}
}

func WithStreamQuality(quality Quality) StreamOption {
	return func(opts *streamOptions) {
		opts.quality = quality
	}
}
