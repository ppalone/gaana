package gaana

import "time"

type Stream struct {
	Quality    string
	URL        string
	Bitrate    string // can be ""
	ExpiryTime time.Time
}
