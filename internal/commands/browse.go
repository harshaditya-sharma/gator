package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/harshaditya-sharma/gator/internal/config"
	"github.com/harshaditya-sharma/gator/internal/database"
)

func HandlerBrowse(s *config.State, cmd Command, user_id uuid.UUID) error {
	if len(cmd.Args) > 1 {
		return fmt.Errorf("usage: %v <number_of_posts>", cmd.Name)
	}
	postLimit := 2
	if len(cmd.Args) == 1 {
		limit, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("usage: %v <number_of_posts>", cmd.Name)
		}
		postLimit = limit
	}

	posts, err := s.Db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user_id,
		Limit:  int32(postLimit),
	})
	if err != nil {
		return fmt.Errorf("Error getting posts: %v", err)
	}

	for _, post := range posts {
		fmt.Printf("%s\n", post.Title)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Printf("%s\n", post.Description.String)
		fmt.Printf("Published: %s\n", post.PublishedAt.Time.Format("Mon Jan 2 2006"))
		fmt.Println("=====================================")
	}
	return nil
}
