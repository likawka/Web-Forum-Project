package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/0-LY/Forum-test/pkg/api"
)

// Funcs Read

func ReadPost(id int64) (api.ParseData, error) {
	parseData := api.ParseData{}

	comments, err := readPostComments(id, "created_date", "DESC", "post_id")
	if err != nil {
		return api.ParseData{}, err
	}
	parseData.Comments = comments

	post := api.Post{}
	row := DB.QueryRow("SELECT * FROM posts WHERE id = ?", id)
	err = row.Scan(&post.Record.Id, &post.Record.UserId, &post.Title,
		&post.Record.Content, &post.Record.UserRecord.CreatedAt, &post.AmountOfComments,
		&post.AmountOfWatches, &post.Record.Rate)
	if err != nil {
		return api.ParseData{}, err
	}
	post.Record.UserRecord.CreateCustomCreateTimeAgo()
	tags, err := readPostTags(id)
	if err != nil {
		return api.ParseData{}, err
	}
	post.Tags = tags
	post.Record.UserName = getNamebyID(post.Record.UserId)
	post.Record.ContentType = "post"
	parseData.Posts = append(parseData.Posts, post)
	return parseData, nil
}

func getNamebyID(id int64) string {
	var name string
	err := DB.QueryRow("SELECT username FROM users WHERE id = ?", id).Scan(&name)
	if err != nil {
		return "?"
	}
	return name
}

func readPostTags(postID int64) ([]string, error) {
	rows, err := DB.Query("SELECT c.name FROM categories c JOIN post_categories pc ON c.id = pc.category_id WHERE pc.post_id = ?", postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tagName string
		if err := rows.Scan(&tagName); err != nil {
			return nil, err
		}
		tags = append(tags, tagName)
	}

	return tags, nil
}

func readPostComments(postID int64, sortBy string, sortOrder string, what string) ([]api.Comment, error) {
	var query string
	if what == "post_id" {
		query = fmt.Sprintf("SELECT * FROM comments WHERE post_id = ? ORDER BY %s %s", sortBy, sortOrder)
	} else {
		query = fmt.Sprintf("SELECT * FROM comments WHERE user_id = ? ORDER BY %s %s", sortBy, sortOrder)
	}
	rows, err := DB.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []api.Comment
	for rows.Next() {
		var comment api.Comment

		if err := rows.Scan(&comment.Record.Id, &comment.PostId, &comment.Record.UserId,
			&comment.Record.Content, &comment.Record.UserRecord.CreatedAt, &comment.Record.Rate); err != nil {
			return nil, err
		}

		comment.Record.ContentType = "comment"
		comment.Record.UserRecord.CreateCustomCreateTimeAgo()
		comment.Record.UserName = getNamebyID(comment.Record.UserId)
		comments = append(comments, comment)
	}
	updatePostInfo(postID, len(comments))

	return comments, nil
}

func readPosts(id int64, sortBy string, sortOrder string) ([]api.Post, error) {
	var query string
	if id != 0 {
		query = fmt.Sprintf("SELECT * FROM posts WHERE user_id = %d ORDER BY %s %s", id, sortBy, sortOrder)
	} else {
		query = fmt.Sprintf("SELECT * FROM posts ORDER BY %s %s", sortBy, sortOrder)
	}
	rows, err := DB.Query(query)
	if err != nil {
		return []api.Post{}, err
	}
	defer rows.Close()

	var posts []api.Post
	for rows.Next() {
		var post api.Post

		if err := rows.Scan(&post.Record.Id, &post.Record.UserId, &post.Title,
			&post.Record.Content, &post.Record.UserRecord.CreatedAt, &post.AmountOfComments,
			&post.AmountOfWatches, &post.Record.Rate); err != nil {
		}
		post.Record.UserRecord.CreateCustomCreateTimeAgo()
		tags, err := readPostTags(post.Record.Id)
		if err != nil {
			return []api.Post{}, err
		}
		post.Tags = tags
		post.Record.UserName = getNamebyID(post.Record.UserId)
		post.Record.ContentType = "posts"
		post.ShortPost()
		posts = append(posts, post)
	}
	return posts, nil
}

func GetUserIDBySessionID(sessionID string) (int64, error) {
	query := "SELECT user_id FROM active_sessions WHERE session_id = ?;"

	row := DB.QueryRow(query, sessionID)

	var userID int64

	err := row.Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("user with session ID %s not found", sessionID)
		}
		return 0, err
	}

	return userID, nil
}

func GetUser(data map[string]interface{}) (int64, error) {
	var userId int64
	var hashedPassword string

	err := DB.QueryRow("SELECT id, password_hash FROM users WHERE email = ?",
		data["email"]).Scan(&userId, &hashedPassword)
	if err != nil {
		return 0, fmt.Errorf("invalid username or password")
	}

	err = checkPassword(data["password"].(string), hashedPassword)
	if err != nil {
		return 0, fmt.Errorf("invalid username or password")
	}

	return userId, nil
}

// Funcs Get

func getCategoryId(name string) (int64, error) {
	var categoryId int64
	err := DB.QueryRow("SELECT id FROM categories WHERE name = ?", name).Scan(&categoryId)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return categoryId, err
}

func GetUserInfoPublic(id int64) (api.User, error) {
	user := api.User{}
	err := DB.QueryRow("SELECT id, username, created_at, amount_of_posts, amount_of_comments  FROM users WHERE id = ?",
		id).Scan(&user.UserId, &user.UserRecord.UserName, &user.UserRecord.CreatedAt, &user.AmountOfPosts, &user.AmountOfComments)
	if err != nil {
		return api.User{}, err
	}
	user.UserRecord.CreateCustomCreateTimeAgo()
	return user, nil
}

func SearchPosts(what string, parserConfig api.ParserConfig) (api.ParseData, error) {
    var sortByField string
    var sortOrder string
    switch parserConfig.AddData.SortBy {
    case "Oldest":
        sortByField = "created_date"
        sortOrder = "ASC"
    case "Popular":
        sortByField = "rate"
        sortOrder = "DESC"
    default:
        sortByField = "created_date"
        sortOrder = "DESC"
    }
	tags, _, _, _ := parseSearchHints(what)

    if len(tags) == 0 {
        return api.ParseData{}, errors.New("empty tags array")
    }

    query := `
        SELECT DISTINCT p.*
        FROM posts p
        INNER JOIN post_categories pc ON p.id = pc.post_id
        INNER JOIN categories c ON pc.category_id = c.id
        WHERE c.name IN (` + strings.Repeat("?,", len(tags)-1) + `?) 
        ORDER BY ` + sortByField + ` ` + sortOrder

    args := make([]interface{}, len(tags))
    for i, tag := range tags {
        args[i] = tag
    }

    rows, err := DB.Query(query, args...)
    if err != nil {
        return api.ParseData{}, err 
    }
    defer rows.Close()

    parsedData := api.ParseData{}
    var posts []api.Post
    for rows.Next() {
        var post api.Post

        if err := rows.Scan(&post.Record.Id, &post.Record.UserId, &post.Title,
            &post.Record.Content, &post.Record.UserRecord.CreatedAt, &post.AmountOfComments,
            &post.AmountOfWatches, &post.Record.Rate); err != nil {
            return api.ParseData{}, err
        }
        post.Record.UserRecord.CreateCustomCreateTimeAgo()
        tags, err := readPostTags(post.Record.Id)
        if err != nil {
            return api.ParseData{}, err 
        }
        post.Tags = tags
        post.Record.UserName = getNamebyID(post.Record.UserId)
        post.Record.ContentType = "posts"
        post.ShortPost()
        posts = append(posts, post)
    }
    parsedData.Posts = posts

    return parsedData, nil
}
