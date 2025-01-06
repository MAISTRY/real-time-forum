package main

import (
	"RTF/DB"
	"RTF/auth"
	"RTF/handlers"
	"RTF/middleware"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	srvr := http.Server{
		Addr:    ":443",
		Handler: Routes(),
	}

	DB.InitDB()
	log.Println("starting server on https://localhost/")
	err := srvr.ListenAndServeTLS("./cert/cert.pem", "./cert/key.pem")
	if err != nil {
		log.Fatalf("error starting server:%v", err)
	}
}

func Routes() http.Handler {

	router := http.NewServeMux()

	router.HandleFunc("/", handlers.HomePage)

	router.HandleFunc("/ws", middleware.WebSocketHandler)
	router.HandleFunc("/auth/status", handlers.CheckAuthHandler)

	router.HandleFunc("/Data-userLogin", handlers.LoginHandler)
	router.HandleFunc("/Data-userLogout", handlers.LogoutHandler)
	router.HandleFunc("/Data-userRegister", handlers.RegisterHandler)

	router.HandleFunc("/auth/google/login", auth.HandleGoogleLogin)
	router.HandleFunc("/auth/github/login", auth.HandleGitHubLogin)
	router.HandleFunc("/auth/google/callback", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		q.Add("provider", "google")
		r.URL.RawQuery = q.Encode()
		auth.HandleOAuthCallback(w, r)
	})
	router.HandleFunc("/auth/github/callback", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		q.Add("provider", "github")
		r.URL.RawQuery = q.Encode()
		auth.HandleOAuthCallback(w, r)
	})

	router.HandleFunc("/Data-Post", handlers.PostHandler)
	router.HandleFunc("/Data-PostLike", handlers.PostLikeHandler)
	router.HandleFunc("/Data-PostDisLike", handlers.PostDisLikeHandler)

	router.HandleFunc("/Data-Comment", handlers.CommentHandler)
	router.HandleFunc("/Data-CommentLike", handlers.CommentLikeHandler)
	router.HandleFunc("/Data-CommentDisLike", handlers.CommentDislikeHandler)

	router.HandleFunc("/Data-CreatPost", handlers.CreatePostHandler)
	router.HandleFunc("/Data-CreatComment", handlers.CreatCommentHandler)

	router.HandleFunc("/Data-Profile", handlers.ProfileHandler)
	router.HandleFunc("/Data-Categories", handlers.CategoriesHandler)

	// ! for testing
	router.HandleFunc("/test", handlers.NotificaionHandler)

	router.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./static/uploads"))))
	router.Handle("/scripts/", http.StripPrefix("/scripts/", http.FileServer(http.Dir("./static/scripts"))))
	router.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("./static/styles"))))
	router.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./static/images"))))

	return middleware.RateLimiter(router)

}
