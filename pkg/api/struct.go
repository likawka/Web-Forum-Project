package api

type Page struct {
	Template   string
	Components []string
}

var Pages = map[string]Page{
	"/": {
		Template:   "all_posts.page.html",
		Components: []string{"sidebar", "header", "sortBox", "postBox", "errorBox"},
	},
	"create": {
		Template:   "create_post.page.html",
		Components: []string{"sidebar", "header", "createPost", "markdownEditor"},
	},
	"post": {
		Template:   "post_and_comments.page.html",
		Components: []string{"sidebar", "header", "createPost", "markdownEditor", "postBox", "sortBox", "transitionBox", "errorBox"},
	},
	"authourisation": {
		Template:   "authourisation.page.html",
		Components: []string{"authBox"},
	},
	"user": {
		Template:   "user.page.html",
		Components: []string{"sidebar", "header", "postBox", "errorBox", "sortBox"},
	},
}

type ParserConfig struct {
	ParseData ParseData
	AddData   AddData
	FormData  map[string]interface{}
}

type Record struct {
	UserRecord
	ContentType string
	Id          int64
	UserId      int64
	Content     string
	Rate        int
	RateType    string
}

type ParseData struct {
	Posts    []Post
	Comments []Comment
	User     []User
}

type Post struct {
	Record
	Title            string
	AmountOfComments int
	AmountOfWatches  int
	Tags             []string
}

type Comment struct {
	Record
	PostId int64
}

type User struct {
	UserRecord
	UserId           int64
	AmountOfPosts    int
	AmountOfComments int
}

type AddData struct {
	UserRole         string
	SessionCookie    string
	UserId           int64
	Error            bool
	ErrorText        string
	SortBy           string
	SortByType       string
	AmountOfPosts    int
	AmountOfComments int
}

type UserRecord struct {
	UserName      string
	CreatedAt     string
	CreateTimeAgo string
}
