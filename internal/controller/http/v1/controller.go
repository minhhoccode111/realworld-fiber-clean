package v1

import (
	"github.com/go-playground/validator/v10"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/usecase"
	"github.com/minhhoccode111/realworld-fiber-clean/pkg/logger"
)

// V1 -.
type V1 struct {
	t   usecase.Translation
	tc  usecase.TranslationClone
	tag usecase.Tag
	l   logger.Interface
	v   *validator.Validate
}
