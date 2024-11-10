package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/0-LY/Forum-test/pkg/db"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	GithubClientID     = "6a31827ffd6ec33bbade"
	GithubClientSecret = "0dd50ce67c4baa5f99f2429aeacf0d8c0b17745f"
)

func GithubLogin(w http.ResponseWriter, r *http.Request) {
	redirectURL := "https://localhost:8080/authourisation/github/callback"
	URL := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s", GithubClientID, redirectURL)
	http.Redirect(w, r, URL, http.StatusSeeOther)
}

func GithubCallbackHandler(w http.ResponseWriter, r *http.Request, sessionID string) {
	code := r.URL.Query().Get("code")
	githubAccessToken := getGithubAccessToken(code)
	githubData := getGithubData(githubAccessToken)
	userId, _ := addApiUserToDB(githubData, "github")
	db.AddActiveSession(userId, sessionID)
}

func getGithubAccessToken(code string) string {
	requestBodyMap := map[string]string{
		"client_id":     GithubClientID,
		"client_secret": GithubClientSecret,
		"code":          code,
	}
	requestJSON, _ := json.Marshal(requestBodyMap)
	req, reqerr := http.NewRequest("POST", "https://github.com/login/oauth/access_token",
		bytes.NewBuffer(requestJSON),
	)
	if reqerr != nil {
		log.Panic("Request creation failed")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		log.Panic("Request failed")
	}
	respbody, _ := ioutil.ReadAll(resp.Body)
	type githubAccessTokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}

	var ghresp githubAccessTokenResponse
	json.Unmarshal(respbody, &ghresp)
	return ghresp.AccessToken
}

func getGithubData(accessToken string) map[string]interface{} {
	req, reqerr := http.NewRequest("GET", "https://api.github.com/user", nil)
	if reqerr != nil {
		log.Panic("API Request creation failed")
	}

	authorizationHeaderValue := fmt.Sprintf("token %s", accessToken)
	req.Header.Set("Authorization", authorizationHeaderValue)

	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		log.Panic("Request failed")
	}
	respbody, _ := ioutil.ReadAll(resp.Body)

	var responseBodyMap map[string]interface{}
	json.Unmarshal(respbody, &responseBodyMap)
	return responseBodyMap
}
