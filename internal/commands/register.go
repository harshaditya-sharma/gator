package commands

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/harshaditya-sharma/gator/internal/config"
	"github.com/harshaditya-sharma/gator/internal/database"
	"github.com/lib/pq"
)

func HandlerRegister(s *config.State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return errors.New("username is required:\nusage: gator register <username>")
	}

	username := cmd.Args[0]
	if _, err := s.Db.GetUser(context.Background(), username); err == nil {
		return fmt.Errorf("user %q already exists", username)
	} else if err != sql.ErrNoRows {
		return fmt.Errorf("checking user: %w", err)
	}

	newUser := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	}
	row, err := s.Db.CreateUser(context.Background(), newUser)
	if err != nil {
		if pe, ok := err.(*pq.Error); ok && pe.Code == "23505" {
			return fmt.Errorf("user %q already exists", username)
		}
		return fmt.Errorf("create user: %w", err)
	}
	if err = s.Cfg.SetUser(username); err != nil {
		return fmt.Errorf("set current user: %w", err)
	}
	fmt.Printf("User created with paremeters:\nuuid: %v\ncreated_at: %v\nupdated_at: %v\nname: %v\n", row.ID, row.CreatedAt, row.UpdatedAt, row.Name)
	return nil
}
