package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/TeraSurror/gator/internal/database"
	"github.com/google/uuid"
)

func addFeedHandler(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
	}
	name := cmd.Args[0]
	url := cmd.Args[1]

	feed := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	}

	createdFeed, err := s.db.CreateFeed(context.Background(), feed)
	if err != nil {
		log.Printf("could not create feed: %v", err)
		return err
	}

	feedFollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    createdFeed.UserID,
		FeedID:    createdFeed.ID,
	}
	_, err = s.db.CreateFeedFollow(context.Background(), feedFollow)
	if err != nil {
		return fmt.Errorf("could not created feed follow: %v", err)
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

func followHandler(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <feed url>", cmd.Name)
	}

	url := cmd.Args[0]

	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("could not fetch feed: %v", err)
	}

	feedFollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	feedFollowRow, err := s.db.CreateFeedFollow(context.Background(), feedFollow)
	if err != nil {
		return fmt.Errorf("could perform follow action: %v", err)
	}

	log.Printf("User: %s followed feed: %s.", feedFollowRow.UserName, feedFollowRow.FeedName)

	return nil
}

func unfollowHandler(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <feed url>", cmd.Name)
	}

	url := cmd.Args[0]

	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("could not fetch feed: %v", err)
	}

	err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{UserID: user.ID, FeedID: feed.ID})
	if err != nil {
		return fmt.Errorf("could unfollow feed - %s: %v", feed.Url, err)
	}

	log.Printf("User: %s unfollowed feed: %s.", user.Name, feed.Url)

	return nil
}

func followsHandler(s *state, cmd command, user database.User) error {
	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		return fmt.Errorf("could not fetch follows for user: %s - %v", user.Name, err)
	}

	log.Printf("User %s follows:\n", user.Name)
	for _, feedFollow := range feedFollows {
		log.Printf("%s\n", feedFollow.FeedName)
	}

	return nil
}
