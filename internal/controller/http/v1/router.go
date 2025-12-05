package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/minhhoccode111/realworld-fiber-clean/config"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/controller/http/middleware"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/usecase"
	"github.com/minhhoccode111/realworld-fiber-clean/pkg/logger"
	"github.com/minhhoccode111/realworld-fiber-clean/pkg/validatorx"
)

// NewV1Routes -.
func NewV1Routes(
	apiV1Group *gin.RouterGroup,
	cfg *config.Config,
	l logger.Interface,

	t usecase.Translation,
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
		translation.GET("/history", r.history)
		translation.POST("/do-translate", r.doTranslate)
	}

	users := apiV1Group.Group("/users")
	{
		users.POST("/", r.postRegisterUser)
		users.POST("/login", r.postLoginUser)
		// users.Post("/logout", r.postLogoutUser)
	}

	user := apiV1Group.Group("/user")
	{
		user.GET("/", auth, r.getCurrentUser)
		user.PUT("/", auth, r.putUpdateUser)
	}

	articles := apiV1Group.Group("/articles")
	{
		articles.POST("/", auth, r.postArticle)
		articles.GET("/", optionalAuth, r.getAllArticles)
		articles.GET("/feed", auth, r.getFeedArticles)
		articles.GET("/:slug", optionalAuth, r.getArticle)
		articles.PUT("/:slug", auth, r.putArticle)
		articles.DELETE("/:slug", auth, r.deleteArticle)
	}

	favorites := apiV1Group.Group("/articles/:slug/favorite")
	{
		favorites.POST("/", auth, r.createFavorite)
		favorites.DELETE("/", auth, r.deleteFavorite)
	}

	comments := apiV1Group.Group("/articles/:slug/comments")
	{
		comments.POST("/", auth, r.postComment)
		comments.GET("/", optionalAuth, r.getAllComments)
		comments.DELETE("/:commentID", auth, r.deleteComment)
	}

	profiles := apiV1Group.Group("/profiles/:username")
	{
		profiles.GET("/", optionalAuth, r.getProfile)
		profiles.POST("/follow", auth, r.postFollowProfile)
		profiles.DELETE("/follow", auth, r.deleteFollowProfile)
	}

	tags := apiV1Group.Group("/tags")
	{
		tags.GET("/", r.getTags)
	}
}
