package content

import (
	"github.com/gin-gonic/gin"
	"github.com/rickyromansyah2045/halocat-backend-go/helper"
	"github.com/rickyromansyah2045/halocat-backend-go/user"
	"gorm.io/gorm"
)

type Repository interface {
	GetAllContent(ctx *gin.Context) ([]Content, error)
	GetContentByID(id int) (Content, error)
	SaveContent(Content) (Content, error)
	UpdateContent(Content) (Content, error)
	DeleteContent(Content) (bool, error)

	GetAllContentImage() ([]ContentImage, error)
	GetContentImageByID(id int) (ContentImage, error)
	CreateContentImage(ContentImage) (ContentImage, error)
	UpdateAllImagesAsNonPrimary(contentID int) (bool, error)
	DeleteContentImage(ContentImage) (bool, error)

	GetAllContentCategory() ([]ContentCategory, error)
	GetContentCategoryByID(id int) (ContentCategory, error)
	SaveContentCategory(ContentCategory) (ContentCategory, error)
	UpdateContentCategory(ContentCategory) (ContentCategory, error)
	DeleteContentCategory(ContentCategory) (bool, error)

	AdminDataTablesContents(ctx *gin.Context) (helper.DataTables, error)
	AdminDataTablesCategories(ctx *gin.Context) (helper.DataTables, error)

	UserDataTablesContents(*gin.Context, user.User) (helper.DataTables, error)
}

type repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{DB: db}
}
