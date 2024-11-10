package handlers

import (
	"net/http"
	"strings"

	"github.com/0-LY/Forum-test/pkg/api"
	"github.com/0-LY/Forum-test/pkg/db"
	"time"
)

func getCurrentPosition(url string) string {
	pathParts := strings.Split(url, "/")
	lastPart := pathParts[len(pathParts)-1]
	return lastPart
}

func createPageConfig(w http.ResponseWriter, r *http.Request) api.ParserConfig {
	ParserConfig := api.ParserConfig{}
	ParserConfig.ReadPage(w, r)
	userId, err := db.GetUserIDBySessionID(ParserConfig.AddData.SessionCookie)
	if err == nil {
		ParserConfig.AddData.UserRole = "user"
		ParserConfig.AddData.UserId = userId
	}
	UserActionArray[ParserConfig.AddData.SessionCookie] = AddAction(ParserConfig, r.URL.Path, "/authourisation")
	ParserConfig.AddData.SortBy = UserActionArray[ParserConfig.AddData.SessionCookie].SortBy
	return ParserConfig
}

var UserActionArray = make(map[string]UserAction)

type UserAction struct {
	UserRole    string
	DDOS        DDOS
	PreviousURL string
	CurrentURL  string
	SortBy      string
	SortByType  string
	FormData    map[string]interface{}
}

type DDOS struct {
	Active  bool
	EndTime time.Time
}

func AddAction(info api.ParserConfig, url, ignoreURL string) UserAction {
	previousURL := ""
	ua, ok := UserActionArray[info.AddData.SessionCookie]
	if ok {
		previousURL = ua.CurrentURL
	}

	var formData map[string]interface{}
	if strings.Contains(url, ignoreURL) {
		formData = ua.FormData
	} else {
		formData = info.FormData
	}

	SortBy := "Newest"
	if !ok || previousURL != url {
		SortBy = "Newest"
	} else if val, ok := info.FormData["SortBy"].(string); ok && val != "" {
		SortBy = val
	} else {
		SortBy = ua.SortBy
	}

	ua = UserAction{
		UserRole:    info.AddData.UserRole,
		PreviousURL: previousURL,
		DDOS:        ua.DDOS,
		CurrentURL:  url,
		SortBy:      SortBy,
		SortByType:  info.AddData.SortByType,
		FormData:    formData,
	}

	if ua.CurrentURL == "/Too many requests" {
		ua.DDOS = DDOS{Active: true, EndTime: time.Now().Add(1 * time.Minute)}
	}
	if time.Now().After(ua.DDOS.EndTime) {
		ua.DDOS = DDOS{Active: false, EndTime: time.Time{}}
	}
	return ua
}
