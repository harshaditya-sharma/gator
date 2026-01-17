package commands

import (
	"context"
	"fmt"

	"github.com/harshaditya-sharma/gator/internal/config"
)

func HandlerUsers(s *config.State, cmd Command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("extra argument found:\nusage:gator users")
	}

	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("unable to fetch users: %w", err)
	}

	for _, user := range users {
		fmt.Printf("* %v", user.Name)
		if s.Cfg.Current_user_name == user.Name {
			fmt.Print(" (current) ")
		}
		fmt.Print("\n")
	}
	return nil
}
