package content

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/rickyromansyah2045/halocat-backend-go/constant"
	"github.com/rickyromansyah2045/halocat-backend-go/helper"
	"github.com/rickyromansyah2045/halocat-backend-go/user"
)

func (svc *service) GetAllContent(ctx *gin.Context) ([]Content, error) {
	contents, err := svc.repo.GetAllContent(ctx)

	if err != nil {
		return contents, err
	}

	return contents, nil
}

func (svc *service) GetContentByID(req RequestGetContentByID) (Content, error) {
	content, err := svc.repo.GetContentByID(req.ID)

	if err != nil {
		return content, err
	}

	return content, nil
}

func (svc *service) CreateContent(req RequestCreateContent) (Content, error) {
	content := Content{}
	content.UserID = req.User.ID
	content.CategoryID = req.CategoryID
	content.Title = req.Title
	content.ShortDescription = req.ShortDescription
	content.Description = req.Description
	content.Status = req.Status

	layoutFormat := "2006-01-02"
	finishedAt, _ := time.Parse(layoutFormat, req.FinishedAt)

	content.FinishedAt = finishedAt
	slugCandidate := fmt.Sprintf("%s %v%d", req.Title, time.Now().Unix(), req.User.ID)
	content.Slug = slug.Make(slugCandidate)
	content.CreatedBy = helper.SetNS(strconv.Itoa(req.User.ID))

	newContentData, err := svc.repo.SaveContent(content)

	if err != nil {
		return newContentData, err
	}

	return newContentData, nil
}

func (svc *service) UpdateContent(reqDetail RequestGetContentByID, reqUpdate RequestUpdateContent) (Content, error) {
	content, err := svc.repo.GetContentByID(reqDetail.ID)

	if err != nil {
		return content, err
	}

	if content.UserID != reqUpdate.User.ID && reqUpdate.User.Role == "user" {
		return content, errors.New("not an owner of the content")
	}

	if reqUpdate.User.Role == "user" {
		content.UserID = reqUpdate.User.ID
	} else {
		content.UserID = reqUpdate.UserID
	}

	content.CategoryID = reqUpdate.CategoryID
	content.Title = reqUpdate.Title
	content.ShortDescription = reqUpdate.ShortDescription
	content.Description = reqUpdate.Description
	content.Status = reqUpdate.Status

	if reqUpdate.FinishedAt != "" {
		layoutFormat := "2006-01-02"
		finishedAt, _ := time.Parse(layoutFormat, reqUpdate.FinishedAt)

		content.FinishedAt = finishedAt
	}

	content.UpdatedBy = helper.SetNS(strconv.Itoa(reqUpdate.User.ID))

	updatedContent, err := svc.repo.UpdateContent(content)

	if err != nil {
		return updatedContent, err
	}

	return updatedContent, nil
}

func (svc *service) DeleteContent(reqDetail RequestGetContentByID, reqDelete RequestDeleteContent) (bool, error) {
	if constant.DELETED_BY {
		content, err := svc.repo.GetContentByID(reqDetail.ID)

		if err != nil {
			return false, err
		}

		content.UpdatedBy = helper.SetNS(strconv.Itoa(reqDelete.User.ID))
		content.DeletedAt = *helper.SetNowNT()
		content.DeletedBy = helper.SetNS(strconv.Itoa(reqDelete.User.ID))

		status, err := svc.repo.DeleteContent(content)

		if err != nil {
			return status, err
		}

		return status, nil
	}

	content, err := svc.repo.GetContentByID(reqDetail.ID)

	if err != nil {
		return false, err
	}

	if content.UserID != reqDelete.User.ID && reqDelete.User.Role == "user" {
		return false, errors.New("not an owner of the content")
	}

	content.ID = reqDetail.ID
	status, err := svc.repo.DeleteContent(content)

	if err != nil {
		return status, err
	}

	return status, nil
}

func (svc *service) SaveContentImage(req RequestCreateContentImage, fileLocation string) (contentImage ContentImage, err error) {
	content, err := svc.repo.GetContentByID(req.ContentID)

	if err != nil {
		return contentImage, err
	}

	if content.UserID != req.User.ID && req.User.Role == "user" {
		return contentImage, errors.New("not an owner of the content")
	}

	isPrimary := 0
	if req.IsPrimary {
		if _, err := svc.repo.UpdateAllImagesAsNonPrimary(req.ContentID); err != nil {
			return contentImage, err
		}
		isPrimary = 1
	}

	contentImage.ContentID = req.ContentID
	contentImage.IsPrimary = isPrimary
	contentImage.FileLocation = fileLocation
	contentImage.CreatedBy = helper.SetNS(strconv.Itoa(req.User.ID))

	newContentImage, err := svc.repo.CreateContentImage(contentImage)

	if err != nil {
		return newContentImage, err
	}

	return newContentImage, nil
}

func (svc *service) GetAllContentImage() ([]ContentImage, error) {
	contentImages, err := svc.repo.GetAllContentImage()

	if err != nil {
		return contentImages, err
	}

	return contentImages, nil
}

func (svc *service) GetContentImageByID(req RequestGetContentImageByID) (ContentImage, error) {
	contentImage, err := svc.repo.GetContentImageByID(req.ID)

	if err != nil {
		return contentImage, err
	}

	return contentImage, nil
}

func (svc *service) DeleteContentImage(reqDetail RequestGetContentImageByID, reqDelete RequestDeleteContentImage) (bool, error) {
	getContentImage, err := svc.repo.GetContentImageByID(reqDetail.ID)

	if err != nil {
		return false, err
	}

	content, err := svc.repo.GetContentByID(getContentImage.ContentID)

	if err != nil {
		return false, err
	}

	if content.UserID != reqDelete.User.ID && reqDelete.User.Role == "user" {
		return false, errors.New("not an owner of the content image")
	}

	contentImage := ContentImage{}
	contentImage.ID = reqDetail.ID
	status, err := svc.repo.DeleteContentImage(contentImage)

	if err != nil {
		return status, err
	}

	return status, nil
}

func (svc *service) GetAllContentCategory() ([]ContentCategory, error) {
	contentCategory, err := svc.repo.GetAllContentCategory()

	if err != nil {
		return contentCategory, err
	}

	return contentCategory, nil
}

func (svc *service) GetContentCategoryByID(req RequestGetContentCategoryByID) (ContentCategory, error) {
	contentCategory, err := svc.repo.GetContentCategoryByID(req.ID)

	if err != nil {
		return contentCategory, err
	}

	return contentCategory, nil
}

func (svc *service) DeleteContentCategory(reqDetail RequestGetContentCategoryByID, reqDelete RequestDeleteContentCategory) (bool, error) {
	if constant.DELETED_BY {
		contentCategory, err := svc.repo.GetContentCategoryByID(reqDetail.ID)

		if err != nil {
			return false, err
		}

		contentCategory.UpdatedBy = helper.SetNS(strconv.Itoa(reqDelete.User.ID))
		contentCategory.DeletedAt = *helper.SetNowNT()
		contentCategory.DeletedBy = helper.SetNS(strconv.Itoa(reqDelete.User.ID))

		status, err := svc.repo.DeleteContentCategory(contentCategory)

		if err != nil {
			return status, err
		}

		return status, nil
	}

	contentCategory := ContentCategory{}
	contentCategory.ID = reqDetail.ID
	status, err := svc.repo.DeleteContentCategory(contentCategory)

	if err != nil {
		return status, err
	}

	return status, nil
}

func (svc *service) CreateContentCategory(req RequestCreateContentCategory) (ContentCategory, error) {
	category := ContentCategory{}
	category.Category = req.Category
	category.CreatedBy = helper.SetNS(strconv.Itoa(req.User.ID))

	newContentCategoryData, err := svc.repo.SaveContentCategory(category)

	if err != nil {
		return newContentCategoryData, err
	}

	return newContentCategoryData, nil
}

func (svc *service) UpdateContentCategory(reqDetail RequestGetContentCategoryByID, reqUpdate RequestUpdateContentCategory) (contentCategory ContentCategory, err error) {
	contentCategory, err = svc.repo.GetContentCategoryByID(reqDetail.ID)

	if err != nil {
		return contentCategory, err
	}

	contentCategory.ID = reqDetail.ID
	contentCategory.Category = reqUpdate.Category
	contentCategory.UpdatedBy = helper.SetNS(strconv.Itoa(reqUpdate.User.ID))

	updatedContentCategory, err := svc.repo.UpdateContentCategory(contentCategory)

	if err != nil {
		return updatedContentCategory, err
	}

	return updatedContentCategory, nil
}

func (svc *service) AdminDataTablesContents(ctx *gin.Context) (helper.DataTables, error) {
	dataTablesContents, err := svc.repo.AdminDataTablesContents(ctx)

	if err != nil {
		return dataTablesContents, err
	}

	return dataTablesContents, nil
}

func (svc *service) AdminDataTablesCategories(ctx *gin.Context) (helper.DataTables, error) {
	dataTablesCategories, err := svc.repo.AdminDataTablesCategories(ctx)

	if err != nil {
		return dataTablesCategories, err
	}

	return dataTablesCategories, nil
}

func (svc *service) UserDataTablesContents(ctx *gin.Context, user user.User) (helper.DataTables, error) {
	dataTablesContents, err := svc.repo.UserDataTablesContents(ctx, user)

	if err != nil {
		return dataTablesContents, err
	}

	return dataTablesContents, nil
}

func (svc *service) GetTotalContent() (res int, err error) {
	res, err = svc.repo.GetTotalContent()

	if err != nil {
		return res, err
	}

	return res, nil
}

func (svc *service) GetContentCompleted() (res int, err error) {
	res, err = svc.repo.GetContentCompleted()

	if err != nil {
		return res, err
	}

	return res, nil
}
