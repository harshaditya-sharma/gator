package commands

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/harshaditya-sharma/gator/internal/config"
	"github.com/harshaditya-sharma/gator/internal/database"
)

func HandlerUnfollow(s *config.State, cmd Command, user_id uuid.UUID) error {
	if len(cmd.Args) > 1 || len(cmd.Args) == 0 {
		return fmt.Errorf("expected 1 Argument got %v\n", len(cmd.Args))
	}
	url := cmd.Args[0]
	feed, err := s.Db.GetFeedFromUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("Cannot find feed with url %v\n", url)
	}

	arg := database.DeleteFeedFollowForUserParams{
		UserID: user_id,
		FeedID: feed.ID,
	}
	if err := s.Db.DeleteFeedFollowForUser(context.Background(), arg); err != nil {
		return fmt.Errorf("failed to unfollow %v (%v)\n", feed.Name, feed.Url)
	}
	fmt.Printf("Successfully unfollowed %v (%v)\n", feed.Name, feed.Url)
	return nil
}
