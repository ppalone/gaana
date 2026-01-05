package gaana

import (
	"fmt"
	"strings"
)

// Search Songs API Response
type searchSongsAPIResponse struct {
	Gr []struct {
		Gd []struct {
			Id  int    `json:"id"`  // ID
			Ti  string `json:"ti"`  // Title
			Aw  string `json:"aw"`  // Album Art
			Sti string `json:"sti"` // Artists
			Seo string `json:"seo"` // Seo Key
		} `json:"gd"`
	} `json:"gr"`
}

func (res *searchSongsAPIResponse) toResult(opts *searchOptions) (SearchSongsResult, error) {
	if len(res.Gr) == 0 {
		return SearchSongsResult{}, fmt.Errorf("no tracks found")
	}

	songs := make([]Song, 0)
	tracks := res.Gr[0]
	for _, track := range tracks.Gd {
		songs = append(songs, Song{
			Id:       track.Id,
			Title:    track.Ti,
			CoverURL: track.Aw,
			Artists:  strings.Split(track.Sti, ","),
			SeoKey:   track.Seo,
		})
	}

	return SearchSongsResult{
		Page:  opts.page,
		Size:  len(songs),
		Songs: songs,
	}, nil
}
