package gaana

import "time"

// Song
type Song struct {
	Id       int
	Title    string
	CoverURL string
	Artists  []string
	SeoKey   string
	Tags     []string
}

// Song Detail
type SongDetail struct {
	// Track details
	Id          int
	Title       string
	SeoKey      string
	Duration    time.Duration
	Language    string
	ReleaseDate time.Time

	// Artwork
	Atw          string
	Artwork      string
	ArtworkWeb   string
	ArtworkLarge string

	// Album details
	AlbumId    int
	AlbumTitle string

	// Artist details
	Artists []Artist

	// TODO: Composer details
	// TODO: Artist + Composer details

	// Streaming details
	StreamURLs []Stream
	PreviewURL Stream

	// Genre
	Genres []Genre

	// Tags
	Tags []Tag
}
