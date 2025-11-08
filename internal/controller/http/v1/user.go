package v1

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/controller/http/middleware"
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
// @Success     201 {object} response.UserAuthResponse
// @Failure     400 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /users [post]
func (r *V1) postRegisterUser(ctx *fiber.Ctx) error {
	var body request.UserRegisterRequest

	if err := ctx.BodyParser(&body); err != nil {
		r.l.Error(err, "http - v1 - postRegisterUser - ctx.BodyParser")

		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	body.User.Trim()

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
				case "password":
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
			return errorResponse(ctx, http.StatusBadRequest, strings.Join(errors, "; "))
		}
		return errorResponse(ctx, http.StatusInternalServerError, "validation error")
	}

	// user.Id generated, user.Password hashed
	user, err := r.u.Register(ctx.UserContext(), entity.User{
		Username: body.User.Username,
		Email:    body.User.Email,
		Password: body.User.Password,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return errorResponse(ctx, http.StatusConflict, "email/username already existed")
			}
		}

		r.l.Error(err, "http - v1 - postRegisterUser - r.u.Register")

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

	return ctx.Status(http.StatusCreated).JSON(response.UserAuthResponse{
		User: response.NewUserAuth(user, token),
	})
}

// @Summary     Login User
// @Description Login User
// @ID          users-login
// @Tags  	    users
// @Accept      json
// @Param       request body request.UserLoginRequest true "Login User"
// @Produce     json
// @Success     200 {object} response.UserAuthResponse
// @Failure     400 {object} response.Error
// @Failure     401 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /users/login [post]
func (r *V1) postLoginUser(ctx *fiber.Ctx) error {
	var body request.UserLoginRequest

	if err := ctx.BodyParser(&body); err != nil {
		r.l.Error(err, "http - v1 - postLoginUser - ctx.BodyParser")

		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	body.User.Trim()

	if err := r.v.Struct(body.User); err != nil {
		r.l.Error(err, "http - v1 - postLoginUser - r.v.Struct")

		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	user, err := r.u.Login(ctx.UserContext(), entity.User{
		Email:    body.User.Email,
		Password: body.User.Password,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errorResponse(ctx, http.StatusUnauthorized, "incorrect email")
		}

		if strings.Contains(err.Error(), "incorrect password") {
			return errorResponse(ctx, http.StatusUnauthorized, "incorrect password")
		}

		r.l.Error(err, "http - v1 - postLoginUser - r.u.Login")

		return errorResponse(ctx, http.StatusInternalServerError, "database problems")
	}

	token, err := util.GenerateJWT(
		user.Id,
		r.cfg.JWT.Secret,
		r.cfg.JWT.Issuer,
		r.cfg.JWT.Expiration,
	)
	if err != nil {
		r.l.Error(err, "http - v1 - postLoginUser - util.GenerateJWT")

		return errorResponse(ctx, http.StatusInternalServerError, "jwt problems")
	}

	return ctx.Status(http.StatusOK).JSON(response.UserAuthResponse{
		User: response.NewUserAuth(user, token),
	})
}

func (r *V1) postLogoutUser(ctx *fiber.Ctx) error {
	return ctx.SendStatus(200)
}

// @Summary     Get current User
// @Description Get current User
// @ID          users-current
// @Tags  	    users
// @Produce     json
// @Success     200 {object} response.UserAuthResponse
// @Failure     500 {object} response.Error
// @Router      /user [get]
// @Security    BearerAuth
func (r *V1) getCurrentUser(ctx *fiber.Ctx) error {
	userId := ctx.Locals(middleware.CtxUserIdKey).(string)
	if userId == "" {
		return errorResponse(ctx, http.StatusInternalServerError, "cannot authorize user in jwt")
	}

	user, err := r.u.Current(ctx.UserContext(), userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errorResponse(ctx, http.StatusNotFound, "userId in token not found")
		}

		r.l.Error(err, "http - v1 - getCurrentUser - r.u.Current")

		return errorResponse(ctx, http.StatusInternalServerError, "database problems")
	}

	token, err := util.GenerateJWT(
		userId,
		r.cfg.JWT.Secret,
		r.cfg.JWT.Issuer,
		r.cfg.JWT.Expiration,
	)
	if err != nil {
		r.l.Error(err, "http - v1 - getCurrentUser - util.GenerateJWT")

		return errorResponse(ctx, http.StatusInternalServerError, "jwt problems")
	}

	return ctx.Status(http.StatusOK).JSON(response.UserAuthResponse{
		User: response.NewUserAuth(user, token),
	})
}

// @Summary     Update User
// @Description Update User
// @ID          users-update
// @Tags  	    users
// @Accept      json
// @Produce     json
// @Param       request body request.UserUpdateRequest true "Update User"
// @Success     200 {object} response.UserAuthResponse
// @Failure     400 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /user [put]
// @Security    BearerAuth
func (r *V1) putUpdateUser(ctx *fiber.Ctx) error {
	var body request.UserUpdateRequest

	if err := ctx.BodyParser(&body); err != nil {
		r.l.Error(err, "http - v1 - putUpdateUser - ctx.BodyParser")

		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	body.User.Trim()

	if body.User.IsAllEmpty() {
		return errorResponse(ctx, http.StatusBadRequest, "no field provided")
	}

	if err := r.v.Struct(body.User); err != nil {
		r.l.Error(err, "http - v1 - putUpdateUser - r.v.Struct")
		if verrs, ok := err.(validator.ValidationErrors); ok {
			errors := make([]string, 0, len(verrs))
			for _, e := range verrs {
				switch e.Tag() {
				case "email":
					errors = append(errors, "invalid email format")
				case "password":
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
			return errorResponse(ctx, http.StatusBadRequest, strings.Join(errors, "; "))
		}
		return errorResponse(ctx, http.StatusInternalServerError, "validation error")
	}

	userId := ctx.Locals(middleware.CtxUserIdKey).(string)
	if userId == "" {
		return errorResponse(ctx, http.StatusInternalServerError, "cannot authorize user in jwt")
	}

	user, err := r.u.Update(ctx.UserContext(), body.User.NewUser(userId))
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return errorResponse(ctx, http.StatusConflict, "email/username alread existed")
			}
		}

		r.l.Error(err, "http - v1 - putUpdateUser - r.u.Update")

		return errorResponse(ctx, http.StatusInternalServerError, "database problems")
	}

	token, err := util.GenerateJWT(
		userId,
		r.cfg.JWT.Secret,
		r.cfg.JWT.Issuer,
		r.cfg.JWT.Expiration,
	)
	if err != nil {
		r.l.Error(err, "http - v1 - putUpdateUser - util.GenerateJWT")

		return errorResponse(ctx, http.StatusInternalServerError, "generate jwt error")
	}

	return ctx.Status(http.StatusOK).JSON(response.UserAuthResponse{
		User: response.NewUserAuth(user, token),
	})
}
