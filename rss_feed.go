package main

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"log"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		log.Printf("could not create request to %s: %v\n", feedURL, err)
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")

	response, err := httpClient.Do(req)
	if err != nil {
		log.Printf("could not fetch content from %s: %v\n", feedURL, err)
		return nil, err
	}
	defer response.Body.Close()

	val, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("could not read request body: %v\n", err)
		return nil, err
	}

	var rssData RSSFeed
	err = xml.Unmarshal(val, &rssData)
	if err != nil {
		log.Printf("could not parse response to xml: %v\n", err)
		return nil, err
	}

	rssData.Channel.Title = html.UnescapeString(rssData.Channel.Title)
	rssData.Channel.Description = html.UnescapeString(rssData.Channel.Description)

	for i, item := range rssData.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		rssData.Channel.Item[i] = item
	}

	return &rssData, nil
}
