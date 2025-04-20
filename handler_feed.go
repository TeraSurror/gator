package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/TeraSurror/gator/internal/database"
	"github.com/google/uuid"
)

func addFeedHandler(s *state, cmd command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
	}
	name := cmd.Args[0]
	url := cmd.Args[1]

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		log.Printf("could not fetch user info: %v", err)
		return err
	}

	feed := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	}

	_, err = s.db.CreateFeed(context.Background(), feed)
	if err != nil {
		log.Printf("could not create feed: %v", err)
		return err
	}

	log.Println("Feed created successfully!")

	return nil
}

func feedHandler(s *state, cmd command) error {
	feedList, err := s.db.GetFeedList(context.Background())
	if err != nil {
		return fmt.Errorf("could not fetch feed: %v", err)
	}

	for _, feedRow := range feedList {
		log.Printf("%s %s %s\n", feedRow.FeedName, feedRow.Url, feedRow.CreatorName)
	}

	return nil
}
