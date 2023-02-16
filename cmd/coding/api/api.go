package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type CatFact struct {
	Text string `json:"text"`
}

const (
	publicApi = "https://cat-fact.herokuapp.com/facts"
)

func main() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", publicApi, nil)
	if err != nil {
		panic(err)
	}

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	d, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var facts []CatFact
	err = json.Unmarshal(d, &facts)
	if err != nil {
		panic(err)
	}

	for _, f := range facts {
		fmt.Println(f.Text)
		fmt.Println()
	}
}
