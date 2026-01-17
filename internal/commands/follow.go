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

func HandlerFollow(s *config.State, cmd Command, user_id uuid.UUID) error {
	if len(cmd.Args) > 1 {
		return errors.New("too many arguments\n")
	}
	if len(cmd.Args) == 0 {
		return errors.New("no URL for feed given.\n")
	}

	feed_url := cmd.Args[0]

	feed, err := s.Db.GetFeedFromUrl(context.Background(), feed_url)
	if err != nil {
		return fmt.Errorf("Cannot get feed with the url %v.\n Please add the feed first using addfeed <name> <url>\n", feed_url)
	}
	return addFollow(s, feed, user_id)
}
func addFollow(s *config.State, feed database.Feed, user_id uuid.UUID) error {

	newFollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user_id,
		FeedID:    feed.ID,
	}
	follow, err := s.Db.CreateFeedFollow(context.Background(), newFollow)
	if err != nil {
		if isUniqueViolation(err) {
			return fmt.Errorf("user %v already follows %v", s.Cfg.Current_user_name, feed.Name)
		}
		return fmt.Errorf("error following feed:\n%w\n", err)
	}
	fmt.Printf("User %v succesfully followed feed %v\n", follow.UserName, follow.FeedName)

	return nil
}
