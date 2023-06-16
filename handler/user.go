package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rickyromansyah2045/halocat-backend-go/auth"
	"github.com/rickyromansyah2045/halocat-backend-go/helper"
	"github.com/rickyromansyah2045/halocat-backend-go/logs"
	"github.com/rickyromansyah2045/halocat-backend-go/user"
)

type userHandler struct {
	userSvc user.Service
	authSvc auth.Service
	logsSvc logs.Service
}

func NewUserHandler(
	userService user.Service,
	authService auth.Service,
	logsService logs.Service,
) *userHandler {
	return &userHandler{
		userSvc: userService,
		authSvc: authService,
		logsSvc: logsService,
	}
}

func (handler *userHandler) Register(ctx *gin.Context) {
	var req user.RequestRegister

	err := ctx.ShouldBind(&req)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Registration failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	isDuplicateEmail, err := handler.userSvc.CheckDuplicateEmail(req.Email)

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Registration failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if isDuplicateEmail {
		response := helper.APIResponseError(http.StatusConflict, "Registration failed!", "Email already registered!")
		ctx.JSON(http.StatusConflict, response)
		return
	}

	newUserData, err := handler.userSvc.Register(req)

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Registration failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	token, err := handler.authSvc.GenerateToken(newUserData.ID)

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Registration failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	userData := user.FormatUserData(newUserData, token)
	response := helper.APIResponse(http.StatusOK, "Registration successfully!", userData)

	// {
	// 	templateData := helper.EmailWelcome{
	// 		Name: userData.Name,
	// 	}
	// 	go helper.SendMail(userData.Email, "Welcome to The Cloud Donation", templateData, "html/welcome.html")
	// }

	handler.logsSvc.CreateActivityLog(ctx, fmt.Sprintf("%s registered to the system.", userData.Name))

	ctx.JSON(http.StatusOK, response)
}

func (handler *userHandler) Login(ctx *gin.Context) {
	var req user.RequestLogin

	err := ctx.ShouldBind(&req)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Login failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	userData, err := handler.userSvc.Login(req)

	if err != nil {
		if helper.IsErrNoRows(err.Error()) {
			response := helper.APIResponseError(http.StatusNotFound, "Login failed!", "User not found!")
			ctx.JSON(http.StatusNotFound, response)
			return
		}
		if err.Error() == "crypto/bcrypt: hashedPassword is not the hash of the given password" {
			response := helper.APIResponseError(http.StatusNotFound, "Login failed!", "Wrong password!")
			ctx.JSON(http.StatusNotFound, response)
			return
		}
		if err.Error() == "email not registered" {
			response := helper.APIResponseError(http.StatusNotFound, "Login failed!", err.Error())
			ctx.JSON(http.StatusNotFound, response)
			return
		}
		response := helper.APIResponseError(http.StatusInternalServerError, "Login failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	token, err := handler.authSvc.GenerateToken(userData.ID)

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Login failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	formatData := user.FormatUserData(userData, token)
	response := helper.APIResponse(http.StatusOK, "Login successfully!", formatData)

	handler.logsSvc.CreateActivityLog(ctx, fmt.Sprintf("%s successfully login to the system.", userData.Name))

	ctx.JSON(http.StatusOK, response)
}

func (handler *userHandler) GetAllUser(ctx *gin.Context) {
	users, err := handler.userSvc.GetAllUser()

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Get users failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	formatData := user.FormatListUserData(users)
	response := helper.APIResponse(http.StatusOK, "Get users successfully!", formatData)

	ctx.JSON(http.StatusOK, response)
}

func (handler *userHandler) GetUserByID(ctx *gin.Context) {
	var req user.RequestGetUserByID

	err := ctx.ShouldBindUri(&req)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Get user detail failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	userDetail, err := handler.userSvc.GetUserByID(req.ID)

	if err != nil {
		if helper.IsErrNoRows(err.Error()) {
			response := helper.APIResponseError(http.StatusNotFound, "Get user detail failed!", "Data not found!")
			ctx.JSON(http.StatusNotFound, response)
			return
		}

		response := helper.APIResponseError(http.StatusInternalServerError, "Get user detail failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	formatData := user.FormatUserFullData(userDetail)
	response := helper.APIResponse(http.StatusOK, "Get user detail successfully!", formatData)

	ctx.JSON(http.StatusOK, response)
}

func (handler *userHandler) CreateUser(ctx *gin.Context) {
	var req user.RequestCreateUser

	err := ctx.ShouldBind(&req)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Create user failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	isDuplicateEmail, err := handler.userSvc.CheckDuplicateEmail(req.Email)

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Create user failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if isDuplicateEmail {
		response := helper.APIResponseError(http.StatusConflict, "Create user failed!", "Email already registered!")
		ctx.JSON(http.StatusConflict, response)
		return
	}

	req.User = ctx.MustGet("userData").(user.User)

	newUserData, err := handler.userSvc.CreateUser(req)

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Create user failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	formatData := user.FormatUserFullData(newUserData)
	response := helper.APIResponse(http.StatusCreated, "Create user successfully!", formatData)

	handler.logsSvc.CreateActivityLog(ctx, fmt.Sprintf("%v creating user id %v.", req.User.Name, newUserData.ID))

	ctx.JSON(http.StatusCreated, response)
}

func (handler *userHandler) UpdateUser(ctx *gin.Context) {
	var reqID user.RequestGetUserByID

	err := ctx.ShouldBindUri(&reqID)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Update user failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	var reqUpdate user.RequestUpdateUser

	err = ctx.ShouldBind(&reqUpdate)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Update user failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	reqUpdate.User = ctx.MustGet("userData").(user.User)

	updatedUser, err := handler.userSvc.UpdateUser(reqID, reqUpdate)

	if err != nil {
		if helper.IsErrNoRows(err.Error()) {
			response := helper.APIResponseError(http.StatusNotFound, "Update user failed!", fmt.Sprintf("User with ID %d not found!", reqID.ID))
			ctx.JSON(http.StatusNotFound, response)
			return
		}

		response := helper.APIResponseError(http.StatusInternalServerError, "Update user failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	formatData := user.FormatUserFullData(updatedUser)
	response := helper.APIResponse(http.StatusOK, "Update user successfully!", formatData)

	handler.logsSvc.CreateActivityLog(ctx, fmt.Sprintf("%v updating user id %v.", reqUpdate.User.Name, reqID.ID))

	ctx.JSON(http.StatusOK, response)
}

func (handler *userHandler) DeleteUser(ctx *gin.Context) {
	var reqID user.RequestGetUserByID

	err := ctx.ShouldBindUri(&reqID)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Delete user failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	var reqDelete user.RequestDeleteUser

	err = ctx.ShouldBind(&reqDelete)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Delete user failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	reqDelete.User = ctx.MustGet("userData").(user.User)

	if _, err = handler.userSvc.DeleteUser(reqID, reqDelete); err != nil {
		if helper.IsErrNoRows(err.Error()) {
			response := helper.APIResponseError(http.StatusNotFound, "Delete user failed!", fmt.Sprintf("User with ID %d not found!", reqID.ID))
			ctx.JSON(http.StatusNotFound, response)
			return
		}

		response := helper.APIResponseError(http.StatusInternalServerError, "Delete user failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.BasicAPIResponse(http.StatusOK, "Delete user successfully!")

	handler.logsSvc.CreateActivityLog(ctx, fmt.Sprintf("%v deleting user id %v.", reqDelete.User.Name, reqID.ID))

	ctx.JSON(http.StatusOK, response)
}

func (handler *userHandler) AdminDataTablesUsers(ctx *gin.Context) {
	dataTablesUsers, err := handler.userSvc.AdminDataTablesUsers(ctx)

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Get datatables users failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	ctx.JSON(http.StatusOK, dataTablesUsers)
}

func (handler *userHandler) GetUserData(ctx *gin.Context) {
	formatData := user.FormatUserFullData(ctx.MustGet("userData").(user.User))
	response := helper.APIResponse(http.StatusOK, "Get user data successfully!", formatData)
	ctx.JSON(http.StatusOK, response)
}

func (handler *userHandler) ChangeUserData(ctx *gin.Context) {
	var reqUpdate user.RequestUpdateUser
	var reqSelfUpdate user.RequestSelfUpdateUser

	err := ctx.ShouldBind(&reqSelfUpdate)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Update self user data failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	reqUpdate.User = ctx.MustGet("userData").(user.User)
	reqUpdate.Name = reqSelfUpdate.Name
	reqUpdate.Email = reqSelfUpdate.Email
	reqUpdate.Role = reqUpdate.User.Role

	if reqSelfUpdate.Password != "" {
		reqUpdate.Password = reqSelfUpdate.Password
	}

	updatedUser, err := handler.userSvc.UpdateUser(user.RequestGetUserByID{ID: reqUpdate.User.ID}, reqUpdate)

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Update self user data failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	formatData := user.FormatUserFullData(updatedUser)
	response := helper.APIResponse(http.StatusOK, "Update self user data successfully!", formatData)

	handler.logsSvc.CreateActivityLog(ctx, fmt.Sprintf("%v self updating user data.", reqUpdate.User.Name))

	ctx.JSON(http.StatusOK, response)
}

func (handler *userHandler) GetNameByID(ctx *gin.Context) {
	var req user.RequestGetUserByID

	err := ctx.ShouldBindUri(&req)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Get user name by user id failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	userDetail, err := handler.userSvc.GetUserByID(req.ID)

	if err != nil {
		if helper.IsErrNoRows(err.Error()) {
			response := helper.APIResponseError(http.StatusNotFound, "Get user name by user id failed!", "Data not found!")
			ctx.JSON(http.StatusNotFound, response)
			return
		}

		response := helper.APIResponseError(http.StatusInternalServerError, "Get user name by user id failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.APIResponse(http.StatusOK, "Get user name by user id successfully!", gin.H{"name": userDetail.Name})

	ctx.JSON(http.StatusOK, response)
}

func (handler *userHandler) CreateForgotPasswordToken(ctx *gin.Context) {
	var req user.RequestCreateForgotPasswordToken

	err := ctx.ShouldBind(&req)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Request forgot password failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	userData, err := handler.userSvc.GetUserByEmail(req.Email)

	if err != nil {
		if helper.IsErrNoRows(err.Error()) {
			response := helper.APIResponseError(http.StatusNotFound, "Request forgot password failed!", fmt.Sprintf("User with email %s not found!", req.Email))
			ctx.JSON(http.StatusNotFound, response)
			return
		}

		response := helper.APIResponseError(http.StatusInternalServerError, "Request forgot password failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	req.User = userData

	dataCreated, err := handler.userSvc.CreateUserForgotPasswordToken(req)

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Request forgot password failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	// {
	// 	templateData := helper.EmailForgotPassword{
	// 		Name: userData.Name,
	// 		URL:  os.Getenv("WEB_URL") + "/auth/forgot-password/" + dataCreated.Token,
	// 	}
	// 	go helper.SendMail(userData.Email, "Forgot Password Request", templateData, "html/forgot_password.html")
	// }

	response := helper.APIResponse(http.StatusCreated, "Request forgot password successfully, please check your email inbox or spam!", dataCreated)

	handler.logsSvc.CreateActivityLog(ctx, fmt.Sprintf("%v make a request token for forgot password.", req.User.Name))

	ctx.JSON(http.StatusCreated, response)
}

func (handler *userHandler) ProcessForgotPasswordToken(ctx *gin.Context) {
	var req user.RequestProcessForgotPasswordToken

	err := ctx.ShouldBindUri(&req)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Process request forgot password failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	forgotPasswordData, err := handler.userSvc.GetDataForgotPasswordByToken(req.Token)

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Process request forgot password failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	expired := forgotPasswordData.CreatedAt.AddDate(0, 0, 1)
	isExpired := time.Now().After(expired)

	if isExpired {
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Process request forgot password failed!", "Token expired!")
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	userData, err := handler.userSvc.GetUserByID(forgotPasswordData.UserID)

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Process request forgot password failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	token, err := handler.authSvc.GenerateToken(userData.ID)

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Process request forgot password failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if _, err = handler.userSvc.DeleteForgotPasswordToken(forgotPasswordData); err != nil {
		if helper.IsErrNoRows(err.Error()) {
			response := helper.APIResponseError(http.StatusNotFound, "Process request forgot password failed!", fmt.Sprintf("Forgot password data with ID %d not found!", forgotPasswordData.ID))
			ctx.JSON(http.StatusNotFound, response)
			return
		}

		response := helper.APIResponseError(http.StatusInternalServerError, "Process request forgot password failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	formatData := user.FormatUserData(userData, token)
	response := helper.APIResponse(http.StatusOK, "Process request forgot password successfully!", formatData)

	handler.logsSvc.CreateActivityLog(ctx, fmt.Sprintf("%v process request forgot password.", req.User.Name))

	ctx.JSON(http.StatusOK, response)
}
