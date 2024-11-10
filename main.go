package main

import (
	"github.com/0-LY/Forum-test/pkg/db"
	"github.com/0-LY/Forum-test/pkg/handlers"
	"github.com/0-LY/Forum-test/pkg/router"
	"log"
	"net/http"
)

const (
	portNumber = ":8080"
	certFile   = "certs/localhost.crt"
	keyFile    = "certs/localhost.key"
)

func main() {
	r := router.Router{}
	r.NewRoute(`/`, handlers.AllPosts)
	r.NewRoute(`/search.*`, handlers.Search)
	r.NewRoute("/create", handlers.CreatePost)
	r.NewRoute(`.*/authourisation`, handlers.LoginRegister)
	r.NewRoute(`.*/authourisation/github`, handlers.LoginApi)
	r.NewRoute(`.*/authourisation/google`, handlers.LoginApi)
	r.NewRoute(`.*/authourisation/github/callback`, handlers.LoginApiCallback)
	r.NewRoute(`.*/authourisation/google/callback`, handlers.LoginApiCallback)
	r.NewRoute(`/post/.*`, handlers.GetPostAndComments)
	r.NewRoute(`/user/(?P<id>\d+)`, handlers.UserProfile)
	r.NewRoute(`.*`, handlers.NotFound)
	r.SetRateLimiter(5, 2)

	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("web/assets"))))
	http.HandleFunc("/", r.Serve)

	err := db.InitDB("pkg/db/data/mydatabase.db", "pkg/db/data/database.sql")
	if err != nil {
		log.Println("Error initializing database:", err)
		return
	}
	log.Println("Ctrl + Click on the link: https://localhost" + portNumber)
	log.Println("To stop the server press `Ctrl + C`")

	err = router.StartServer(portNumber, certFile, keyFile)
	if err != nil {
		log.Fatal("Server error: ", err)
	}
}
