package main

import (
	"context"
	"log"
)

func aggHandler(s *state, cmd command) error {
	rssFeed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	log.Printf("%+v\n", rssFeed)

	return nil
}
