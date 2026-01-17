package commands

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/harshaditya-sharma/gator/internal/config"
	"github.com/harshaditya-sharma/gator/internal/database"
)

func HandlerSearch(s *config.State, cmd Command, user_id uuid.UUID) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("usage: %v <search_term>", cmd.Name)
	}

	searchTerm := cmd.Args[0]
	limit := int32(10) // Default limit

	params := database.GetPostsForUserMatchingParams{
		UserID:     user_id,
		SearchTerm: sql.NullString{String: searchTerm, Valid: true},
		Limit:      limit,
	}

	posts, err := s.Db.GetPostsForUserMatching(context.Background(), params)
	if err != nil {
		return fmt.Errorf("error searching posts: %v", err)
	}

	if len(posts) == 0 {
		fmt.Println("No posts found matching your query.")
		return nil
	}

	printPosts(posts)
	return nil
}
