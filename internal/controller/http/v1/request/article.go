package request

import (
	"strings"
)

type ArticleCreate struct {
	Title       string   `json:"title"       validate:"required,max=255"         example:"This is title - to generate slug"`
	Description string   `json:"description" validate:"required,max=255"         example:"This is description"`
	Body        string   `json:"body"        validate:"required,max=50000"       example:"this is article content"`
	TagList     []string `json:"tagList"     validate:"dive,required,max=50,tag" example:"go,fiber,api,clean-arch"`
}

func (a *ArticleCreate) Trim() {
	a.Title = strings.TrimSpace(a.Title)
	a.Description = strings.TrimSpace(a.Description)
	a.Body = strings.TrimSpace(a.Body)
	cleaned := []string{}

	for _, v := range a.TagList {
		v = strings.TrimSpace(v)
		if v != "" {
			cleaned = append(cleaned, v)
		}
	}

	a.TagList = cleaned
}

type ArticleCreateRequest struct {
	Article ArticleCreate `json:"article"`
}

type ArticleUpdate struct {
	Title       string `json:"title"       validate:"omitempty,max=255"   example:"This is title - to generate slug"`
	Description string `json:"description" validate:"omitempty,max=255"   example:"This is description"`
	Body        string `json:"body"        validate:"omitempty,max=50000" example:"this is article content"`
}

func (a *ArticleUpdate) Trim() {
	a.Title = strings.TrimSpace(a.Title)
	a.Description = strings.TrimSpace(a.Description)
	a.Body = strings.TrimSpace(a.Body)
}

type ArticleUpdateRequest struct {
	Article ArticleUpdate `json:"article"`
}
