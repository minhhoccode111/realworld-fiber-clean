// Package repo implements application outer layer logic. Each logic group in own file.
package repo

import (
	"context"

	"github.com/minhhoccode111/realworld-fiber-clean/internal/entity"
)

//go:generate mockgen -source=contracts.go -destination=../usecase/mocks_repo_test.go -package=usecase_test

type (
	// TranslationRepo -.
	TranslationRepo interface {
		Store(context.Context, entity.Translation) error
		GetHistory(context.Context) ([]entity.Translation, error)
	}

	// TranslationWebAPI -.
	TranslationWebAPI interface {
		Translate(entity.Translation) (entity.Translation, error)
	}

	// UserRepo -.
	UserRepo interface {
		StoreRegister(context.Context, *entity.User) error
		GetUserByEmail(context.Context, string) (*entity.User, error)
		GetUserByID(context.Context, string) (*entity.User, error)
		StoreUpdate(context.Context, *entity.User) error
	}

	// ArticleRepo -.
	ArticleRepo interface {
		StoreCreate(ctx context.Context, dto *entity.Article, tags []string) (err error)
		CanSlugBeUsed(ctx context.Context, articleID, slug string) (bool, error)
		GetDetailBySlug(ctx context.Context, userID, slug string) (*entity.ArticleDetail, error)
		StoreTagsList(ctx context.Context, tags []string) (ids []string, err error)
		StoreArticleTagsList(ctx context.Context, articleID string, tagIDs []string) error
		GetList(
			ctx context.Context,
			isFeed bool,
			userID, tag, author, favorited string,
			limit, offset uint64,
		) ([]entity.ArticlePreview, uint64, error)
		GetBasicBySlug(ctx context.Context, slug string) (*entity.Article, error)
		StoreUpdate(ctx context.Context, dto *entity.Article) (err error)
		StoreDelete(ctx context.Context, userID, slug string) (err error)
	}

	FavoriteRepo interface {
		StoreCreate(ctx context.Context, userID, slug string) error
		StoreDelete(ctx context.Context, userID, slug string) error
	}

	CommentRepo interface {
		StoreCreate(ctx context.Context, slug string, dto *entity.Comment) (string, error)
		GetDetailByID(ctx context.Context, userID, commentID string) (*entity.CommentDetail, error)
		GetList(
			ctx context.Context,
			userID, slug string,
			limit, offset uint64,
		) ([]entity.CommentDetail, uint64, error)
		StoreDelete(ctx context.Context, userID, slug, commentID string) (err error)
	}

	ProfileRepo interface {
		IsExisted(ctx context.Context, username string) error
		GetDetail(ctx context.Context, userID, username string) (*entity.ProfilePreview, error)
		StoreCreate(ctx context.Context, userID, username string) error
		StoreDelete(ctx context.Context, userID, username string) error
	}

	// TagRepo -.
	TagRepo interface {
		GetList(
			ctx context.Context,
			limit, offset uint64,
		) (tags []entity.TagName, total uint64, err error)
	}
)
