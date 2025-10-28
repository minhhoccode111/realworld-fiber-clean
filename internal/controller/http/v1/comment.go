package v1

import (
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/controller/http/middleware"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/controller/http/v1/request"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/controller/http/v1/response"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/entity"
)

// @Summary     Create Comment
// @Description Create Comment
// @ID          comments-create
// @Tags  	    comments
// @Accept      json
// @Produce     json
// @Param       request body request.CommentCreateRequest true "Create Comment"
// @Success     201 {object} response.CommentDetailResponse
// @Failure     400 {object} response.Error
// @Failure     401 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /articles/{slug}/comments [post]
// @Security    BearerAuth
func (r *V1) postCreateComment(ctx *fiber.Ctx) error {
	var body request.CommentCreateRequest

	if err := ctx.BodyParser(&body); err != nil {
		r.l.Error(err, "http - v1 - postCreateComment - ctx.BodyParser")

		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	body.Comment.Trim()

	if err := r.v.Struct(body.Comment); err != nil {
		r.l.Error(err, "http - v1 - postCreateComment - r.v.Struct")
		if verrs, ok := err.(validator.ValidationErrors); ok {
			errors := make([]string, 0, len(verrs))
			for _, e := range verrs {
				switch e.Tag() {
				case "required":
					errors = append(errors, e.Field()+" is required")
				default:
					errors = append(errors, e.Field()+" is invalid")
				}
			}
			return errorResponse(ctx, http.StatusBadRequest, strings.Join(errors, "; "))
		}
		return errorResponse(ctx, http.StatusInternalServerError, "validation error")
	}

	userId := ctx.Locals(middleware.CtxUserIdKey).(string)
	if userId == "" {
		return errorResponse(ctx, http.StatusUnauthorized, "cannot authorize user in jwt")
	}

	slug := ctx.Params("slug")
	if slug == "" {
		return errorResponse(ctx, http.StatusBadRequest, "slug is required")
	}

	comment, err := r.c.Create(
		ctx.UserContext(),
		slug,
		entity.Comment{AuthorId: userId, Body: body.Comment.Body},
	)
	if err != nil {
		if strings.Contains(err.Error(), "notfound") {
			return errorResponse(ctx, http.StatusNotFound, "Article not found")
		}

		r.l.Error(err, "http - v1 - postCreateComment - r.c.Create")

		return errorResponse(ctx, http.StatusInternalServerError, "database problems")
	}

	return ctx.Status(http.StatusCreated).JSON(response.CommentDetailResponse{
		Comment: comment,
	})
}
