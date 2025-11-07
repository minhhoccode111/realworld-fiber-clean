package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/minhhoccode111/realworld-fiber-clean/config"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/controller/http/middleware"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/usecase"
	"github.com/minhhoccode111/realworld-fiber-clean/pkg/logger"
	"github.com/minhhoccode111/realworld-fiber-clean/pkg/validatorx"
)

// NewV1Routes -.
func NewV1Routes(
	apiV1Group fiber.Router,
	cfg *config.Config,
	l logger.Interface,

	t usecase.Translation,
	tc usecase.TranslationClone,
	u usecase.User,
	a usecase.Article,
	f usecase.Favorite,
	c usecase.Comment,
	p usecase.Profile,
	tag usecase.Tag,
) {
	r := &V1{
		cfg: cfg,
		l:   l,
		v:   validatorx.New(),

		t:   t,
		tc:  tc,
		u:   u,
		a:   a,
		f:   f,
		c:   c,
		p:   p,
		tag: tag,
	}

	jwtSecret := r.cfg.JWT.Secret
	auth := middleware.AuthMiddleware(l, jwtSecret, false)
	optionalAuth := middleware.AuthMiddleware(l, jwtSecret, true)

	translation := apiV1Group.Group("/translation")
	{
		translation.Get("/history", r.history)
		translation.Post("/do-translate", r.doTranslate)
	}

	translationClone := apiV1Group.Group("/translation-clone")
	{
		translationClone.Get("/history", r.getHistory)
		translationClone.Post("/translate", r.postTranslate)
	}

	users := apiV1Group.Group("/users")
	{
		users.Post("/", r.postRegisterUser)
		users.Post("/login", r.postLoginUser)
		users.Post("/logout", r.postLogoutUser)
	}

	user := apiV1Group.Group("/user")
	{
		user.Get("/", auth, r.getCurrentUser)
		user.Put("/", auth, r.putUpdateUser)
	}

	articles := apiV1Group.Group("/articles")
	{
		articles.Post("/", auth, r.postArticle)
		articles.Get("/", optionalAuth, r.getAllArticles)
		articles.Get("/feed", auth, r.getFeedArticles)
		articles.Get("/:slug", optionalAuth, r.getArticle)
		articles.Put("/:slug", auth, r.putArticle)
		articles.Delete("/:slug", auth, r.deleteArticle)
	}

	favorites := apiV1Group.Group("/articles/:slug/favorite")
	{
		favorites.Post("/", auth, r.createFavorite)
		favorites.Delete("/", auth, r.deleteFavorite)
	}

	comments := apiV1Group.Group("/articles/:slug/comments")
	{
		comments.Post("/", auth, r.postComment)
		comments.Get("/", optionalAuth, r.getAllComments)
		comments.Delete("/:commentId", auth, r.deleteComment)
	}

	profiles := apiV1Group.Group("/profiles/:username")
	{
		profiles.Get("/", optionalAuth, r.getProfile)
		profiles.Post("/follow", auth, r.postFollowProfile)
		profiles.Delete("/follow", auth, r.deleteFollowProfile)
	}

	tags := apiV1Group.Group("/tags")
	{
		tags.Get("/", r.getTags)
	}
}
