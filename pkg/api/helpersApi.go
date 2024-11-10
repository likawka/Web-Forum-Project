package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (ur *UserRecord) CreateCustomCreateTimeAgo() {
	dataTime, _ := time.Parse(time.RFC3339, ur.CreatedAt)
	diff := time.Now().Sub(dataTime)

	years := diff / (365 * 24 * time.Hour)
	diff -= years * (365 * 24 * time.Hour)

	months := diff / (30 * 24 * time.Hour)
	diff -= months * (30 * 24 * time.Hour)

	days := diff / (24 * time.Hour)
	diff -= days * (24 * time.Hour)

	hours := diff / time.Hour
	diff -= hours * time.Hour

	minutes := diff / time.Minute
	diff -= minutes * time.Minute

	seconds := diff / time.Second

	var result string

	if years > 0 {
		result += fmt.Sprintf("%d year", years)
		if years > 1 {
			result += "s"
		}
		if months > 0 {
			result += fmt.Sprintf(", %d month", months)
			if months > 1 {
				result += "s"
			}
		}
	} else if months > 0 {
		result += fmt.Sprintf("%d month", months)
		if months > 1 {
			result += "s"
		}
	} else if days > 0 {
		result += fmt.Sprintf("%d day", days)
		if days > 1 {
			result += "s"
		}
	} else if hours > 0 {
		result += fmt.Sprintf("%d hour", hours)
		if hours > 1 {
			result += "s"
		}
	} else if minutes > 0 {
		result += fmt.Sprintf("%d minute", minutes)
		if minutes > 1 {
			result += "s"
		}
	} else if seconds > 0 {
		result += fmt.Sprintf("%d second", seconds)
		if seconds > 1 {
			result += "s"
		}
	}

	if result != "" {
		result += " ago"
	} else {
		result = "just now"
	}
	ur.CreatedAt = dataTime.Format("2006-01-02")
	ur.CreateTimeAgo = result
}

func (ad *AddData) GetAmount(ParseData ParseData) {
	ad.AmountOfPosts = len(ParseData.Posts)
	ad.AmountOfComments = len(ParseData.Comments)
}

func (pc *ParserConfig) ReadPage(w http.ResponseWriter, r *http.Request) bool {
	pc.AddData.CheckSession(w, r)
	r.ParseForm()
	pc.FormData = make(map[string]interface{})
	for key, values := range r.Form {
		if len(values) > 0 {
			pc.FormData[key] = values[0]
		}
	}
	if r.Method != http.MethodPost {
		return false
	}
	return true
}

func (ad *AddData) CheckSession(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie("session_id")
	if err != nil && err != http.ErrNoCookie {
		return err
	}
	if cookie == nil {
		ad.createSession(w)
	}

	ad.UserRole = "guest"
	if cookie != nil {
		ad.SessionCookie = cookie.Value
	}

	return nil
}

func (pc *Post) ShortPost() {
	maxContentLength := 200
	maxLineCount := 2
	ch := false

	lines := strings.Split(pc.Content, "\n")
	if len(lines) > maxLineCount {
		pc.Content = strings.Join(lines[:maxLineCount-1], "\n")
		pc.Content += lines[maxLineCount-1][:len(lines[maxLineCount-1])-1] + " "
		ch = true
	}

	if len(pc.Content) > maxContentLength {
		pc.Content = pc.Content[:maxContentLength]
		ch = true
	}
	if ch {
		pc.Content += "..."
	}

}

func (ad *AddData) createSession(w http.ResponseWriter) {
	sessionID := uuid.New().String()
	cookie := http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, &cookie)
	ad.SessionCookie = sessionID
}
