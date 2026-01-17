package commands

import (
	"context"
	"errors"
	"fmt"

	"github.com/harshaditya-sharma/gator/internal/config"
)

func HandlerReset(s *config.State, cmd Command) error {
	if len(cmd.Args) > 0 {
		return errors.New("extra argument found:\nusage: gator reset")
	}
	if err := s.Db.ResetUsers(context.Background()); err != nil {
		return fmt.Errorf("reset users failed: %w", err)
	}
	fmt.Print("succesfully reset users\n")
	return nil
}
