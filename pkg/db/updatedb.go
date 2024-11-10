package db

import (
	"fmt"
)

func UpdatePostInfoW(id int64) (err error) {
	_, err = DB.Exec("UPDATE posts SET amount_of_watches = amount_of_watches + 1 WHERE id = ?", id)
	return
}

func updatePostInfo(id int64, countOfAnswers int) (err error) {
	_, err = DB.Exec("UPDATE posts SET amount_of_comments = ? WHERE id = ?", countOfAnswers, id)
	return
}

func DeleteActiveSession(user_id int64) (err error) {
	_, err = DB.Exec("DELETE FROM active_sessions WHERE user_id = ?", user_id)
	return
}

func UpdateUserAmounts(userId int64, what string) (err error) {
	query := fmt.Sprintf("UPDATE users SET %s = %s + 1 WHERE id = ?", what, what)
	_, err = DB.Exec(query, userId)
	return
}

func UpdateRate(what, rate, id string) (err error) {
	query := fmt.Sprintf("UPDATE %ss SET rate = rate + ? WHERE id = ?", what)
	_, err = DB.Exec(query, rate, id)
	return
}
