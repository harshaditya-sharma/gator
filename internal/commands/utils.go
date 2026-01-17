package commands

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/harshaditya-sharma/gator/internal/config"
	"github.com/harshaditya-sharma/gator/internal/database"
	"github.com/lib/pq"
)

func printFeed(feed database.Feed) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* UserID:        %s\n", feed.UserID)
}

func isUniqueViolation(e error) bool {
	var pqErr *pq.Error
	if errors.As(e, &pqErr) {
		return pqErr.Code == "23505"
	}
	return false
}

func MiddlewareLoggedIn(handler func(s *config.State, cmd Command, user_id uuid.UUID) error) func(*config.State, Command) error {
	return func(s *config.State, cmd Command) error {
		user_id, err := s.Db.GetUserID(context.Background(), s.Cfg.Current_user_name)
		if err != nil {
			return fmt.Errorf("cannot fetch info for current user\n%w\n", err)
		}
		return handler(s, cmd, user_id)
	}
}
