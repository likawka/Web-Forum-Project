package utils

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/0-LY/Forum-test/pkg/db"
)

func addApiUserToDB(githubData map[string]interface{}, way string) (int64, error) {
	var (
		username    string
		err         error
		existingID  int64
	)

	err = db.DB.QueryRow("SELECT id FROM users WHERE id_github = ?", githubData["id"]).Scan(&existingID)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	if existingID != 0 {
		return existingID, nil
	}

	for {
		err := db.DB.QueryRow("SELECT username FROM users WHERE username = ?", githubData["name"]).Scan(&username)
		if err != nil {
			if err != sql.ErrNoRows {
				return 0, err
			}

			break
		}

		githubData["name"] = fmt.Sprintf("%s%d", githubData["name"].(string), 1)

		err = db.DB.QueryRow("SELECT username FROM users WHERE username = ?", githubData["name"]).Scan(&username)
		if err != nil {
			if err != sql.ErrNoRows {
				return 0, err
			}

			break
		}
	}

	result, err := db.DB.Exec("INSERT INTO users (username, id_github, created_at) VALUES (?, ?, ?)",
		githubData["name"].(string), githubData["id"], time.Now())
	if err != nil {
		return 0, err
	}

	userId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return userId, nil
}
