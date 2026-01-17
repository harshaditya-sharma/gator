package commands

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/harshaditya-sharma/gator/internal/config"
)

func HandlerFollowing(s *config.State, cmd Command, user_id uuid.UUID) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("Got %v arguments where none expected\n", len(cmd.Args))
	}

	follows, err := s.Db.GetFeedFollowsForUser(context.Background(), user_id)
	if err != nil {
		return fmt.Errorf("Cannot fetch followd feeds.\n%w\n", err)
	}
	fmt.Printf("User follows %v feeds.\n", len(follows))
	for i, follow := range follows {
		feed, err := s.Db.GetFeed(context.Background(), follow.FeedID)
		fmt.Printf("Feed: %v\n", i)
		if err != nil {
			fmt.Printf("Error fetching feed with ID:%v\n", follow.FeedID)
		} else {
			printFeed(feed)
		}
	}
	return nil
}
