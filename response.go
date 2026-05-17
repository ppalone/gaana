package gaana

import (
	"fmt"
	"strings"
	"time"
)

// Search Songs API Response
type searchSongsAPIResponse struct {
	Gr []struct {
		Gd []struct {
			Id   int      `json:"id"`   // ID
			Ti   string   `json:"ti"`   // Title
			Aw   string   `json:"aw"`   // Album Art
			Sti  string   `json:"sti"`  // Artists
			Seo  string   `json:"seo"`  // Seo Key
			Tags []string `json:"tags"` // Tags
		} `json:"gd"`
	} `json:"gr"`
}

func (res *searchSongsAPIResponse) toResult() (SearchSongsResult, error) {
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
			Tags:     track.Tags,
		})
	}

	return SearchSongsResult{
		Size:  len(songs),
		Songs: songs,
	}, nil
}

// Song Detail API Response
type songDetailAPIResponse struct {
	Count  int `json:"count"`
	Tracks []struct {
		Tags []struct {
			TagID   int    `json:"tag_id"` // already an int
			TagName string `json:"tag_name"`
		} `json:"tags"`
		Atw          string `json:"atw"`
		TrackID      int    `json:"track_id,string"`
		Seokey       string `json:"seokey"`
		AlbumSeokey  string `json:"albumseokey"`
		TrackTitle   string `json:"track_title"`
		AlbumID      int    `json:"album_id,string"`
		AlbumTitle   string `json:"album_title"`
		Language     string `json:"language"`
		Artwork      string `json:"artwork"`
		ArtworkWeb   string `json:"artwork_web"`
		ArtworkLarge string `json:"artwork_large"`
		Artist       []struct {
			ArtistID      int    `json:"artist_id,string"`
			Seokey        string `json:"seokey"`
			Name          string `json:"name"`
			Artwork       string `json:"atw"`
			FavoriteCount string `json:"favorite_count"`
		} `json:"artist"`
		Gener []struct {
			GenreID int    `json:"genre_id,string"`
			Name    string `json:"name"`
		} `json:"gener"` // genre is mispelled as "gener"
		LyricsURL           string `json:"lyrics_url"`
		YoutubeID           string `json:"youtube_id"`
		Popularity          string `json:"popularity"`
		Rating              string `json:"rating"`
		TotalFavouriteCount int    `json:"total_favourite_count"`
		StreamType          string `json:"stream_type"`
		Duration            int    `json:"duration,string"` // in seconds
		Isrc                string `json:"isrc"`
		IsMostPopular       int    `json:"is_most_popular"`
		TrackFormat         struct {
			Mp3 struct {
				Normal  string `json:"normal"`
				Medium  string `json:"medium"`
				High    string `json:"high"`
				Extreme string `json:"extreme"`
				Auto    string `json:"auto"`
			} `json:"mp3"`
			Mp4Aac struct {
				Normal  string `json:"normal"`
				Medium  string `json:"medium"`
				High    string `json:"high"`
				Extreme string `json:"extreme"`
				Auto    string `json:"auto"`
			} `json:"mp4_aac"`
		} `json:"track_format"`
		Urls map[string]struct {
			Message    string `json:"message"`
			BitRate    string `json:"bitRate"`
			ExpiryTime int    `json:"expiryTime"`
		} `json:"urls"`
		ReleaseDate string `json:"release_date"`
		Composer    []struct {
			ComposerID int    `json:"composer_id,string"`
			Name       string `json:"name"`
			Seokey     string `json:"seokey"`
		} `json:"composer"`
		PlayCt            string `json:"play_ct"`
		SecondaryLanguage string `json:"secondary_language"`
		ArtistDetail      []struct {
			ArtistID       int    `json:"artist_id,string"`
			Seokey         string `json:"seokey"`
			Name           string `json:"name"`
			Artwork        string `json:"artwork"`
			Artwork175X175 string `json:"artwork_175x175"`
			Atw            string `json:"atw"`
			Role           string `json:"role"`
		} `json:"artist_detail"`
		TotalDownloads int `json:"total_downloads"`
		IsPremium      int `json:"is_premium"`
		PreviewURL     struct {
			Message    string `json:"message"`
			ExpiryTime int    `json:"expiryTime"`
		} `json:"preview_url"`
	} `json:"tracks"`
	Status     int `json:"status"`
	StatusCode int `json:"status_code"`
}

func (res *songDetailAPIResponse) toSongDetail() (SongDetail, error) {
	if res.StatusCode != 0 {
		return SongDetail{}, fmt.Errorf("song details not found")
	}

	if len(res.Tracks) == 0 {
		return SongDetail{}, fmt.Errorf("song details not found")
	}

	data := res.Tracks[0]

	artists := make([]Artist, 0)
	for _, a := range data.Artist {
		artists = append(artists, Artist{
			Id:            a.ArtistID,
			Seokey:        a.Seokey,
			Name:          a.Name,
			Artwork:       a.Artwork,
			FavoriteCount: a.FavoriteCount,
		})
	}

	genres := make([]Genre, 0)
	for _, g := range data.Gener {
		genres = append(genres, Genre{
			Id:   g.GenreID,
			Name: g.Name,
		})
	}

	tags := make([]Tag, 0)
	for _, t := range data.Tags {
		tags = append(tags, Tag{
			Id:   t.TagID,
			Name: t.TagName,
		})
	}

	streamURLs := make([]Stream, 0)
	for quality, data := range data.Urls {
		url, err := decAESCBCPKCS(data.Message)
		if err != nil {
			continue
		}
		streamURLs = append(streamURLs, Stream{
			Quality:    quality,
			URL:        url,
			Bitrate:    data.BitRate,
			ExpiryTime: time.Unix(int64(data.ExpiryTime), 0),
		})
	}

	// ignore error
	releaseDate, _ := time.Parse("2006-01-02", data.ReleaseDate)

	return SongDetail{
		Id:           data.TrackID,
		Title:        data.TrackTitle,
		SeoKey:       data.Seokey,
		Duration:     time.Duration(data.Duration) * time.Second,
		Language:     data.Language,
		ReleaseDate:  releaseDate,
		Atw:          data.Atw,
		Artwork:      data.Artwork,
		ArtworkWeb:   data.ArtworkWeb,
		ArtworkLarge: data.ArtworkLarge,
		AlbumId:      data.AlbumID,
		AlbumTitle:   data.AlbumTitle,
		Artists:      artists,
		Genres:       genres,
		Tags:         tags,
		StreamURLs:   streamURLs,
	}, nil
}
