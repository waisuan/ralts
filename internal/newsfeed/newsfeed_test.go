package newsfeed

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"ralts/internal/config"
	"ralts/internal/dependencies"
	testHelper "ralts/internal/testing"
	"sort"
	"testing"
	"time"
)

var cfg = config.NewConfig(true)

func SeedNewsFeed(count int, deps *dependencies.Dependencies) {
	for i := 0; i < count; i++ {
		author := fmt.Sprintf("author_%d", i)
		title := fmt.Sprintf("title_%d", i)
		description := fmt.Sprintf("description_%d", i)
		publishedAt := time.Now().Add(time.Hour * time.Duration(i))

		err := deps.Storage.Exec(context.Background(), `
    	insert into news_feed (author, title, description, url, published_at)
    	VALUES ($1, $2, $3, $4, $5)
        `, author, title, description, "", publishedAt)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func TestNewNewsFeed(t *testing.T) {
	t.Run("LoadAllArticles", func(t *testing.T) {
		assert := assert.New(t)

		deps := dependencies.NewDependencies(cfg)
		defer deps.Disconnect()

		th := testHelper.TestHelper(cfg)
		defer th()

		count := 10
		SeedNewsFeed(count, deps)

		nf := NewNewsFeed(deps)

		a, err := nf.LoadAllArticles()
		assert.Nil(err)
		assert.Len(a, count)
		assert.NotEmpty(a[0].Author)
		assert.NotEmpty(a[0].Title)
		assert.NotEmpty(a[0].Description)
		assert.NotEmpty(a[0].PublishedAt)

		publishedDates := make([]time.Time, count)
		for _, v := range a {
			publishedDates = append(publishedDates, v.PublishedAt)
		}
		sort.SliceIsSorted(publishedDates, func(i, j int) bool {
			return i < j
		})
	})
}
