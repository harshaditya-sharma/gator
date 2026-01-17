package commands

import (
	"context"
	"fmt"
	"reflect"
	"strconv"

	"github.com/google/uuid"
	"github.com/harshaditya-sharma/gator/internal/config"
	"github.com/harshaditya-sharma/gator/internal/database"
)

func HandlerLiked(s *config.State, cmd Command, user_id uuid.UUID) error {
	limit := 2
	page := 1

	if len(cmd.Args) >= 1 {
		l, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("usage: %v <number_of_posts> <page_number>", cmd.Name)
		}
		limit = l
	}

	if len(cmd.Args) == 2 {
		p, err := strconv.Atoi(cmd.Args[1])
		if err != nil {
			return fmt.Errorf("usage: %v <number_of_posts> <page_number>", cmd.Name)
		}
		if p < 1 {
			return fmt.Errorf("page number must be 1 or greater")
		}
		page = p
	}

	offset := (page - 1) * limit

	posts, err := s.Db.GetLikedPostsForUser(context.Background(), database.GetLikedPostsForUserParams{
		UserID: user_id,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return fmt.Errorf("error getting liked posts: %v", err)
	}

	printLikedPosts(posts)
	return nil
}

func printLikedPosts(posts interface{}) {
	v := reflect.ValueOf(posts)
	if v.Kind() != reflect.Slice {
		fmt.Println("Error: posts is not a slice")
		return
	}

	for i := 0; i < v.Len(); i++ {
		post := v.Index(i)

		title := post.FieldByName("Title").String()
		url := post.FieldByName("Url").String()
		description := post.FieldByName("Description").FieldByName("String").String()
		publishedAt := post.FieldByName("PublishedAt").FieldByName("Time").MethodByName("Format").Call([]reflect.Value{reflect.ValueOf("Mon Jan 2 2006")})[0].String()
		feedName := post.FieldByName("FeedName").String()

		fmt.Printf("%s (from %s)\n", title, feedName)
		fmt.Printf("Link: %s\n", url)
		fmt.Printf("%s\n", description)
		fmt.Printf("Published: %s\n", publishedAt)
		fmt.Println("=====================================")
	}
}
