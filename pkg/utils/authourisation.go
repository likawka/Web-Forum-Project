package utils

import (
	"fmt"
	"github.com/0-LY/Forum-test/pkg/db"
	"regexp"
)

func RegisterUser(data map[string]interface{}, sessionID string) (int64, error) {
	if err := checkAuthInput(data); err != nil {
		return 0, err
	}
	userId, err := db.AddUser(data)
	if err != nil {
		return 0, err
	}
	db.AddActiveSession(userId, sessionID)
	return userId, nil
}

func LoginUser(data map[string]interface{}, sesionID string) (int64, error) {
	userId, err := db.GetUser(data)
	if err != nil {
		return 0, err
	}
	db.AddActiveSession(userId, sesionID)
	return userId, nil
}

func checkAuthInput(data map[string]interface{}) error {
	if username, ok := data["username"].(string); ok {
		if len(username) > 16 {
			return fmt.Errorf("the username length should max 16 characters")
		}

	}
	if email, ok := data["email"].(string); ok {
		matched, _ := regexp.MatchString(`[^@\s]+@[^@\s]+\.[^@\s]+`, email)
		if !matched {
			return fmt.Errorf("error email format")
		}
	}

	if password, ok := data["password"].(string); ok {
		if reppassword, ok := data["reppassword"].(string); ok {
			if password != reppassword {
				return fmt.Errorf("passwords do not match")
			}

			regex := regexp.MustCompile(`^.{8,16}$`)
			hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
			hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
			hasDigit := regexp.MustCompile(`\d`).MatchString(password)
			hasSpecial := regexp.MustCompile(`[@$!%*?&]`).MatchString(password)

			if !regex.MatchString(password) || !hasLower || !hasUpper || !hasDigit || !hasSpecial {
				return fmt.Errorf("password should be at least one number and one uppercase and lowercase letter, special character and number of char 8-16")
			}
		}
	}

	return nil
}

func CheckLogOut(id int64, status map[string]interface{}) bool {
	if str, ok := status["Exit"].(string); ok && str == "Log out" {
		db.DeleteActiveSession(id)
		return true
	}
	return false
}
