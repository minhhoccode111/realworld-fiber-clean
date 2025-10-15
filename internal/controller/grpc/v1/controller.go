package v1

import (
	v1 "github.com/minhhoccode111/realworld-fiber-clean/docs/proto/v1"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/usecase"
	"github.com/minhhoccode111/realworld-fiber-clean/pkg/logger"
	"github.com/go-playground/validator/v10"
)

// V1 -.
type V1 struct {
	v1.TranslationServer

	t usecase.Translation
	l logger.Interface
	v *validator.Validate
}
