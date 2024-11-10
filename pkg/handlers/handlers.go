package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/0-LY/Forum-test/pkg/api"
	"github.com/0-LY/Forum-test/pkg/db"
	"github.com/0-LY/Forum-test/pkg/utils"
)

func handleRender(w http.ResponseWriter, r *http.Request, page string, ParserConfig api.ParserConfig) {
	UserAction := UserActionArray[ParserConfig.AddData.SessionCookie]

	if utils.CheckLogOut(ParserConfig.AddData.UserId, ParserConfig.FormData) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if UserActionArray[ParserConfig.AddData.SessionCookie].UserRole != "guest" && page == "authourisation" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if UserAction.UserRole == "guest" && page == "create" {
		http.Redirect(w, r, "/authourisation", http.StatusSeeOther)
	}
	if page == "authourisation" {
		index := strings.Index(UserAction.CurrentURL, "/authourisation")
		if index != -1 {
			UserAction.CurrentURL = UserAction.CurrentURL[:index]
		}
	}

	UserActionArray[ParserConfig.AddData.SessionCookie] = UserAction
	ParserConfig.AddData.SortBy = UserAction.SortBy

	if UserAction.DDOS.Active {
		ParserConfig.AddData.Error = true
		page = "/"
	}

	if strings.Contains(UserAction.CurrentURL, "/authourisation") && page != "authourisation" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	//
	// writeStructToFile(ParserConfig, "testParsingData")
	// writeStructToFile(UserActionArray, "testUserAction")
	//

	renderTemplate(w, page, ParserConfig)
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	ParserConfig := createPageConfig(w, r)
	ParserConfig.AddData.Error = true
	handleRender(w, r, "/", ParserConfig)
}

func AllPosts(w http.ResponseWriter, r *http.Request) {
	ParserConfig := createPageConfig(w, r)
	posts, _ := db.SortBy(ParserConfig, "posts")
	ParserConfig.ParseData = posts
	ParserConfig.AddData.GetAmount(ParserConfig.ParseData)
	handleRender(w, r, "/", ParserConfig)
}

func Search(w http.ResponseWriter, r *http.Request) {
	ParserConfig := createPageConfig(w, r)
	var searchResult api.ParseData
	var err error
	if searchResult, err = db.SearchPosts(r.URL.Query()["search"][0], ParserConfig); err != nil {
		NotFound(w, r)
		return
	}
	ParserConfig.ParseData = searchResult
	ParserConfig.AddData.GetAmount(ParserConfig.ParseData)
	handleRender(w, r, "/", ParserConfig)
}

func GetPostAndComments(w http.ResponseWriter, r *http.Request) {
	ParserConfig := createPageConfig(w, r)
	position := getCurrentPosition(r.URL.Path)
	id, _ := strconv.ParseInt(position, 10, 64)
	post, err := db.ReadPost(id)
	ParserConfig.ParseData = post
	if err != nil {
		NotFound(w, r)
		return
	}
	if UserActionArray[ParserConfig.AddData.SessionCookie].PreviousURL != r.URL.Path {
		db.UpdatePostInfoW(id)
	}
	comments, _ := db.SortBy(ParserConfig, "comments")
	ParserConfig.ParseData = comments
	if _, ok := ParserConfig.FormData["content"]; ok {
		if ParserConfig.AddData.UserRole == "guest" {
			http.Redirect(w, r, position+"/authourisation", http.StatusSeeOther)
			return
		}
		if err, ok := db.AddComments(ParserConfig.FormData, id, ParserConfig.AddData.SessionCookie); err != nil {
			if !ok {
				http.Redirect(w, r, position+"/authourisation", http.StatusSeeOther)
				return
			}
			ParserConfig.AddData.ErrorText = err.Error()
			handleRender(w, r, "post", ParserConfig)
			return
		}
		http.Redirect(w, r, position, http.StatusSeeOther)
	}
	ParserConfig.AddData.GetAmount(ParserConfig.ParseData)
	rateInfo, ok := ParserConfig.FormData["Rate"].(string)
	if ok && rateInfo != "" {
		rateFields := strings.Split(rateInfo, ",")
		if len(rateFields) == 3 {
			if ParserConfig.AddData.UserRole != "guest" {
				if err := db.AddRateTo(ParserConfig.FormData, id, ParserConfig.AddData.UserId); err == nil {
					http.Redirect(w, r, "/post/"+position, http.StatusSeeOther)
					return
				}
			} else {
				http.Redirect(w, r, r.URL.Path+"/authourisation", http.StatusSeeOther)
				return
			}
		}
	}
	ParserConfig = db.CheckIfIRate(ParserConfig)
	handleRender(w, r, "post", ParserConfig)
}

func LoginRegister(w http.ResponseWriter, r *http.Request) {
	ParserConfig := createPageConfig(w, r)
	if len(ParserConfig.FormData) >= 3 {
		if ParserConfig.FormData["action"] == "REGISTER" {
			_, err := utils.RegisterUser(ParserConfig.FormData, ParserConfig.AddData.SessionCookie)
			if err != nil {
				ParserConfig.FormData["authStatus"] = "REGISTER"
				ParserConfig.AddData.ErrorText = err.Error()
				handleRender(w, r, "authourisation", ParserConfig)
				return
			}
			http.Redirect(w, r, UserActionArray[ParserConfig.AddData.SessionCookie].PreviousURL, http.StatusSeeOther)
		} else {
			_, err := utils.LoginUser(ParserConfig.FormData, ParserConfig.AddData.SessionCookie)
			if err != nil {
				ParserConfig.AddData.ErrorText = err.Error()
				ParserConfig.FormData["authStatus"] = "LOGIN"
				handleRender(w, r, "authourisation", ParserConfig)
				return
			}
			http.Redirect(w, r, UserActionArray[ParserConfig.AddData.SessionCookie].PreviousURL, http.StatusSeeOther)
		}
	}
	handleRender(w, r, "authourisation", ParserConfig)
}

func LoginApi(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, "/authourisation/github") {
		utils.GithubLogin(w, r)
	} else if strings.HasSuffix(r.URL.Path, "/authourisation/google") {
		utils.GoogleLogin(w, r)
	}
}

func LoginApiCallback(w http.ResponseWriter, r *http.Request) {
	ParserConfig := createPageConfig(w, r)
	if strings.HasSuffix(r.URL.Path, "/authourisation/github/callback") {
		utils.GithubCallbackHandler(w, r, ParserConfig.AddData.SessionCookie)
	}
	if strings.HasSuffix(r.URL.Path, "/authourisation/google/callback") {
		utils.GoogleCallback(w, r, ParserConfig.AddData.SessionCookie)
	}
	http.Redirect(w, r, UserActionArray[ParserConfig.AddData.SessionCookie].PreviousURL, http.StatusSeeOther)
}
func CreatePost(w http.ResponseWriter, r *http.Request) {
	ParserConfig := createPageConfig(w, r)
	if len(ParserConfig.FormData) == 3 {
		id, err := db.AddPost(ParserConfig.FormData, ParserConfig.AddData.UserId)
		if err != nil {
			ParserConfig.AddData.ErrorText = err.Error()
			handleRender(w, r, "create", ParserConfig)
			return
		}
		http.Redirect(w, r, "/post/"+strconv.Itoa(int(id)), http.StatusSeeOther)
	}
	handleRender(w, r, "create", ParserConfig)
}

func UserProfile(w http.ResponseWriter, r *http.Request) {
	ParserConfig := createPageConfig(w, r)
	position := getCurrentPosition(r.URL.Path)
	id, _ := strconv.ParseInt(position, 10, 64)
	userInfoPublic, err := db.GetUserInfoPublic(id)
	if err != nil {
		NotFound(w, r)
		return
	}
	ParserConfig.AddData.SortByType = "Posts"
	ParserConfig.ParseData.User = append(ParserConfig.ParseData.User, userInfoPublic)
	postsAndCommends, _ := db.SortBy(ParserConfig, "postsAndCommends")
	ParserConfig.ParseData.Posts = postsAndCommends.Posts
	ParserConfig.ParseData.Comments = postsAndCommends.Comments
	ParserConfig.AddData.GetAmount(ParserConfig.ParseData)
	handleRender(w, r, "user", ParserConfig)

}
