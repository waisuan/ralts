package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"net/http"
	"net/url"
	"ralts/internal/config"
	"ralts/internal/dependencies"
	"ralts/internal/newsfeed"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	cfg := config.NewConfig(false)
	deps := dependencies.NewDependencies(cfg)
	defer deps.Disconnect()

	baseURL, _ := url.Parse("http://api.mediastack.com")
	baseURL.Path += "v1/news"

	params := url.Values{}
	params.Add("access_key", "5b480460652aeaad979240ea1e246461") // TODO: Move to env file
	params.Add("limit", "50")
	baseURL.RawQuery = params.Encode()

	req, _ := http.NewRequest("GET", baseURL.String(), nil)
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	var newsFeed newsfeed.NewsFeed
	body, _ := io.ReadAll(res.Body)
	fmt.Println(string(body))

	if err := json.Unmarshal(body, &newsFeed); err != nil {
		panic(err)
	}

	for _, n := range newsFeed.Data {
		j, _ := json.MarshalIndent(n, "", "\t")
		fmt.Println(string(j))

		//t, err := time.Parse(time.RFC3339, n.PublishedAt)
		//if err != nil {
		//	panic(err)
		//}

		err = deps.Storage.Exec(context.Background(), `
    	insert into news_feed (author, title, description, url, published_at)
    	VALUES ($1, $2, $3, $4, $5)
        `, n.Author, n.Title, n.Description, n.Url, n.PublishedAt)
		if err != nil {
			panic(err)
		}
	}
}
