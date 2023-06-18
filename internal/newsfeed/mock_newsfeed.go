package newsfeed

import (
	"fmt"
	"time"
)

type MockNewsFeedHandler struct {
	Config *MockNewsFeedConfig
}

type MockNewsFeedConfig struct {
	Seeded bool
}

func NewMockNewsFeedHandler(cfg *MockNewsFeedConfig) *MockNewsFeedHandler {
	return &MockNewsFeedHandler{
		Config: cfg,
	}
}

func (mnf *MockNewsFeedHandler) LoadAllArticles() (Articles, error) {
	if mnf.Config.Seeded == false {
		return []Article{}, nil
	}

	var articles Articles
	for i := 0; i < 10; i++ {
		author := fmt.Sprintf("author_%d", i)
		title := fmt.Sprintf("title_%d", i)
		description := fmt.Sprintf("description_%d", i)
		publishedAt := time.Now().Add(time.Hour * time.Duration(i))

		articles = append(articles, Article{
			Author:      author,
			Title:       title,
			Description: description,
			PublishedAt: publishedAt,
		})
	}

	return articles, nil
}
