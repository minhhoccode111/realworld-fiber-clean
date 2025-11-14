package v1

import (
	"github.com/go-playground/validator/v10"
	"github.com/minhhoccode111/realworld-fiber-clean/config"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/usecase"
	"github.com/minhhoccode111/realworld-fiber-clean/pkg/logger"
)

// V1 -.
type V1 struct {
	l   logger.Interface
	v   *validator.Validate
	cfg *config.Config

	t   usecase.Translation
	u   usecase.User
	a   usecase.Article
	f   usecase.Favorite
	c   usecase.Comment
	p   usecase.Profile
	tag usecase.Tag
}
