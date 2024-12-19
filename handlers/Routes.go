package handlers

import (
	"forum/auth"
	mdlware "forum/middleware"
	"net/http"
)

func Routes() http.Handler {

	router := http.NewServeMux()

	router.HandleFunc("/", HomePage)

	router.HandleFunc("/ws", mdlware.WebSocketHandler)
	router.HandleFunc("/auth/status", CheckAuthHandler)

	router.HandleFunc("/Data-userLogin", LoginHandler)
	router.HandleFunc("/Data-userLogout", LogoutHandler)
	router.HandleFunc("/Data-userRegister", RegisterHandler)

	// ! START Google and Github auth
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
	// ! END Google and Github auth

	router.HandleFunc("/Data-Post", PostHandler)
	router.HandleFunc("/Data-PostLike", PostLikeHandler)
	router.HandleFunc("/Data-PostDisLike", PostDisLikeHandler)

	router.HandleFunc("/Data-Comment", CommentHandler)
	router.HandleFunc("/Data-CommentLike", CommentLikeHandler)
	router.HandleFunc("/Data-CommentDisLike", CommentDislikeHandler)

	router.HandleFunc("/Data-CreatPost", CreatePostHandler)
	router.HandleFunc("/Data-CreatComment", CreatCommentHandler)

	router.HandleFunc("/Data-Profile", ProfileHandler)
	router.HandleFunc("/Data-Categories", CategoriesHandler)

	// ! for testing
	router.HandleFunc("/test", NotificaionHandler)

	router.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./static/uploads"))))
	router.Handle("/scripts/", http.StripPrefix("/scripts/", http.FileServer(http.Dir("./static/scripts"))))
	router.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("./static/styles"))))
	router.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./static/images"))))

	return mdlware.RateLimiter(router)

}
