package webapi

import (
	"fmt"

	translator "github.com/Conight/go-googletrans"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/entity"
)

// TranslationWebAPI -.
type TranslationCloneWebAPI struct {
	conf translator.Config
}

// New -.
func NewClone() *TranslationCloneWebAPI {
	conf := translator.Config{
		UserAgent: []string{
			"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:15.0) Gecko/20100101 Firefox/15.0.1",
		},
		ServiceUrls: []string{"translate.google.com"},
	}

	return &TranslationCloneWebAPI{
		conf: conf,
	}
}

// Translate -.
func (t *TranslationCloneWebAPI) DoTranslate(
	translation entity.TranslationClone,
) (entity.TranslationClone, error) {
	trans := translator.New(t.conf)

	result, err := trans.Translate(
		translation.Original,
		translation.Source,
		translation.Destination,
	)
	if err != nil {
		return entity.TranslationClone{}, fmt.Errorf(
			"TranslationCloneWebAPI - DoTranslate - trans.Translate: %w",
			err,
		)
	}

	translation.Translation = result.Text

	return translation, nil
}
