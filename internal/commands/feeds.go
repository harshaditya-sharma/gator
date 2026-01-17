package commands

import (
	"context"
	"fmt"

	"github.com/harshaditya-sharma/gator/internal/config"
)

func HandlerFeeds(s *config.State, cmd Command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("got %v arguments where none expected.\n", len(cmd.Args))
	}
	feeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Cannot get feed details from datbase\n%w\n", err)
	}

	for i, feed := range feeds {
		fmt.Printf("Feed %d\nname: %v\nurl %s\n", i, feed.Name, feed.Url)
		creator, err := s.Db.GetUserFromID(context.Background(), feed.UserID)
		if err != nil {
			fmt.Printf("creator: Unknown User\n")
		} else {
			fmt.Printf("creator: %v\n", creator.Name)
		}
	}

	return nil
}
