package db

import (
	"fmt"
	"strings"

	"github.com/0-LY/Forum-test/pkg/api"
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

func SortBy(parserConfig api.ParserConfig, what string) (api.ParseData, error) {
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
	var id int64
	if len(parserConfig.ParseData.User) == 1 {
		id = parserConfig.ParseData.User[0].UserId
	} else {
		id = 0
	}
	if what == "posts" {
		parserConfig.ParseData.Posts, _ = readPosts(id, sortByField, sortOrder)
	}
	if what == "comments" && len(parserConfig.ParseData.Posts) > 0 {
		parserConfig.ParseData.Comments, _ = readPostComments(parserConfig.ParseData.Posts[0].Id, sortByField, sortOrder, "post_id")
	}
	if what == "postsAndCommends" {
		parserConfig.ParseData.Posts, _ = readPosts(id, sortByField, sortOrder)
		parserConfig.ParseData.Comments, _ = readPostComments(id, sortByField, sortOrder, "user_id")
	}

	return parserConfig.ParseData, nil
}

func checkInputData(data map[string]interface{}) error {
	checkLength := func(s string, min, max int) error {
		s = strings.TrimSpace(s)
		if len(s) < min || len(s) > max {
			return fmt.Errorf("length should be between %d and %d characters", min, max)
		}
		return nil
	}

	checkCategoryLength := func(s string) error {
		return checkLength(s, 2, 50)
	}

	checkContentLength := func(s string, min, max int) error {
		err := checkLength(s, min, max)
		if err != nil {
			return err
		}
		return nil
	}

	if title, ok := data["title"].(string); ok {
		if err := checkContentLength(title, 5, 100); err != nil {
			return err
		}
	}

	if contentAns, ok := data["content_ans"].(string); ok {
		if err := checkContentLength(contentAns, 100, 1000); err != nil {
			return err
		}
	}

	if content, ok := data["content"].(string); ok {
		if err := checkContentLength(content, 100, 10000); err != nil {
			return err
		}
	}

	if categories, ok := data["categories"].(string); ok {
		if err := checkCategoryLength(categories); err != nil {
			return err
		}

		tags := removeDuplicatesCategories(strings.Split(categories, " "))
		for _, tag := range tags {
			if len(tag) <= 1 {
				continue
			}
			if tag[0] != '#' {
				return fmt.Errorf("each tag should start with # symbol")
			}
		}
	}

	return nil
}

func removeDuplicatesCategories(elements []string) []string {
	uniqueTags := make(map[string]bool)

	for _, tag := range elements {
		uniqueTags[tag] = true
	}

	uniqueTagsArray := make([]string, 0, len(uniqueTags))
	for tag := range uniqueTags {
		uniqueTagsArray = append(uniqueTagsArray, tag)
	}

	return uniqueTagsArray
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func checkPassword(password, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

func UserExists(username, email string) (bool, string, error) {
	var existingField string
	var count int
	err := DB.QueryRow("SELECT CASE WHEN username = ? THEN 'username' ELSE 'email' END, COUNT(*) FROM users WHERE username = ? OR email = ?", username, username, email).Scan(&existingField, &count)
	if err != nil {
		return false, "", err
	}
	return count > 0, existingField, nil
}

func CheckIfIRate(Parse api.ParserConfig) api.ParserConfig {
	for id, post := range Parse.ParseData.Posts {
		var existingStatus string
		err := DB.QueryRow("SELECT status FROM rates WHERE user_id = ? AND post_id = ?", Parse.AddData.UserId, post.Id).Scan(&existingStatus)
		if err != nil {
			continue
		}
		Parse.ParseData.Posts[id].RateType = existingStatus
	}
	for id, comment := range Parse.ParseData.Comments {
		var existingStatus string
		err := DB.QueryRow("SELECT status FROM rates WHERE user_id = ? AND comment_id = ?", Parse.AddData.UserId, comment.Id).Scan(&existingStatus)
		if err != nil {
			continue
		}
		Parse.ParseData.Comments[id].RateType = existingStatus
	}
	return Parse
}

func parseSearchHints(what string) ([]string, string, string, string) {
	var tags []string
	var title, date, user string

	tagRegex := regexp.MustCompile(`\[[^\]]+\]`)
	matches := tagRegex.FindAllString(what, -1)
	for _, match := range matches {
		tags = append(tags, strings.ToLower(match[1:len(match)-1]))
	}

	paramRegex := regexp.MustCompile(`\b(user|title|date):\w+\b`)
	params := paramRegex.FindAllString(what, -1)
	for _, param := range params {
		parts := strings.Split(param, ":")
		switch parts[0] {
		case "user":
			user = parts[1]
		case "title":
			title = parts[1]
		case "date":
			date = parts[1]
		}
	}

	return tags, user, title, date
}
