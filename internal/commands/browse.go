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

func HandlerBrowse(s *config.State, cmd Command, user_id uuid.UUID) error {
	limit := 2
	page := 1
	args := cmd.Args
	var sortOrder string
	var feedFilter string

	// Parse args manually to handle flags
	newArgs := []string{}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--sort" || arg == "-s" {
			if i+1 < len(args) {
				sortOrder = args[i+1]
				i++
			}
		} else if arg == "--feed" || arg == "-f" {
			if i+1 < len(args) {
				feedFilter = args[i+1]
				i++
			}
		} else {
			newArgs = append(newArgs, arg)
		}
	}

	if len(newArgs) > 2 {
		return fmt.Errorf("usage: %v <number_of_posts> <page_number>", cmd.Name)
	}

	if len(newArgs) >= 1 {
		l, err := strconv.Atoi(newArgs[0])
		if err != nil {
			return fmt.Errorf("usage: %v <number_of_posts> <page_number>", cmd.Name)
		}
		limit = l
	}

	if len(newArgs) == 2 {
		p, err := strconv.Atoi(newArgs[1])
		if err != nil {
			return fmt.Errorf("usage: %v <number_of_posts> <page_number>", cmd.Name)
		}
		if p < 1 {
			return fmt.Errorf("page number must be 1 or greater")
		}
		page = p
	}

	offset := (page - 1) * limit

	var posts []database.GetPostsForUserRow
	var err error

	// Default to 'desc' if not provided
	if sortOrder == "" {
		sortOrder = "desc"
	}

	// Choose the correct query based on flags
	if feedFilter != "" {
		if sortOrder == "asc" {
			var p []database.GetPostsForUserByFeedAscRow
			p, err = s.Db.GetPostsForUserByFeedAsc(context.Background(), database.GetPostsForUserByFeedAscParams{
				UserID: user_id,
				Name:   feedFilter,
				Limit:  int32(limit),
				Offset: int32(offset),
			})
			if err == nil {
				printPosts(p)
				return nil
			}
		} else {
			var p []database.GetPostsForUserByFeedRow
			p, err = s.Db.GetPostsForUserByFeed(context.Background(), database.GetPostsForUserByFeedParams{
				UserID: user_id,
				Name:   feedFilter,
				Limit:  int32(limit),
				Offset: int32(offset),
			})
			if err == nil {
				printPosts(p)
				return nil
			}
		}
	} else {
		if sortOrder == "asc" {
			var p []database.GetPostsForUserAscRow
			p, err = s.Db.GetPostsForUserAsc(context.Background(), database.GetPostsForUserAscParams{
				UserID: user_id,
				Limit:  int32(limit),
				Offset: int32(offset),
			})
			if err == nil {
				printPosts(p)
				return nil
			}
		} else {
			// Default case
			posts, err = s.Db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
				UserID: user_id,
				Limit:  int32(limit),
				Offset: int32(offset),
			})
		}
	}

	if err != nil {
		return fmt.Errorf("Error getting posts: %v", err)
	}

	printPosts(posts)
	return nil
}

func printPosts(posts interface{}) {
	v := reflect.ValueOf(posts)
	if v.Kind() != reflect.Slice {
		fmt.Println("Error: posts is not a slice")
		return
	}

	for i := 0; i < v.Len(); i++ {
		post := v.Index(i)
		// Access fields by name
		title := post.FieldByName("Title").String()
		url := post.FieldByName("Url").String()
		description := post.FieldByName("Description").FieldByName("String").String()
		publishedAt := post.FieldByName("PublishedAt").FieldByName("Time").MethodByName("Format").Call([]reflect.Value{reflect.ValueOf("Mon Jan 2 2006")})[0].String()

		fmt.Printf("%s\n", title)
		fmt.Printf("Link: %s\n", url)
		fmt.Printf("%s\n", description)
		fmt.Printf("Published: %s\n", publishedAt)
		fmt.Println("=====================================")
	}
}
