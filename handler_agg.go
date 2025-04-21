package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/TeraSurror/gator/internal/database"
)

func aggHandler(s *state, cmd command) error {
	if len(cmd.Args) < 1 || len(cmd.Args) > 2 {
		return fmt.Errorf("usage: %s <time_between_requests>", cmd.Name)
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return err
	}

	log.Printf("Collecting feeds every %s...\n", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)

	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}

}

func scrapeFeeds(s *state) {
	rssFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Printf("error getting next feed to fetch: %v\n", err)
		return
	}

	log.Printf("found feed to fetch")

	scrapeFeed(s.db, rssFeed)
}

func scrapeFeed(db *database.Queries, rssFeed database.Feed) {
	err := db.MarkFeedFetched(context.Background(), rssFeed.ID)
	if err != nil {
		log.Println("oh bhai heavy ho gaya :(")
		return
	}

	feedData, err := fetchFeed(context.Background(), rssFeed.Url)
	if err != nil {
		log.Printf("error fetching feed data: %v\n", err)
		return
	}

	for _, item := range feedData.Channel.Item {
		log.Printf("Found post: %s\n", item.Title)
	}
}
