package gaana_test

import (
	"context"
	"strings"
	"testing"

	"github.com/ppalone/gaana"
	"github.com/stretchr/testify/assert"
)

func Test_SearchSongs(t *testing.T) {
	c := gaana.NewClient(nil)

	t.Run("with no options", func(t *testing.T) {
		q := "martin garrix"
		res, err := c.SearchSongs(context.Background(), q)
		assert.NoError(t, err)
		assert.NotEmpty(t, res.Songs)
		assert.Equal(t, 0, res.Page)

		for _, track := range res.Songs {
			assert.NotEmpty(t, track.Id)
			assert.NotEmpty(t, track.Title)
			assert.NotEmpty(t, track.Artists)
		}
	})

	t.Run("with start index option", func(t *testing.T) {
		// since the API returns the same results even after passing the start index option
		// this appears to bug from the Gaana's API itself
		// skipping the test for now
		t.Skip()

		q := "martin garrix"
		res, err := c.SearchSongs(context.Background(), q, gaana.WithSearchStartIndex(50))
		assert.NoError(t, err)
		assert.NotEmpty(t, res.Songs)
	})
}

func Test_GetSongDetailByTrackId(t *testing.T) {
	c := gaana.NewClient(nil)

	t.Run("with valid track id", func(t *testing.T) {
		id := 1783362 // "Animals" by Martin Garrix
		detail, err := c.GetSongDetailByTrackId(context.Background(), id)
		assert.NoError(t, err)
		assert.NotEmpty(t, detail)

		assert.NotEmpty(t, detail.Id)
		assert.NotEmpty(t, detail.Title)
		assert.NotEmpty(t, detail.SeoKey)
		assert.NotEmpty(t, detail.AlbumId)
		assert.NotEmpty(t, detail.Artists)
		assert.NotEmpty(t, detail.Genres)
		assert.NotEmpty(t, detail.Tags)
		assert.NotEmpty(t, detail.ReleaseDate)
		assert.NotEmpty(t, detail.Duration)
		assert.NotEmpty(t, detail.StreamURLs)
	})

	t.Run("with invalid track id", func(t *testing.T) {
		id := 999999999 // invalid track id
		detail, err := c.GetSongDetailByTrackId(context.Background(), id)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "song details not found")
		assert.Empty(t, detail)
	})
}

func Test_GetSongDetailBySeoKey(t *testing.T) {
	c := gaana.NewClient(nil)

	t.Run("with valid seo key", func(t *testing.T) {
		seoKey := "animals-11" // "Animals" by Martin Garrix
		detail, err := c.GetSongDetailBySeoKey(context.Background(), seoKey)
		assert.NoError(t, err)
		assert.NotEmpty(t, detail)

		assert.NotEmpty(t, detail.Id)
		assert.NotEmpty(t, detail.Title)
		assert.NotEmpty(t, detail.SeoKey)
		assert.NotEmpty(t, detail.AlbumId)
		assert.NotEmpty(t, detail.Artists)
		assert.NotEmpty(t, detail.Genres)
		assert.NotEmpty(t, detail.Tags)
		assert.NotEmpty(t, detail.ReleaseDate)
		assert.NotEmpty(t, detail.Duration)
		assert.NotEmpty(t, detail.StreamURLs)
	})

	t.Run("with invalid seo key", func(t *testing.T) {
		seoKey := "invalid-seo-key" // invalid seo key
		detail, err := c.GetSongDetailBySeoKey(context.Background(), seoKey)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "song details not found")
		assert.Empty(t, detail)
	})
}

func Test_CompareSongDetail(t *testing.T) {
	c := gaana.NewClient(nil)

	t.Run("compare song detail by track id and seo key", func(t *testing.T) {
		id := 1783362          // "Animals" by Martin Garrix
		seoKey := "animals-11" // "Animals" by Martin Garrix

		detailById, err := c.GetSongDetailByTrackId(context.Background(), id)
		assert.NoError(t, err)
		assert.NotEmpty(t, detailById)

		detailBySeoKey, err := c.GetSongDetailBySeoKey(context.Background(), seoKey)
		assert.NoError(t, err)
		assert.NotEmpty(t, detailBySeoKey)

		assert.Equal(t, detailById.Id, detailBySeoKey.Id)
		assert.Equal(t, detailById.Title, detailBySeoKey.Title)
		assert.Equal(t, detailById.SeoKey, detailBySeoKey.SeoKey)
		assert.Equal(t, detailById.AlbumId, detailBySeoKey.AlbumId)
		assert.Equal(t, detailById.Artists, detailBySeoKey.Artists)
		assert.Equal(t, detailById.Genres, detailBySeoKey.Genres)
		assert.Equal(t, detailById.Tags, detailBySeoKey.Tags)
		assert.Equal(t, detailById.ReleaseDate, detailBySeoKey.ReleaseDate)
		assert.Equal(t, detailById.Duration, detailBySeoKey.Duration)
		assert.ElementsMatch(t, detailById.StreamURLs, detailBySeoKey.StreamURLs)
	})
}

func Test_StreamURLs(t *testing.T) {
	c := gaana.NewClient(nil)

	ids := []int{
		3200230,
		59685660,
		73865153,
	}
	for _, id := range ids {
		detail, err := c.GetSongDetailByTrackId(context.Background(), id)
		assert.NoError(t, err)
		assert.NotEmpty(t, detail)

		for _, stream := range detail.StreamURLs {
			assert.NotEmpty(t, stream.Quality)
			assert.NotEmpty(t, stream.URL)
			assert.True(t, strings.HasPrefix(stream.URL, "https://"))
			assert.True(t, strings.Contains(stream.URL, ".m3u8"))
		}
	}
}
