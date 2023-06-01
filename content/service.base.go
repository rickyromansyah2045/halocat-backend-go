package content

import (
	"github.com/gin-gonic/gin"
	"github.com/rickyromansyah2045/halocat-backend-go/helper"
	"github.com/rickyromansyah2045/halocat-backend-go/user"
)

type Service interface {
	GetAllContent(ctx *gin.Context) ([]Content, error)
	GetContentByID(RequestGetContentByID) (Content, error)
	CreateContent(RequestCreateContent) (Content, error)
	UpdateContent(RequestGetContentByID, RequestUpdateContent) (Content, error)
	DeleteContent(RequestGetContentByID, RequestDeleteContent) (bool, error)

	GetAllContentImage() ([]ContentImage, error)
	GetContentImageByID(RequestGetContentImageByID) (ContentImage, error)
	SaveContentImage(RequestCreateContentImage, string) (ContentImage, error)
	DeleteContentImage(RequestGetContentImageByID, RequestDeleteContentImage) (bool, error)

	GetAllContentCategory() ([]ContentCategory, error)
	GetContentCategoryByID(RequestGetContentCategoryByID) (ContentCategory, error)
	CreateContentCategory(RequestCreateContentCategory) (ContentCategory, error)
	UpdateContentCategory(RequestGetContentCategoryByID, RequestUpdateContentCategory) (ContentCategory, error)
	DeleteContentCategory(RequestGetContentCategoryByID, RequestDeleteContentCategory) (bool, error)

	AdminDataTablesContents(*gin.Context) (helper.DataTables, error)
	AdminDataTablesCategories(*gin.Context) (helper.DataTables, error)

	UserDataTablesContents(*gin.Context, user.User) (helper.DataTables, error)
}

type service struct {
	repo     Repository
	userRepo user.Repository
}

func NewService(
	repository Repository,
	userRepository user.Repository,
) *service {
	return &service{
		repo:     repository,
		userRepo: userRepository,
	}
}
