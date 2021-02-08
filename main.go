package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
)

const twitterURL = "https://twitter.com/nutcrackify/status/745937647068155904"

type tweet struct {
	Name     string
	Username string
	Message  string
}

func main() {
	col := colly.NewCollector(
		colly.Debugger(&debug.LogDebugger{}),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.146 Safari/537.36"),
	)

	messages := []tweet{}

	col.OnHTML(".tweet", func(e *colly.HTMLElement) {
		messages = append(messages, tweet{
			Name:     e.ChildText(".account-group .fullname"),
			Username: e.ChildText(".account-group .username"),
			Message:  e.ChildText(".tweet-text"),
		})
	})

	col.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	col.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", string(r.Body), "\nError:", err)
	})

	if err := col.Visit(twitterURL); err != nil {
		log.Panicln(err)
	}

	col.Wait()

	bs, err := json.MarshalIndent(messages, "", "\t")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bs))
	fmt.Println("Number of tweets :", len(messages))
}
