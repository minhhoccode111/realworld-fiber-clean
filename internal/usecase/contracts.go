// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"github.com/minhhoccode111/realworld-fiber-clean/internal/entity"
)

//go:generate mockgen -source=contracts.go -destination=./mocks_usecase_test.go -package=usecase_test

type (
	// Translation defines translation-related business actions.
	Translation interface {
		// Translate performs translation via a web API then persists the result.
		Translate(context.Context, entity.Translation) (entity.Translation, error)
		// History retrieves previously translated entries.
		History(context.Context) (entity.TranslationHistory, error)
	}

	// User encapsulates business actions for user accounts.
	User interface {
		// Register creates a new user after hashing credentials.
		Register(context.Context, *entity.User) error
		// Login authenticates a user by email and password.
		Login(context.Context, *entity.User) (*entity.User, error)
		// Current fetches user details by identifier.
		Current(context.Context, string) (*entity.User, error)
		// Update mutates user profile fields and persists changes.
		Update(context.Context, *entity.User) (*entity.User, error)
	}

	// Article defines business operations for articles.
	Article interface {
		// Create stores a new article and returns its detail with tags.
		Create(context.Context, *entity.Article, []string) (*entity.ArticleDetail, error)
		// List returns article previews with pagination and filter options.
		List(
			ctx context.Context,
			isFeed bool,
			userID, tag, author, favorited string,
			limit, offset uint64,
		) ([]entity.ArticlePreview, uint64, error)
		// Detail returns full detail for a given slug and viewer.
		Detail(ctx context.Context, userID, slug string) (*entity.ArticleDetail, error)
		// Update modifies an article identified by slug.
		Update(
			ctx context.Context,
			userID, oldSlug string,
			dto *entity.Article,
		) (*entity.ArticleDetail, error)
		// Delete removes an article when permitted by ownership or admin role.
		Delete(ctx context.Context, userID, slug string, userRole entity.Role) error
	}

	// Favorite defines operations to toggle article favorites.
	Favorite interface {
		// Create marks an article as favorited by a user.
		Create(ctx context.Context, userID, slug string) error
		// Delete removes a favorite relation.
		Delete(ctx context.Context, userID, slug string) error
	}

	// Comment defines operations for managing article comments.
	Comment interface {
		// Create adds a comment to an article.
		Create(ctx context.Context, slug string, dto *entity.Comment) (*entity.CommentDetail, error)
		// List returns comments for an article with pagination.
		List(
			ctx context.Context,
			userID, slug string,
			limit, offset uint64,
		) ([]entity.CommentDetail, uint64, error)
		// Delete removes a comment if authorized.
		Delete(ctx context.Context, userID, slug, commentID string, userRole entity.Role) error
	}

	// Profile defines operations around user profile views and follows.
	Profile interface {
		// Detail returns a profile preview with follow status.
		Detail(ctx context.Context, userID, username string) (*entity.ProfilePreview, error)
		// Follow creates a following relationship.
		Follow(ctx context.Context, userID, username string) error
		// Unfollow removes a following relationship.
		Unfollow(ctx context.Context, userID, username string) error
	}

	// Tag defines operations for listing tags.
	Tag interface {
		// List returns tags with pagination.
		List(
			ctx context.Context,
			limit, offset uint64,
		) (tags []entity.TagName, total uint64, err error)
	}
)
