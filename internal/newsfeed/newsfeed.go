package newsfeed

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"ralts/internal/dependencies"
	"time"
)

type Article struct {
	Author      string    `json:"author"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Url         string    `json:"url"`
	PublishedAt time.Time `json:"published_at"`
}

type Articles []Article

type NewsFeedHandler interface {
	LoadAllArticles() (Articles, error)
}

type NewsFeed struct {
	deps *dependencies.Dependencies
	Data Articles `json:"data"`
}

func NewNewsFeed(deps *dependencies.Dependencies) *NewsFeed {
	return &NewsFeed{
		deps: deps,
	}
}

func (nf *NewsFeed) Print() {
	for _, n := range nf.Data {
		j, _ := json.MarshalIndent(n, "", "\t")
		fmt.Println(string(j))
	}
}

func (nf *NewsFeed) LoadAllArticles() (Articles, error) {
	rows, err := nf.deps.Storage.Query(context.Background(), `
		select author, title, description, url, published_at
		from news_feed 
		order by published_at;
`)
	if err != nil {
		log.Error(fmt.Sprintf("Unable to execute query -> %s", err.Error()))
		return nil, err
	}
	defer rows.Close()

	var articles Articles
	for rows.Next() {
		var a Article
		err := rows.Scan(&a.Author, &a.Title, &a.Description, &a.Url, &a.PublishedAt)
		if err != nil {
			log.Error(fmt.Sprintf("Unable load query results -> %s", err.Error()))
			return nil, err
		}
		articles = append(articles, a)
	}
	if err := rows.Err(); err != nil {
		log.Error(fmt.Sprintf("Something went wrong -> %s", err.Error()))
		return nil, err
	}

	return articles, nil
}
