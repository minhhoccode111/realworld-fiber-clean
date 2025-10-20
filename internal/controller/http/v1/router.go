package v1

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/usecase"
	"github.com/minhhoccode111/realworld-fiber-clean/pkg/logger"
)

// NewTranslationRoutes -.
func NewTranslationRoutes(
	apiV1Group fiber.Router,
	l logger.Interface,
	t usecase.Translation,
	tc usecase.TranslationClone,
) {
	r := &V1{t: t, tc: tc, l: l, v: validator.New(validator.WithRequiredStructEnabled())}

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
}
