package config

import (
	"github.com/harshaditya-sharma/gator/internal/database"
)

type State struct {
	Db  *database.Queries
	Cfg *Config
}
