package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/minhhoccode111/realworld-fiber-clean/config"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/controller/http/middleware"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/usecase"
	"github.com/minhhoccode111/realworld-fiber-clean/pkg/logger"
	"github.com/minhhoccode111/realworld-fiber-clean/pkg/validatorx"
)

// NewV1Routes -.
func NewV1Routes(
	apiV1Group fiber.Router,
	cfg *config.Config,
	l logger.Interface,

	t usecase.Translation,
	tc usecase.TranslationClone,
	u usecase.User,
	a usecase.Article,
	tag usecase.Tag,
) {
	r := &V1{
		cfg: cfg,
		l:   l,
		v:   validatorx.New(),

		t:   t,
		tc:  tc,
		u:   u,
		a:   a,
		tag: tag,
	}

	jwtSecret := r.cfg.JWT.Secret
	auth := middleware.AuthMiddleware(l, jwtSecret, false)
	// optionalAuth := middleware.AuthMiddleware(l, jwtSecret, true)

	translation := apiV1Group.Group("/translation")
	{
		translation.Get("/history", r.history)
		translation.Post("/do-translate", r.doTranslate)
	}

	translationClone := apiV1Group.Group("/translation-clone")
	{
		translationClone.Get("/history", r.getHistory)
		translationClone.Post("/translate", r.postTranslate)
	}

	users := apiV1Group.Group("/users")
	{
		users.Post("/", r.postRegisterUser)
		users.Post("/login", r.postLoginUser)
	}

	user := apiV1Group.Group("/user")
	{
		user.Get("/", auth, r.getCurrentUser)
		user.Put("/", auth, r.putUpdateUser)
	}

	articles := apiV1Group.Group("/articles")
	{
		articles.Post("/", auth, r.postCreateArticle)
	}

	tags := apiV1Group.Group("/tags")
	{
		tags.Get("/", r.getTags)
	}
}
