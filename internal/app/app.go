// Package app configures and runs application.
package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/minhhoccode111/realworld-fiber-clean/config"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/controller/http"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/repo/persistent"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/repo/webapi"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/usecase/article"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/usecase/comment"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/usecase/favorite"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/usecase/profile"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/usecase/tag"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/usecase/translation"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/usecase/user"
	"github.com/minhhoccode111/realworld-fiber-clean/pkg/httpserver"
	"github.com/minhhoccode111/realworld-fiber-clean/pkg/logger"
	"github.com/minhhoccode111/realworld-fiber-clean/pkg/postgres"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) { //nolint: gocyclo,cyclop,funlen,gocritic,nolintlint
	l := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// Use-Case
	translationUseCase := translation.New(persistent.New(pg), webapi.New())
	userUseCase := user.New(persistent.NewUserRepo(pg))
	articleUseCase := article.New(persistent.NewArticleRepo(pg))
	favoriteUseCase := favorite.New(persistent.NewFavoriteRepo(pg))
	commentUseCase := comment.New(persistent.NewCommentRepo(pg))
	profileUseCase := profile.New(persistent.NewProfileRepo(pg))
	tagUseCase := tag.New(persistent.NewTagRepo(pg))

	// HTTP Server
	httpServer := httpserver.New(
		l,
		httpserver.Port(cfg.HTTP.Port),
		httpserver.Prefork(cfg.HTTP.UsePreforkMode),
	)
	http.NewRouter(
		httpServer.App,
		cfg,
		l,

		translationUseCase,
		userUseCase,
		articleUseCase,
		favoriteUseCase,
		commentUseCase,
		profileUseCase,
		tagUseCase,
	)

	// Start servers
	httpServer.Start()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: %s", s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
