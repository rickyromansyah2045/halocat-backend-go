package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/rickyromansyah2045/halocat-backend-go/content"
	"github.com/rickyromansyah2045/halocat-backend-go/helper"
	"github.com/rickyromansyah2045/halocat-backend-go/logs"
	"github.com/rickyromansyah2045/halocat-backend-go/user"
)

type contentHandler struct {
	contentSvc content.Service
	userSvc    user.Service
	logsSvc    logs.Service
}

func NewContentHandler(
	contentService content.Service,
	userService user.Service,
	logsService logs.Service,
) *contentHandler {
	return &contentHandler{
		contentSvc: contentService,
		userSvc:    userService,
		logsSvc:    logsService,
	}
}

func (handler *contentHandler) GetAllContent(ctx *gin.Context) {
	contents, err := handler.contentSvc.GetAllContent(ctx)

	if err != nil {
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Get contents failed!", err.Error())
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	formatData := content.FormatMultipleContentData(contents)
	response := helper.APIResponse(http.StatusOK, "Get contents successfully!", formatData)

	ctx.JSON(http.StatusOK, response)
}

func (handler *contentHandler) GetContentByID(ctx *gin.Context) {
	var req content.RequestGetContentByID

	err := ctx.ShouldBindUri(&req)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Get detail content failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	contentDetail, err := handler.contentSvc.GetContentByID(req)

	if err != nil {
		if helper.IsErrNoRows(err.Error()) {
			response := helper.APIResponseError(http.StatusNotFound, "Get detail content failed!", "Data not found!")
			ctx.JSON(http.StatusNotFound, response)
			return
		}

		response := helper.APIResponseError(http.StatusInternalServerError, "Get detail content failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	formatData := content.FormatContentData(contentDetail)
	response := helper.APIResponse(http.StatusOK, "Get detail content successfully!", formatData)

	ctx.JSON(http.StatusOK, response)
}

func (handler *contentHandler) CreateContent(ctx *gin.Context) {
	var req content.RequestCreateContent

	err := ctx.ShouldBind(&req)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Create content failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	userData := ctx.MustGet("userData").(user.User)

	if req.UserID != 0 {
		if userData.Role == "user" {
			response := helper.APIResponseError(http.StatusBadRequest, "Create content failed!", "Bad Request!")
			ctx.JSON(http.StatusBadRequest, response)
			return
		}

		if _, err := handler.userSvc.GetUserByID(req.UserID); err != nil {
			if helper.IsErrNoRows(err.Error()) {
				response := helper.APIResponseError(http.StatusNotFound, "Create content failed!", fmt.Sprintf("User with ID %d not found!", req.UserID))
				ctx.JSON(http.StatusNotFound, response)
				return
			}

			response := helper.APIResponseError(http.StatusInternalServerError, "Create content failed!", err.Error())
			ctx.JSON(http.StatusInternalServerError, response)
			return
		}

		req.User = user.User{ID: req.UserID}
	} else {
		req.User = userData
	}

	newContentData, err := handler.contentSvc.CreateContent(req)

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Create content failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	formatData := content.FormatContentData(newContentData)
	response := helper.APIResponse(http.StatusCreated, "Create content successfully!", formatData)

	handler.logsSvc.CreateActivityLog(ctx, fmt.Sprintf("%v creating content id %v.", userData.Name, newContentData.ID))

	ctx.JSON(http.StatusCreated, response)
}

func (handler *contentHandler) UpdateContent(ctx *gin.Context) {
	var reqID content.RequestGetContentByID

	err := ctx.ShouldBindUri(&reqID)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Update content failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	var reqUpdate content.RequestUpdateContent

	err = ctx.ShouldBind(&reqUpdate)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Update content failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	reqUpdate.User = ctx.MustGet("userData").(user.User)

	oldContent, err := handler.contentSvc.GetContentByID(reqID)

	if err != nil {
		if helper.IsErrNoRows(err.Error()) {
			response := helper.APIResponseError(http.StatusNotFound, "Update content failed!", fmt.Sprintf("Content with ID %d not found!", reqID.ID))
			ctx.JSON(http.StatusNotFound, response)
			return
		}

		response := helper.APIResponseError(http.StatusInternalServerError, "Update content failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	updatedContent, err := handler.contentSvc.UpdateContent(reqID, reqUpdate)

	if err != nil {
		if helper.IsErrNoRows(err.Error()) {
			response := helper.APIResponseError(http.StatusNotFound, "Update content failed!", fmt.Sprintf("Content with ID %d not found!", reqID.ID))
			ctx.JSON(http.StatusNotFound, response)
			return
		}

		response := helper.APIResponseError(http.StatusInternalServerError, "Update content failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if oldContent.Status != "active" && updatedContent.Status == "active" {
		// ownerContentUserData, err := handler.userSvc.GetUserByID(updatedContent.UserID)

		if err != nil {
			response := helper.APIResponseError(http.StatusInternalServerError, "Update content failed!", err.Error())
			ctx.JSON(http.StatusInternalServerError, response)
			return
		}

		// {
		// 	templateData := helper.EmailContentActive{
		// 		Name:        ownerContentUserData.Name,
		// 		Content:     updatedContent,
		// 		ContentLink: os.Getenv("WEB_URL") + "/donate/" + strconv.Itoa(updatedContent.ID),
		// 	}
		// 	go helper.SendMail(ownerContentUserData.Email, "Your Content Now Active!", templateData, "html/content_active.html")
		// }
	}

	formatData := content.FormatContentData(updatedContent)
	response := helper.APIResponse(http.StatusOK, "Update content successfully!", formatData)

	handler.logsSvc.CreateActivityLog(ctx, fmt.Sprintf("%v updating content id %v.", reqUpdate.User.Name, reqID.ID))

	ctx.JSON(http.StatusOK, response)
}

func (handler *contentHandler) DeleteContent(ctx *gin.Context) {
	var reqID content.RequestGetContentByID

	err := ctx.ShouldBindUri(&reqID)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Delete content failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	var reqDelete content.RequestDeleteContent

	err = ctx.ShouldBind(&reqDelete)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Delete content failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	reqDelete.User = ctx.MustGet("userData").(user.User)

	if _, err = handler.contentSvc.DeleteContent(reqID, reqDelete); err != nil {
		if helper.IsErrNoRows(err.Error()) {
			response := helper.APIResponseError(http.StatusNotFound, "Delete content failed!", fmt.Sprintf("Content with ID %d not found!", reqID.ID))
			ctx.JSON(http.StatusNotFound, response)
			return
		}

		response := helper.APIResponseError(http.StatusInternalServerError, "Delete content failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.BasicAPIResponse(http.StatusOK, "Delete content successfully!")

	handler.logsSvc.CreateActivityLog(ctx, fmt.Sprintf("%v deleting content id %v.", reqDelete.User.Name, reqID.ID))

	ctx.JSON(http.StatusOK, response)
}

func (handler *contentHandler) GetAllContentImage(ctx *gin.Context) {
	contentImages, err := handler.contentSvc.GetAllContentImage()

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Get content images failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	formatData := content.FormatMultipleContentImageData(contentImages)
	response := helper.APIResponse(http.StatusOK, "Get content images successfully!", formatData)

	ctx.JSON(http.StatusOK, response)
}

func (handler *contentHandler) GetContentImageByID(ctx *gin.Context) {
	var req content.RequestGetContentImageByID

	err := ctx.ShouldBindUri(&req)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Get detail content image failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	contentImageDetail, err := handler.contentSvc.GetContentImageByID(req)

	if err != nil {
		if helper.IsErrNoRows(err.Error()) {
			response := helper.APIResponseError(http.StatusNotFound, "Get detail content image failed!", "Data not found!")
			ctx.JSON(http.StatusNotFound, response)
			return
		}

		response := helper.APIResponseError(http.StatusInternalServerError, "Get detail content image failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	formatData := content.FormatContentImageData(contentImageDetail)
	response := helper.APIResponse(http.StatusOK, "Get detail content image successfully!", formatData)

	ctx.JSON(http.StatusOK, response)
}

func (handler *contentHandler) UploadImage(ctx *gin.Context) {
	var req content.RequestCreateContentImage

	err := ctx.ShouldBind(&req)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Upload content image failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	req.User = ctx.MustGet("userData").(user.User)

	file, err := ctx.FormFile("file")
	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Upload content image failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	slug := slug.Make(fmt.Sprintf("%d %v %s", 1, time.Now().Unix(), file.Filename[:len(file.Filename)-len(filepath.Ext(file.Filename))]))
	path := fmt.Sprintf("images/%s%v", slug, filepath.Ext(file.Filename))

	if err := ctx.SaveUploadedFile(file, path); err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Upload content image failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	uploadedContentImage, err := handler.contentSvc.SaveContentImage(req, path)

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Upload content image failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	formatData := content.FormatContentImageData(uploadedContentImage)
	response := helper.APIResponse(http.StatusOK, "Upload content image successfully!", formatData)

	handler.logsSvc.CreateActivityLog(ctx, fmt.Sprintf("%v uploading image id %v for content id %v.", req.User.Name, uploadedContentImage.ID, uploadedContentImage.ContentID))

	ctx.JSON(http.StatusOK, response)
}

func (handler *contentHandler) DeleteContentImage(ctx *gin.Context) {
	var reqID content.RequestGetContentImageByID

	err := ctx.ShouldBindUri(&reqID)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Delete content image failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	var reqDelete content.RequestDeleteContentImage

	err = ctx.ShouldBind(&reqDelete)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Delete content image failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	reqDelete.User = ctx.MustGet("userData").(user.User)

	if _, err = handler.contentSvc.DeleteContentImage(reqID, reqDelete); err != nil {
		if helper.IsErrNoRows(err.Error()) {
			response := helper.APIResponseError(http.StatusNotFound, "Delete content image failed!", fmt.Sprintf("Content image with ID %d not found!", reqID.ID))
			ctx.JSON(http.StatusNotFound, response)
			return
		}

		response := helper.APIResponseError(http.StatusInternalServerError, "Delete content image failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.BasicAPIResponse(http.StatusOK, "Delete content image successfully!")

	handler.logsSvc.CreateActivityLog(ctx, fmt.Sprintf("%v deleting content image id %v.", reqDelete.User.Name, reqID.ID))

	ctx.JSON(http.StatusOK, response)
}

func (handler *contentHandler) GetAllContentCategory(ctx *gin.Context) {
	contentCategories, err := handler.contentSvc.GetAllContentCategory()

	if err != nil {
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Get content categories failed!", err.Error())
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	formatData := content.FormatMultipleContentCategoryData(contentCategories)
	response := helper.APIResponse(http.StatusOK, "Get content categories successfully!", formatData)

	ctx.JSON(http.StatusOK, response)
}

func (handler *contentHandler) GetContentCategoryByID(ctx *gin.Context) {
	var req content.RequestGetContentCategoryByID

	err := ctx.ShouldBindUri(&req)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Get detail content category failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	contentCategoryDetail, err := handler.contentSvc.GetContentCategoryByID(req)

	if err != nil {
		if helper.IsErrNoRows(err.Error()) {
			response := helper.APIResponseError(http.StatusNotFound, "Get detail content category failed!", "Data not found!")
			ctx.JSON(http.StatusNotFound, response)
			return
		}

		response := helper.APIResponseError(http.StatusInternalServerError, "Get detail content category failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	formatData := content.FormatContentCategoryData(contentCategoryDetail)
	response := helper.APIResponse(http.StatusOK, "Get detail content category successfully!", formatData)

	ctx.JSON(http.StatusOK, response)
}

func (handler *contentHandler) DeleteContentCategory(ctx *gin.Context) {
	var reqID content.RequestGetContentCategoryByID

	err := ctx.ShouldBindUri(&reqID)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Delete content category failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	var reqDelete content.RequestDeleteContentCategory

	err = ctx.ShouldBind(&reqDelete)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Delete content category failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	reqDelete.User = ctx.MustGet("userData").(user.User)

	if _, err = handler.contentSvc.DeleteContentCategory(reqID, reqDelete); err != nil {
		if helper.IsErrNoRows(err.Error()) {
			response := helper.APIResponseError(http.StatusNotFound, "Delete content category failed!", fmt.Sprintf("Content category with ID %d not found!", reqID.ID))
			ctx.JSON(http.StatusNotFound, response)
			return
		}

		response := helper.APIResponseError(http.StatusInternalServerError, "Delete content category failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.BasicAPIResponse(http.StatusOK, "Delete content category successfully!")

	handler.logsSvc.CreateActivityLog(ctx, fmt.Sprintf("%v deleting category id %v.", reqDelete.User.Name, reqID.ID))

	ctx.JSON(http.StatusOK, response)
}

func (handler *contentHandler) CreateContentCategory(ctx *gin.Context) {
	var req content.RequestCreateContentCategory

	err := ctx.ShouldBind(&req)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Create content category failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	req.User = ctx.MustGet("userData").(user.User)

	newContentCategoryData, err := handler.contentSvc.CreateContentCategory(req)

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Create content category failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	formatData := content.FormatContentCategoryData(newContentCategoryData)
	response := helper.APIResponse(http.StatusCreated, "Create content category successfully!", formatData)

	handler.logsSvc.CreateActivityLog(ctx, fmt.Sprintf("%v creating category id %v.", req.User.Name, newContentCategoryData.ID))

	ctx.JSON(http.StatusCreated, response)
}

func (handler *contentHandler) UpdateContentCategory(ctx *gin.Context) {
	var reqID content.RequestGetContentCategoryByID

	err := ctx.ShouldBindUri(&reqID)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Update content category failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	var reqUpdate content.RequestUpdateContentCategory

	err = ctx.ShouldBind(&reqUpdate)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Update content category failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	reqUpdate.User = ctx.MustGet("userData").(user.User)

	updatedContentCategory, err := handler.contentSvc.UpdateContentCategory(reqID, reqUpdate)

	if err != nil {
		if helper.IsErrNoRows(err.Error()) {
			response := helper.APIResponseError(http.StatusNotFound, "Update content category failed!", fmt.Sprintf("Content category with ID %d not found!", reqID.ID))
			ctx.JSON(http.StatusNotFound, response)
			return
		}

		response := helper.APIResponseError(http.StatusInternalServerError, "Update content category failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	formatData := content.FormatContentCategoryData(updatedContentCategory)
	response := helper.APIResponse(http.StatusOK, "Update content category successfully!", formatData)

	handler.logsSvc.CreateActivityLog(ctx, fmt.Sprintf("%v updating category id %v.", reqUpdate.User.Name, reqID.ID))

	ctx.JSON(http.StatusOK, response)
}

func (handler *contentHandler) AdminDataTablesContents(ctx *gin.Context) {
	dataTablesContents, err := handler.contentSvc.AdminDataTablesContents(ctx)

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Get datatables contents failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	ctx.JSON(http.StatusOK, dataTablesContents)
}

func (handler *contentHandler) AdminDataTablesCategories(ctx *gin.Context) {
	dataTablesCategories, err := handler.contentSvc.AdminDataTablesCategories(ctx)

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Get datatables categories failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	ctx.JSON(http.StatusOK, dataTablesCategories)
}

func (handler *contentHandler) UserDataTablesContents(ctx *gin.Context) {
	userData := ctx.MustGet("userData").(user.User)
	dataTablesContents, err := handler.contentSvc.UserDataTablesContents(ctx, userData)

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Get datatables contents failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	ctx.JSON(http.StatusOK, dataTablesContents)
}
