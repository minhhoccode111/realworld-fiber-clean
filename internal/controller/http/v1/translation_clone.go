package v1

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/controller/http/v1/request"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/entity"
)

// @Summary     Show history
// @Description Show all translation history
// @ID          history-clone
// @Tags  	    translation-clone
// @Accept      json
// @Produce     json
// @Param       limit  query int false "Number of items to return"  minimum(1)  default(10)
// @Param       offset query int false "Number of items to skip"    minimum(0)  default(0)
// @Success     200 {object} entity.TranslationCloneHistory
// @Failure     500 {object} response.Error
// @Router      /translation-clone/history [get]
func (r *V1) getHistory(ctx *fiber.Ctx) error {
	limit, err := strconv.ParseUint(ctx.Query("limit", "10"), 10, 64)
	if err != nil {
		limit = 10
	}

	offset, err := strconv.ParseUint(ctx.Query("offset", "0"), 10, 64)
	if err != nil {
		offset = 0
	}

	translationCloneHistory, err := r.tc.GetHistory(ctx.UserContext(), limit, offset)
	if err != nil {
		r.l.Error(err, "http - v1 - getHistory")

		return errorResponse(ctx, http.StatusInternalServerError, "database problems")
	}

	return ctx.Status(http.StatusOK).JSON(translationCloneHistory)
}

// @Summary     Do Translate
// @Description Translate a text
// @ID          translate-clone
// @Tags  	    translation-clone
// @Accept      json
// @Produce     json
// @Param       request body request.TranslateClone true "Set up translation"
// @Success     200 {object} entity.TranslationClone
// @Failure     400 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /translation-clone/translate [post]
func (r *V1) postTranslate(ctx *fiber.Ctx) error {
	var body request.TranslateClone

	if err := ctx.BodyParser(&body); err != nil {
		r.l.Error(err, "http - v1 - postTranslate")

		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	if err := r.v.Struct(body); err != nil {
		r.l.Error(err, "http - v1 - postTranslate")

		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	translation, err := r.tc.PostTranslate(ctx.UserContext(), entity.TranslationClone{
		Source:      body.Source,
		Destination: body.Destination,
		Original:    body.Original,
	})
	if err != nil {
		r.l.Error(err, "http - v1 - postTranslate")

		return errorResponse(ctx, http.StatusInternalServerError, "translation service problems")
	}

	return ctx.Status(http.StatusOK).JSON(translation)
}
