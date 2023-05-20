package user

import (
	"github.com/gin-gonic/gin"
	"github.com/rickyromansyah2045/halocat-backend-go/helper"
	"gorm.io/gorm"
)

type Repository interface {
	SaveUser(user User) (User, error)
	GetUserByEmail(email string) (User, error)
	GetUserByID(id int) (User, error)

	GetAllUser() ([]User, error)
	UpdateUser(User) (User, error)
	DeleteUser(User) (bool, error)

	GetDataForgotPasswordByToken(token string) (UserForgotPasswordToken, error)
	CreateForgotPasswordToken(UserForgotPasswordToken) (UserForgotPasswordToken, error)
	DeleteForgotPasswordToken(UserForgotPasswordToken) (bool, error)

	AdminDataTablesUsers(ctx *gin.Context) (helper.DataTables, error)

	GetUserRegistered(condition string) (int, error)
}

type repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{DB: db}
}
