package commands

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/harshaditya-sharma/gator/internal/config"
)

func HandlerLogin(s *config.State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return errors.New("username is required: usage: gator login <username>")
	}

	username := cmd.Args[0]
	user, err := s.Db.GetUser(context.Background(), username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user %q does not exist", username)
		}
		return fmt.Errorf("get user: %w", err)
	}

	if err = s.Cfg.SetUser(user.Name); err != nil {
		return fmt.Errorf("failed to set user: %w", err)
	}

	fmt.Printf("User has been set to %s\n", username)
	return nil
}
