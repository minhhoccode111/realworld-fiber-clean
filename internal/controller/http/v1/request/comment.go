package request

import "strings"

type CommentCreate struct {
	Body string `json:"body" validate:"required,max=10000" example:"this is a comment"`
}

func (c *CommentCreate) Trim() {
	c.Body = strings.TrimSpace(c.Body)
}

type CommentCreateRequest struct {
	Comment CommentCreate `json:"comment"`
}
