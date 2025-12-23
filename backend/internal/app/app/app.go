package app

import (
	"app/internal/auth"
	"app/internal/config"
	"app/internal/db"
	"app/internal/neuronet"
	"app/internal/refreshtoken"
	"app/internal/router"
	"app/internal/user"
	"app/internal/word"

	"github.com/labstack/gommon/log"
)

type App struct {
	Router *router.Router
	Cfg    *config.Config
	Db     *db.Db
}

func NewApp() *App {
	// TODO: init config
	cfg := config.MustLoad()
	log.Info(cfg.Env)

	// TODO: init db
	db := db.NewDbClient(cfg.Database.Url)
	log.Info("Open connect to db", db.Client.Schema)

	// TODO: jwt
	jwtSvc := auth.New(cfg.Jwt.Secret, cfg.Jwt.AccessTtlHours)

	// TODO: refresh token
	refreshTokenRepo := refreshtoken.NewPostgresRepo(db)
	refreshTokenService := refreshtoken.NewService(refreshTokenRepo, cfg.Jwt.RefreshTtlHours)

	// TODO: init chi
	userRepo := user.NewPostgresRepo(db)
	userService := user.NewSercie(userRepo)
	userHandler := user.NewHandler(userService, jwtSvc, refreshTokenService)

	// TODO: word
	neuroService, err := neuronet.NewGeminiRepository(cfg.Gemini.ApiKey)
	if err != nil {
		panic(err)
	}
	wordRepo := word.NewPostgresRepository(db)
	wordService := word.NewService(wordRepo, neuroService)
	wordHandler := word.NewHandler(wordService)

	r := router.NewRouter(userHandler, wordHandler, jwtSvc)

	return &App{
		Router: r,
		Cfg:    cfg,
		Db:     db,
	}
}
