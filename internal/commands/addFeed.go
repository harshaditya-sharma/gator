package commands

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/harshaditya-sharma/gator/internal/config"
	"github.com/harshaditya-sharma/gator/internal/database"
)

func HandlerAddFeed(s *config.State, cmd Command, user_id uuid.UUID) error {
	if len(cmd.Args) < 2 {
		return errors.New("insufficient parameters given.")
	}

	name := cmd.Args[0]
	url := cmd.Args[1]
	if name == "" {
		return errors.New("Cannot add feed without a name.\n")
	}
	if url == "" {
		return errors.New("Cannot add feed with empty url.\n")
	}

	if duplicate_feed, err := s.Db.GetFeedFromUrl(context.Background(), url); err == nil {
		creator_name, err2 := s.Db.GetUserFromID(context.Background(), duplicate_feed.UserID)
		if err2 != nil {
			return fmt.Errorf("feed with url %v already created by %v", url, creator_name)
		}
		return fmt.Errorf("feed with url %s created by unknown user already exists\n", url)
	}

	newFeed := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user_id,
	}
	feed, err := s.Db.CreateFeed(context.Background(), newFeed)
	if err != nil {
		return fmt.Errorf("Could not add feed with name %s and url %s\n%w\n", name, url, err)
	}

	fmt.Printf("Created Feed\n")
	printFeed(feed)
	return addFollow(s, feed, user_id)
}
