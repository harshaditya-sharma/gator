package commands

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/harshaditya-sharma/gator/internal/config"
	"github.com/harshaditya-sharma/gator/internal/database"
	"github.com/harshaditya-sharma/gator/internal/structs"
)

func HandlerAgg(s *config.State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %v <time_between_reqs>", cmd.Name)
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}

	fmt.Printf("Collecting feeds every %s...\n", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)

	for ; ; <-ticker.C {
		if err := scrapeFeeds(s); err != nil {
			return err
		}
	}
}

func scrapeFeeds(s *config.State) error {
	feed, err := s.Db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("Error getting next feed to fetch\n%w", err)
	}
	if s.Db.MarkFeedFetched(context.Background(), feed.ID) != nil {
		return fmt.Errorf("Error marking feed %s as fetched.", feed.Name)
	}
	if err := scrapeFeed(s.Db, feed); err != nil {
		return err
	}
	return nil
}

func scrapeFeed(db *database.Queries, feed database.Feed) error {
	err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return fmt.Errorf("Couldn't mark feed %s fetched: %v", feed.Name, err)
	}

	feedData, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("Couldn't collect feed %s: %v", feed.Name, err)
	}
	for _, item := range feedData.Channel.Item {
		fmt.Printf("Found post: %s\n", item.Title)
		pubDate, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			pubDate, err = time.Parse(time.RFC1123, item.PubDate)
			if err != nil {
				fmt.Printf("Could not parse date %v: %v\n", item.PubDate, err)
			}
		}

		description := sql.NullString{
			String: item.Description,
			Valid:  item.Description != "",
		}

		publishedAt := sql.NullTime{
			Time:  pubDate,
			Valid: err == nil,
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: description,
			PublishedAt: publishedAt,
			FeedID:      feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			return fmt.Errorf("Couldn't create post: %v", err)
		}
	}
	fmt.Printf("Feed %s collected, %v posts found", feed.Name, len(feedData.Channel.Item))
	return nil
}

func fetchFeed(ctx context.Context, feedURL string) (*structs.RSSFeed, error) {
	if feedURL == "" {
		return nil, fmt.Errorf("Cannot fetch resource at %s.\nEmpty URL given.", feedURL)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Error fetching response\n%w", err)
	}
	req.Header.Set("User-Agent", "gator")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error fetching response\n%w", err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error parsing reponse\n%w", err)
	}

	var feed structs.RSSFeed
	if err = xml.Unmarshal(b, &feed); err != nil {
		return nil, fmt.Errorf("Error parsing xml\n%w", err)
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for _, item := range feed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
	}

	return &feed, nil
}
