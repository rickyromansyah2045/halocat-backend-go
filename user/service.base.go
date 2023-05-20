package user

import (
	"github.com/gin-gonic/gin"
	"github.com/rickyromansyah2045/halocat-backend-go/helper"
)

type Service interface {
	Register(req RequestRegister) (User, error)
	Login(req RequestLogin) (User, error)
	GetUserByID(ID int) (User, error)
	GetUserByEmail(Email string) (User, error)
	CheckDuplicateEmail(email string) (duplicate bool, err error)

	GetAllUser() ([]User, error)
	CreateUser(RequestCreateUser) (User, error)
	UpdateUser(RequestGetUserByID, RequestUpdateUser) (User, error)
	DeleteUser(RequestGetUserByID, RequestDeleteUser) (bool, error)

	GetDataForgotPasswordByToken(token string) (UserForgotPasswordToken, error)
	CreateUserForgotPasswordToken(RequestCreateForgotPasswordToken) (UserForgotPasswordToken, error)
	DeleteForgotPasswordToken(UserForgotPasswordToken) (bool, error)

	AdminDataTablesUsers(*gin.Context) (helper.DataTables, error)

	GetUserRegistered(condition string) (int, error)
}

type service struct {
	repo Repository
}

func NewService(
	repository Repository,
) *service {
	return &service{
		repo: repository,
	}
}
