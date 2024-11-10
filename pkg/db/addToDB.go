package db

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

func AddPost(data map[string]interface{}, UserId int64) (int64, error) {
	if err := checkInputData(data); err != nil {
		return 0, err
	}
	result, err := DB.Exec("INSERT INTO posts (title, content, created_date, user_id) VALUES (?, ?, ?, ?)",
		data["title"], data["content"], time.Now(), UserId)
	if err != nil {
		return 0, err
	}

	postId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	if categories, ok := data["categories"]; ok {
		categoriesStr, ok := categories.(string)
		if !ok {
			return 0, fmt.Errorf("invalid data format for categories")
		}
		categoryList := removeDuplicatesCategories(strings.Split(strings.TrimSpace(categoriesStr), "#"))
		for _, category := range categoryList {
			category = strings.TrimSpace(category)
			if len(category) > 0 {
				categoryId, err := getCategoryId(category)
				if err != nil {
					return 0, err
				}
				if categoryId == 0 {
					categoryId, err = AddCategory(category)
					if err != nil {
						return 0, err
					}
				}
				_, err = DB.Exec("INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)", postId, categoryId)
				if err != nil {
					return 0, err
				}
			}
		}
	}
	UpdateUserAmounts(UserId, "amount_of_posts")
	return postId, nil
}

func AddCategory(name string) (int64, error) {
	categoryId, err := getCategoryId(name)
	if err != nil {
		return 0, err
	}

	if categoryId != 0 {
		return categoryId, nil
	}

	result, err := DB.Exec("INSERT INTO categories (name) VALUES (?)", name)
	if err != nil {
		return 0, err
	}

	categoryId, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return categoryId, nil
}

func AddComments(data map[string]interface{}, postId int64, cookie string) (error, bool) {
	content, ok := data["content"].(string)
	if !ok {
		return fmt.Errorf("mandatory parameter 'content' is missing or not a string"), true
	}

	if err := checkInputData(data); err != nil {
		return err, true
	}

	userId, err := GetUserIDBySessionID(cookie)
	if err != nil {
		return err, false
	}

	query := "INSERT INTO comments (post_id, content, user_id, created_date) VALUES (?, ?, ?, ?)"
	_, err = DB.Exec(query, postId, content, userId, time.Now())
	if err != nil {
		return err, true
	}
	UpdateUserAmounts(userId, "amount_of_comments")
	return nil, true
}

func AddUser(data map[string]interface{}) (int64, error) {
	exists, field, err := UserExists(data["username"].(string), data["email"].(string))
	if err != nil {
		return 0, err
	}

	if exists {
		return 0, fmt.Errorf("%s already exists", field)
	}

	hashPassword, _ := hashPassword(data["password"].(string))
	result, err := DB.Exec("INSERT INTO users (username, email, password_hash, created_at) VALUES (?, ?, ?, ?)", data["username"], data["email"], hashPassword, time.Now())
	if err != nil {
		return 0, err
	}

	userId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return userId, nil
}

func AddActiveSession(userID int64, sessionID string) error {
	var userExists bool
	err := DB.QueryRow("SELECT EXISTS (SELECT 1 FROM active_sessions WHERE user_id = ?)", userID).Scan(&userExists)
	if err != nil {
		return err
	}

	if userExists {
		_, err = DB.Exec("UPDATE active_sessions SET session_id = ?, created_at = ?  WHERE user_id = ?", sessionID, time.Now(), userID)
		return err
	} else {
		_, err := DB.Exec("INSERT INTO active_sessions (user_id, session_id, created_at) VALUES (?, ?, ?)", userID, sessionID, time.Now())
		return err
	}
}
func AddRateTo(data map[string]interface{}, postId, userId int64) error {
	rateInfo := strings.Split(data["Rate"].(string), ",")
	var existingStatus string
	err := DB.QueryRow("SELECT status FROM rates WHERE user_id = ? AND "+rateInfo[1]+"_id = ?", userId, rateInfo[2]).Scan(&existingStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err = DB.Exec("INSERT INTO rates (user_id, "+rateInfo[1]+"_id, status, rated_at) VALUES (?, ?, ?, ?)",
				userId, rateInfo[2], rateInfo[0], time.Now())
			if err != nil {
				return err
			}
			UpdateRate(rateInfo[1], rateInfo[0], rateInfo[2])
		}
	} else {
		_, err := DB.Exec("DELETE FROM rates WHERE user_id = ? AND "+rateInfo[1]+"_id = ?", userId, rateInfo[2])
		if err != nil {
			return err
		}
		if existingStatus == "-1" {
			rateInfo[0] = "1"
		} else {
			rateInfo[0] = "-1"
		}

		UpdateRate(rateInfo[1], rateInfo[0], rateInfo[2])
	}

	return nil
}
