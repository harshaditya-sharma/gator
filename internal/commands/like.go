package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/harshaditya-sharma/gator/internal/config"
	"github.com/harshaditya-sharma/gator/internal/database"
)

func HandlerLike(s *config.State, cmd Command, user_id uuid.UUID) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %v <post_url>", cmd.Name)
	}
	url := cmd.Args[0]

	post, err := s.Db.GetPostByURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("error finding post: %v", err)
	}

	_, err = s.Db.CreatePostLike(context.Background(), database.CreatePostLikeParams{
		UserID:    user_id,
		PostID:    post.ID,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return fmt.Errorf("error liking post: %v", err)
	}

	fmt.Printf("Post liked: %s\n", post.Title)
	return nil
}

func HandlerUnlike(s *config.State, cmd Command, user_id uuid.UUID) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %v <post_url>", cmd.Name)
	}
	url := cmd.Args[0]

	post, err := s.Db.GetPostByURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("error finding post: %v", err)
	}

	err = s.Db.DeletePostLike(context.Background(), database.DeletePostLikeParams{
		UserID: user_id,
		PostID: post.ID,
	})
	if err != nil {
		return fmt.Errorf("error unliking post: %v", err)
	}

	fmt.Printf("Post unliked: %s\n", post.Title)
	return nil
}
