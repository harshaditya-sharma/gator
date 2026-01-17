#!/bin/bash
set -e

# Build the app
echo "Building gator..."
go build -o gator main.go

# Reset the DB (to start fresh)
echo "Resetting DB..."
./gator reset

# Register User 1
echo "Registering user1..."
./gator register user1
./gator login user1
./gator addfeed "Lane's Blog" "https://www.wagslane.dev/index.xml"

# Verify feed added
echo "Feeds for user1:"
./gator feeds

# Register User 2
echo "Registering user2..."
./gator register user2
./gator login user2
./gator addfeed "BBC World" "http://feeds.bbci.co.uk/news/world/rss.xml"

# User 2 follows Lane's Blog
echo "User2 following Lane's Blog..."
./gator follow "https://www.wagslane.dev/index.xml"

# Verify following
echo "User2 following:"
./gator following

# Run Aggregation (Fetch posts)
# We run it in background and kill it after 10 seconds
echo "Starting aggregation (5s)..."
./gator agg 1s &
AGG_PID=$!
sleep 5
kill $AGG_PID
echo "Aggregation stopped."

# Browse posts as User 2
./gator login user2
echo "Browsing posts for user2 (limit 2):"
./gator browse 2

# Browse posts as User 1
./gator login user1
echo "Browsing posts for user1 (limit 2):"
./gator browse 2

echo "Test complete!"
