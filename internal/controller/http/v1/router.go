package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/minhhoccode111/realworld-fiber-clean/config"
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
	tag usecase.Tag,
) {
	r := &V1{
		cfg: cfg,
		l:   l,
		v:   validatorx.New(),

		t:   t,
		tc:  tc,
		u:   u,
		tag: tag,
	}

	translationGroup := apiV1Group.Group("/translation")

	{
		translationGroup.Get("/history", r.history)
		translationGroup.Post("/do-translate", r.doTranslate)
	}

	translationGroupClone := apiV1Group.Group("/translation-clone")
	{
		translationGroupClone.Get("/history", r.getHistory)
		translationGroupClone.Post("/translate", r.postTranslate)
	}

	usersGroup := apiV1Group.Group("/users")
	{
		usersGroup.Post("/", r.postRegisterUser)
	}

	tagsGroup := apiV1Group.Group("/tags")
	{
		tagsGroup.Get("/", r.getTags)
	}
}
