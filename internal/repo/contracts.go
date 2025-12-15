// Package repo implements application outer layer logic. Each logic group in own file.
package repo

import (
	"context"

	"github.com/minhhoccode111/realworld-fiber-clean/internal/entity"
)

//go:generate mockgen -source=contracts.go -destination=../usecase/mocks_repo_test.go -package=usecase_test

type (
	// TranslationRepo defines persistence operations for translations.
	TranslationRepo interface {
		// Store saves a translation entry.
		Store(context.Context, entity.Translation) error
		// GetHistory fetches a list of stored translations.
		GetHistory(context.Context) ([]entity.Translation, error)
	}

	// TranslationWebAPI abstracts translation calls to an external service.
	TranslationWebAPI interface {
		// Translate performs translation via a remote provider.
		Translate(entity.Translation) (entity.Translation, error)
	}

	// UserRepo defines persistence operations for users.
	UserRepo interface {
		// StoreRegister persists a newly registered user.
		StoreRegister(context.Context, *entity.User) error
		// GetUserByEmail finds a user by email.
		GetUserByEmail(context.Context, string) (*entity.User, error)
		// GetUserByID finds a user by identifier.
		GetUserByID(context.Context, string) (*entity.User, error)
		// StoreUpdate writes back changes to a user.
		StoreUpdate(context.Context, *entity.User) error
	}

	// ArticleRepo defines persistence operations for articles and related data.
	ArticleRepo interface {
		// StoreCreate inserts an article and its tags.
		StoreCreate(ctx context.Context, dto *entity.Article, tags []string) (err error)
		// CanSlugBeUsed checks if a slug is available for use.
		CanSlugBeUsed(ctx context.Context, articleID, slug string) (bool, error)
		// GetDetailBySlug returns a fully populated article detail.
		GetDetailBySlug(ctx context.Context, userID, slug string) (*entity.ArticleDetail, error)
		// StoreTagsList inserts or updates tags and returns their IDs.
		StoreTagsList(ctx context.Context, tags []string) (ids []string, err error)
		// StoreArticleTagsList associates tags with an article.
		StoreArticleTagsList(ctx context.Context, articleID string, tagIDs []string) error
		// GetList returns article previews and total count with filters.
		GetList(
			ctx context.Context,
			isFeed bool,
			userID, tag, author, favorited string,
			limit, offset uint64,
		) ([]entity.ArticlePreview, uint64, error)
		// GetBasicBySlug returns minimal article data for edits or deletes.
		GetBasicBySlug(ctx context.Context, slug string) (*entity.Article, error)
		// StoreUpdate writes article changes to the store.
		StoreUpdate(ctx context.Context, dto *entity.Article) (err error)
		// StoreDelete marks an article as deleted.
		StoreDelete(ctx context.Context, slug string) (err error)
	}

	// FavoriteRepo defines persistence operations for article favorites.
	FavoriteRepo interface {
		// StoreCreate creates a favorite relation.
		StoreCreate(ctx context.Context, userID, slug string) error
		// StoreDelete removes a favorite relation.
		StoreDelete(ctx context.Context, userID, slug string) error
	}

	// CommentRepo defines persistence operations for comments.
	CommentRepo interface {
		// StoreCreate inserts a comment and returns its ID.
		StoreCreate(ctx context.Context, slug string, dto *entity.Comment) (string, error)
		// GetDetailByID returns detailed comment information.
		GetDetailByID(ctx context.Context, userID, commentID string) (*entity.CommentDetail, error)
		// GetList returns comments for an article and the total count.
		GetList(
			ctx context.Context,
			userID, slug string,
			limit, offset uint64,
		) ([]entity.CommentDetail, uint64, error)
		// GetBasicByID returns minimal comment data for authorization checks.
		GetBasicByID(ctx context.Context, commentID string) (*entity.Comment, error)
		// StoreDelete marks a comment as deleted if it belongs to the article.
		StoreDelete(ctx context.Context, slug, commentID string) (err error)
	}

	// ProfileRepo defines persistence operations for profiles and follows.
	ProfileRepo interface {
		// IsExisted reports whether a profile exists for the username.
		IsExisted(ctx context.Context, username string) error
		// GetDetail returns profile preview data.
		GetDetail(ctx context.Context, userID, username string) (*entity.ProfilePreview, error)
		// StoreCreate creates a follow relation.
		StoreCreate(ctx context.Context, userID, username string) error
		// StoreDelete removes a follow relation.
		StoreDelete(ctx context.Context, userID, username string) error
	}

	// TagRepo defines persistence operations for tags.
	TagRepo interface {
		// GetList returns tags and total count with pagination.
		GetList(
			ctx context.Context,
			limit, offset uint64,
		) (tags []entity.TagName, total uint64, err error)
	}
)
