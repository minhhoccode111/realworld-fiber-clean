package v1

import (
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/controller/http/v1/request"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/controller/http/v1/response"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/entity"
	"github.com/minhhoccode111/realworld-fiber-clean/pkg/util"
)

// @Summary     Register User
// @Description Register User
// @ID          users-register
// @Tags  	    users
// @Accept      json
// @Produce     json
// @Param       request body request.UserRegisterRequest true "Register User"
// @Success     200 {object} response.UserAuth
// @Failure     400 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /users [post]
func (r *V1) postRegisterUser(ctx *fiber.Ctx) error {
	var body request.UserRegisterRequest

	if err := ctx.BodyParser(&body); err != nil {
		r.l.Error(err, "http - v1 - postRegisterUser - ctx.BodyParser")

		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	if err := r.v.Struct(body.User); err != nil {
		r.l.Error(err, "http - v1 - postRegisterUser - r.v.Struct")
		if verrs, ok := err.(validator.ValidationErrors); ok {
			errors := make([]string, 0, len(verrs))
			for _, e := range verrs {
				switch e.Tag() {
				case "required":
					errors = append(errors, e.Field()+" is required")
				case "email":
					errors = append(errors, "invalid email format")
				case "passwd":
					errors = append(
						errors,
						"password must include upper, lower, digit, and special char",
					)
				case "username":
					errors = append(
						errors,
						"username can only contain letters, numbers, and underscore",
					)
				default:
					errors = append(errors, e.Field()+" is invalid")
				}
			}
			return errorResponse(ctx, http.StatusBadRequest, strings.Join(errors, ", "))
		}
		return errorResponse(ctx, http.StatusInternalServerError, "validation error")
	}

	e := entity.User{
		Username: body.User.Username,
		Email:    body.User.Email,
		Password: body.User.Password,
	}

	// user.Id generated, user.Password hashed
	user, err := r.u.RegisterUser(ctx.UserContext(), e)
	if err != nil {
		r.l.Error(err, "http - v1 - postRegisterUser - r.u.RegisterUser")

		return errorResponse(ctx, http.StatusInternalServerError, "database problems")
	}

	token, err := util.GenerateJWT(
		user.Id,
		r.cfg.JWT.Secret,
		r.cfg.JWT.Issuer,
		r.cfg.JWT.Expiration,
	)
	if err != nil {
		r.l.Error(err, "http - v1 - postRegisterUser - util.GenerateJWT")

		return errorResponse(ctx, http.StatusInternalServerError, "generate jwt error")
	}

	userAuthResponse := response.UserAuthResponse{
		User: response.NewUserAuth(user, token),
	}

	return ctx.Status(http.StatusOK).JSON(userAuthResponse)
}

// @Summary     Login User
// @Description Login User
// @ID          users-login
// @Tags  	    users
// @Accept      json
// @Produce     json
// @Success     200 {object} response.UserAuth
// @Failure     400 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /users/login [post]
func (r *V1) postLoginUser(ctx *fiber.Ctx) error {
	return nil
}
