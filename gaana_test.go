package gaana_test

import (
	"context"
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
	})

	t.Run("with page option", func(t *testing.T) {
		// since the API returns the same results even after passing the page option
		// this appears to bug from the Gaana's API itself
		// skipping the test for now
		t.Skip()

		q := "martin garrix"
		res, err := c.SearchSongs(context.Background(), q, gaana.WithSearchPage(1))
		assert.NoError(t, err)
		assert.NotEmpty(t, res.Songs)
	})
}
