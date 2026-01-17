#!/bin/bash
CUR_DIR=$(pwd)
cd ~/Documents/gator/sql/schema/ || { echo "Schema directory not found!"; exit 1; }
goose postgres "postgres://postgres:postgres@localhost:5432/gator" down
goose postgres "postgres://postgres:postgres@localhost:5432/gator" up
cd "$CUR_DIR"
echo "Database migrations reset successfully and returned to $CUR_DIR"

