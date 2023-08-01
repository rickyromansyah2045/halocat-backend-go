package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rickyromansyah2045/halocat-backend-go/content"
	"github.com/rickyromansyah2045/halocat-backend-go/helper"
	"github.com/rickyromansyah2045/halocat-backend-go/logs"
	"github.com/rickyromansyah2045/halocat-backend-go/user"
)

type webAndCMSHandler struct {
	contentSvc content.Service
	userSvc    user.Service
	logsSvc    logs.Service
}

func NewWebAndCMSHandler(
	contentService content.Service,
	userService user.Service,
	logsService logs.Service,
) *webAndCMSHandler {
	return &webAndCMSHandler{
		contentSvc: contentService,
		userSvc:    userService,
		logsSvc:    logsService,
	}
}

func (handler *webAndCMSHandler) GetStatisticsForHomePage(ctx *gin.Context) {
	totalDonation, err := handler.contentSvc.GetTotalContent()

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Get statistics for home page failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	donationCompleted, err := handler.contentSvc.GetContentCompleted()

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Get statistics for home page failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	userRegistered, err := handler.userSvc.GetUserRegistered("")

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Get statistics for home page failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.APIResponse(http.StatusOK, "Get statistics for home page successfully!", gin.H{
		"total_donation":     totalDonation,
		"donation_completed": donationCompleted,
		"user_registered":    userRegistered,
	})

	ctx.JSON(http.StatusOK, response)
}

func (handler *webAndCMSHandler) GetStatisticsForAdminDashboard(ctx *gin.Context) {
	totalDonation, err := handler.contentSvc.GetTotalContent()

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Get statistics for admin dashboard failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	donationCompleted, err := handler.contentSvc.GetContentCompleted()

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Get statistics for admin dashboard failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	userRegistered, err := handler.userSvc.GetUserRegistered("")

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Get statistics for admin dashboard failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	userRegisteredRoleAdmin, err := handler.userSvc.GetUserRegistered("AND role = 'admin'")

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Get statistics for admin dashboard failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.APIResponse(http.StatusOK, "Get statistics for admin dashboard successfully!", gin.H{
		"total_donation":        totalDonation,
		"donation_completed":    donationCompleted,
		"user_registered":       userRegistered,
		"user_admin_registered": userRegisteredRoleAdmin,
	})

	ctx.JSON(http.StatusOK, response)
}
