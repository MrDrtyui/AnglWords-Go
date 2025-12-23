package router

import (
	"app/internal/auth"
	"app/internal/middleware"
	"app/internal/user"
	"app/internal/word"
	"net/http"

	_ "app/docs" // Import generated docs

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Router struct {
	chi *chi.Mux
}

func NewRouter(userHandler *user.Handler, wordHandler *word.Handler, jwt *auth.JWT) *Router {
	r := chi.NewRouter()

	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", userHandler.Register)
		r.Post("/login", userHandler.Login)
		r.Post("/refresh", userHandler.Refresh)
		r.Post("/logout", userHandler.Logout)
	})

	r.Route("/words", func(r chi.Router) {
		r.Use(middleware.Auth(jwt))
		r.Post("/", wordHandler.CreateWord)
		r.Get("/my", wordHandler.GetMyWords)
		r.Get("/", wordHandler.GetAllWords)
		r.Get("/{word}", wordHandler.GetWord)
	})

	return &Router{chi: r}
}

func (r *Router) Handler() http.Handler {
	return r.chi
}
