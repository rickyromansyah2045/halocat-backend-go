package content

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rickyromansyah2045/halocat-backend-go/constant"
	"github.com/rickyromansyah2045/halocat-backend-go/helper"
	"github.com/rickyromansyah2045/halocat-backend-go/user"
)

func (repo *repository) GetAllContent(ctx *gin.Context) (contents []Content, err error) {
	additionalQuery := ""

	if ctx.Query("search") != "" {
		s := ctx.Query("search")
		additionalQuery += fmt.Sprintf(" AND (title LIKE '%%%v%%' OR short_description LIKE '%%%v%%')", s, s)
	}

	if ctx.Query("status") != "" {
		additionalQuery += fmt.Sprintf(" AND status = '%v'", ctx.Query("status"))
	}

	if ctx.Query("category") != "" {
		additionalQuery += fmt.Sprintf(" AND category_id = %v", ctx.Query("category"))
	}

	if ctx.Query("limit") != "" {
		additionalQuery += fmt.Sprintf(" LIMIT %v", ctx.Query("limit"))

		if ctx.Query("offset") != "" {
			additionalQuery += fmt.Sprintf(" OFFSET %v", ctx.Query("offset"))
		}
	}

	rows, err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryGetAll + additionalQuery)).Rows()

	if err != nil {
		return contents, err
	}

	defer rows.Close()

	for rows.Next() {
		tmp := Content{}
		err := rows.Scan(
			&tmp.ID,
			&tmp.UserID,
			&tmp.CategoryID,
			&tmp.Title,
			&tmp.Slug,
			&tmp.ShortDescription,
			&tmp.Description,
			&tmp.Status,
			&tmp.FinishedAt,
			&tmp.CreatedAt,
			&tmp.CreatedBy,
			&tmp.UpdatedAt,
			&tmp.UpdatedBy,
			&tmp.DeletedAt,
			&tmp.DeletedBy,
		)

		if err != nil {
			return contents, err
		}

		imageRows, err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryGetContentImages), tmp.ID).Rows()

		if err != nil {
			return contents, err
		}

		defer imageRows.Close()

		contentImages := []ContentImage{}

		for imageRows.Next() {
			tmpContentImages := ContentImage{}
			errContentImages := imageRows.Scan(
				&tmpContentImages.ID,
				&tmpContentImages.ContentID,
				&tmpContentImages.FileLocation,
				&tmpContentImages.IsPrimary,
				&tmpContentImages.CreatedAt,
				&tmpContentImages.CreatedBy,
				&tmpContentImages.UpdatedAt,
				&tmpContentImages.UpdatedBy,
				&tmpContentImages.DeletedAt,
				&tmpContentImages.DeletedBy,
			)

			if errContentImages != nil {
				return contents, err
			}

			contentImages = append(contentImages, tmpContentImages)
		}

		tmp.ContentImages = contentImages
		contents = append(contents, tmp)
	}

	return contents, nil
}

func (repo *repository) GetContentByID(id int) (content Content, err error) {
	row := repo.DB.Raw(helper.ConvertToInLineQuery(QueryGetContentByID), id).Row()

	err = row.Scan(
		&content.ID,
		&content.UserID,
		&content.CategoryID,
		&content.Title,
		&content.Slug,
		&content.ShortDescription,
		&content.Description,
		&content.Status,
		&content.FinishedAt,
		&content.CreatedAt,
		&content.CreatedBy,
		&content.UpdatedAt,
		&content.UpdatedBy,
		&content.DeletedAt,
		&content.DeletedBy,
	)

	if err != nil {
		return content, err
	}

	imageRows, err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryGetContentImages), content.ID).Rows()

	if err != nil {
		return content, err
	}

	defer imageRows.Close()

	contentImages := []ContentImage{}

	for imageRows.Next() {
		tmpContentImages := ContentImage{}
		errContentImages := imageRows.Scan(
			&tmpContentImages.ID,
			&tmpContentImages.ContentID,
			&tmpContentImages.FileLocation,
			&tmpContentImages.IsPrimary,
			&tmpContentImages.CreatedAt,
			&tmpContentImages.CreatedBy,
			&tmpContentImages.UpdatedAt,
			&tmpContentImages.UpdatedBy,
			&tmpContentImages.DeletedAt,
			&tmpContentImages.DeletedBy,
		)

		if errContentImages != nil {
			return content, err
		}

		contentImages = append(contentImages, tmpContentImages)
	}

	content.ContentImages = contentImages

	return content, nil
}

func (repo *repository) SaveContent(content Content) (Content, error) {
	if err := repo.DB.Create(&content).Error; err != nil {
		return content, err
	}
	return content, nil
}

func (repo *repository) UpdateContent(content Content) (Content, error) {
	if err := repo.DB.Save(&content).Error; err != nil {
		return content, err
	}
	return content, nil
}

func (repo *repository) DeleteContent(content Content) (bool, error) {
	if constant.DELETED_BY {
		if err := repo.DB.Save(&content).Error; err != nil {
			return false, err
		}
		return true, nil
	}

	if err := repo.DB.Delete(&content).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (repo *repository) CreateContentImage(contentImage ContentImage) (ContentImage, error) {
	err := repo.DB.Create(&contentImage).Error
	if err != nil {
		return contentImage, err
	}
	return contentImage, nil
}

func (repo *repository) UpdateAllImagesAsNonPrimary(contentID int) (bool, error) {
	err := repo.DB.Model(&ContentImage{}).Where("content_id = ?", contentID).Update("is_primary", false).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func (repo *repository) GetAllContentImage() (contentImages []ContentImage, err error) {
	rows, err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryGetAllImage)).Rows()

	if err != nil {
		return contentImages, err
	}

	defer rows.Close()

	for rows.Next() {
		tmp := ContentImage{}
		err := rows.Scan(
			&tmp.ID,
			&tmp.ContentID,
			&tmp.FileLocation,
			&tmp.IsPrimary,
			&tmp.CreatedAt,
			&tmp.CreatedBy,
			&tmp.UpdatedAt,
			&tmp.UpdatedBy,
			&tmp.DeletedAt,
			&tmp.DeletedBy,
		)

		if err != nil {
			return contentImages, err
		}

		contentImages = append(contentImages, tmp)
	}

	return contentImages, nil
}

func (repo *repository) GetContentImageByID(id int) (contentImage ContentImage, err error) {
	row := repo.DB.Raw(helper.ConvertToInLineQuery(QueryGetContentImageByID), id).Row()

	err = row.Scan(
		&contentImage.ID,
		&contentImage.ContentID,
		&contentImage.FileLocation,
		&contentImage.IsPrimary,
		&contentImage.CreatedAt,
		&contentImage.CreatedBy,
		&contentImage.UpdatedAt,
		&contentImage.UpdatedBy,
		&contentImage.DeletedAt,
		&contentImage.DeletedBy,
	)

	if err != nil {
		return contentImage, err
	}

	return contentImage, nil
}

func (repo *repository) DeleteContentImage(contentImage ContentImage) (bool, error) {
	if err := repo.DB.Delete(&contentImage).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (repo *repository) GetAllContentCategory() (contentCategories []ContentCategory, err error) {
	rows, err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryGetAllCategory)).Rows()

	if err != nil {
		return contentCategories, err
	}

	defer rows.Close()

	for rows.Next() {
		tmp := ContentCategory{}
		err := rows.Scan(
			&tmp.ID,
			&tmp.Category,
			&tmp.CreatedAt,
			&tmp.CreatedBy,
			&tmp.UpdatedAt,
			&tmp.UpdatedBy,
			&tmp.DeletedAt,
			&tmp.DeletedBy,
		)

		if err != nil {
			return contentCategories, err
		}

		contentCategories = append(contentCategories, tmp)
	}

	return contentCategories, nil
}

func (repo *repository) GetContentCategoryByID(id int) (contentCategory ContentCategory, err error) {
	row := repo.DB.Raw(helper.ConvertToInLineQuery(QueryGetContentCategoryByID), id).Row()

	err = row.Scan(
		&contentCategory.ID,
		&contentCategory.Category,
		&contentCategory.CreatedAt,
		&contentCategory.CreatedBy,
		&contentCategory.UpdatedAt,
		&contentCategory.UpdatedBy,
		&contentCategory.DeletedAt,
		&contentCategory.DeletedBy,
	)

	if err != nil {
		return contentCategory, err
	}

	return contentCategory, nil
}

func (repo *repository) DeleteContentCategory(contentCategory ContentCategory) (bool, error) {
	tmpContentCategory := ContentCategory{}

	if err := repo.DB.Where("id = ?", contentCategory.ID).Find(&tmpContentCategory).Error; err != nil {
		return false, err
	}

	if tmpContentCategory.ID == 0 {
		return false, errors.New("sql: no rows in result set")
	}

	if constant.DELETED_BY {
		if err := repo.DB.Save(&contentCategory).Error; err != nil {
			return false, err
		}
		return true, nil
	}

	if err := repo.DB.Delete(&contentCategory).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (repo *repository) SaveContentCategory(category ContentCategory) (ContentCategory, error) {
	if err := repo.DB.Create(&category).Error; err != nil {
		return category, err
	}
	return category, nil
}

func (repo *repository) UpdateContentCategory(category ContentCategory) (ContentCategory, error) {
	tmpContentCategory := ContentCategory{}

	if err := repo.DB.Where("id = ?", category.ID).Find(&tmpContentCategory).Error; err != nil {
		return category, err
	}

	if tmpContentCategory.ID == 0 {
		return category, errors.New("sql: no rows in result set")
	}

	if err := repo.DB.Save(&category).Error; err != nil {
		return category, err
	}
	return category, nil
}

func (repo *repository) AdminDataTablesContents(ctx *gin.Context) (result helper.DataTables, err error) {
	var (
		query string = QueryAdminDataTablesContents
		likes string = ""
		order string = ""
		limit string = ""
	)

	var (
		no       int = 1
		total    int = 0
		filtered int = 0
	)

	var data []map[string]any

	listOrder := []string{"", "title", "total_image", "status", ""}

	searchValue := ctx.Query("search[value]")
	orderColumn := ctx.Query("order[0][column]")
	starting, _ := strconv.Atoi(ctx.Query("start"))

	if searchValue != "" {
		likes = fmt.Sprintf("(title LIKE '%%%s%%' OR short_description LIKE '%%%s%%' OR status LIKE '%%%s%%')", searchValue, searchValue, searchValue)
	}

	if orderColumn != "" {
		orderType := ctx.Query("order[0][dir]")
		orderColumn, _ := strconv.Atoi(orderColumn)
		order = fmt.Sprintf("ORDER BY %s %s", listOrder[orderColumn], strings.ToUpper(orderType))
	} else {
		order = "ORDER BY id DESC"
	}

	if starting != -1 {
		length, _ := strconv.Atoi(ctx.Query("length"))
		limit = fmt.Sprintf("LIMIT %v OFFSET %v", length, starting)
		no = starting + 1
	}

	if likes != "" {
		query = fmt.Sprintf("%s AND %s", query, likes)

		if err := repo.DB.Raw(fmt.Sprintf("%s AND %s", helper.ConvertToInLineQuery(QueryCountAllAdminDataTablesContents), likes)).Scan(&filtered).Error; err != nil {
			return result, err
		}

		if err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryCountAllAdminDataTablesContents)).Scan(&total).Error; err != nil {
			return result, err
		}
	} else {
		if err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryCountAllAdminDataTablesContents)).Scan(&filtered).Error; err != nil {
			return result, err
		}

		total = filtered
	}

	if order != "" {
		query = fmt.Sprintf("%s %s", query, order)
	}

	query = fmt.Sprintf("%s %s", query, limit)

	rows, err := repo.DB.Raw(helper.ConvertToInLineQuery(query)).Rows()

	if err != nil {
		return result, err
	}

	defer rows.Close()

	for rows.Next() {
		var (
			totalImage int = 0
			finishedAt sql.NullTime
			status     string
		)

		tmp := Content{}
		err := rows.Scan(
			&tmp.ID,
			&tmp.UserID,
			&tmp.CategoryID,
			&tmp.Title,
			&tmp.Slug,
			&tmp.ShortDescription,
			&tmp.Description,
			&totalImage,
			&status,
			&finishedAt,
			&tmp.CreatedAt,
			&tmp.CreatedBy,
			&tmp.UpdatedAt,
			&tmp.UpdatedBy,
			&tmp.DeletedAt,
			&tmp.DeletedBy,
		)

		if err != nil {
			return result, err
		}

		data = append(data, map[string]any{
			"no":                no,
			"id":                tmp.ID,
			"user_id":           tmp.UserID,
			"category_id":       tmp.CategoryID,
			"title":             tmp.Title,
			"slug":              tmp.Slug,
			"short_description": tmp.ShortDescription,
			"description":       tmp.Description,
			"total_image":       totalImage,
			"status":            status,
			"finished_at":       helper.HNTime(finishedAt),
			"created_at":        helper.HNTime(tmp.CreatedAt),
			"created_by":        helper.HNString(tmp.CreatedBy),
			"updated_at":        helper.HNTime(tmp.UpdatedAt),
			"updated_by":        helper.HNString(tmp.UpdatedBy),
			"deleted_at":        helper.HNTimeGDeletedAt(tmp.DeletedAt),
			"deleted_by":        helper.HNString(tmp.DeletedBy),
		})

		no++
	}

	return helper.BuildDatatTables(data, filtered, total), nil
}

func (repo *repository) AdminDataTablesCategories(ctx *gin.Context) (result helper.DataTables, err error) {
	var (
		query string = QueryAdminDataTablesCategories
		likes string = ""
		order string = ""
		limit string = ""
	)

	var (
		no       int = 1
		total    int = 0
		filtered int = 0
	)

	var data []map[string]any

	listOrder := []string{"", "category", ""}

	searchValue := ctx.Query("search[value]")
	orderColumn := ctx.Query("order[0][column]")
	starting, _ := strconv.Atoi(ctx.Query("start"))

	if searchValue != "" {
		likes = fmt.Sprintf("(category LIKE '%%%s%%')", searchValue)
	}

	if orderColumn != "" {
		orderType := ctx.Query("order[0][dir]")
		orderColumn, _ := strconv.Atoi(orderColumn)
		order = fmt.Sprintf("ORDER BY %s %s", listOrder[orderColumn], strings.ToUpper(orderType))
	} else {
		order = "ORDER BY id DESC"
	}

	if starting != -1 {
		length, _ := strconv.Atoi(ctx.Query("length"))
		limit = fmt.Sprintf("LIMIT %v OFFSET %v", length, starting)
		no = starting + 1
	}

	if likes != "" {
		query = fmt.Sprintf("%s AND %s", query, likes)

		if err := repo.DB.Raw(fmt.Sprintf("%s AND %s", helper.ConvertToInLineQuery(QueryCountAllAdminDataTablesCategories), likes)).Scan(&filtered).Error; err != nil {
			return result, err
		}

		if err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryCountAllAdminDataTablesCategories)).Scan(&total).Error; err != nil {
			return result, err
		}
	} else {
		if err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryCountAllAdminDataTablesCategories)).Scan(&filtered).Error; err != nil {
			return result, err
		}

		total = filtered
	}

	if order != "" {
		query = fmt.Sprintf("%s %s", query, order)
	}

	query = fmt.Sprintf("%s %s", query, limit)

	rows, err := repo.DB.Raw(helper.ConvertToInLineQuery(query)).Rows()

	if err != nil {
		return result, err
	}

	defer rows.Close()

	for rows.Next() {
		tmp := ContentCategory{}
		err := rows.Scan(
			&tmp.ID,
			&tmp.Category,
			&tmp.CreatedAt,
			&tmp.CreatedBy,
			&tmp.UpdatedAt,
			&tmp.UpdatedBy,
			&tmp.DeletedAt,
			&tmp.DeletedBy,
		)

		if err != nil {
			return result, err
		}

		data = append(data, map[string]any{
			"no":         no,
			"id":         tmp.ID,
			"category":   tmp.Category,
			"created_at": helper.HNTime(tmp.CreatedAt),
			"created_by": helper.HNString(tmp.CreatedBy),
			"updated_at": helper.HNTime(tmp.UpdatedAt),
			"updated_by": helper.HNString(tmp.UpdatedBy),
			"deleted_at": helper.HNTimeGDeletedAt(tmp.DeletedAt),
			"deleted_by": helper.HNString(tmp.DeletedBy),
		})

		no++
	}

	return helper.BuildDatatTables(data, filtered, total), nil
}

func (repo *repository) UserDataTablesContents(ctx *gin.Context, user user.User) (result helper.DataTables, err error) {
	var (
		query string = QueryUserDataTablesContents
		likes string = ""
		order string = ""
		limit string = ""
	)

	var (
		no       int = 1
		total    int = 0
		filtered int = 0
	)

	var data []map[string]any

	listOrder := []string{"", "title", "total_image", "status", ""}

	searchValue := ctx.Query("search[value]")
	orderColumn := ctx.Query("order[0][column]")
	starting, _ := strconv.Atoi(ctx.Query("start"))

	if searchValue != "" {
		likes = fmt.Sprintf("(title LIKE '%%%s%%' OR short_description LIKE '%%%s%%' OR status LIKE '%%%s%%')", searchValue, searchValue, searchValue)
	}

	if orderColumn != "" {
		orderType := ctx.Query("order[0][dir]")
		orderColumn, _ := strconv.Atoi(orderColumn)
		order = fmt.Sprintf("ORDER BY %s %s", listOrder[orderColumn], strings.ToUpper(orderType))
	} else {
		order = "ORDER BY id DESC"
	}

	if starting != -1 {
		length, _ := strconv.Atoi(ctx.Query("length"))
		limit = fmt.Sprintf("LIMIT %v OFFSET %v", length, starting)
		no = starting + 1
	}

	if likes != "" {
		query = fmt.Sprintf("%s AND %s", query, likes)

		if err := repo.DB.Raw(fmt.Sprintf("%s AND %s", helper.ConvertToInLineQuery(QueryCountAllUserDataTablesContents), likes), user.ID).Scan(&filtered).Error; err != nil {
			return result, err
		}

		if err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryCountAllUserDataTablesContents), user.ID).Scan(&total).Error; err != nil {
			return result, err
		}
	} else {
		if err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryCountAllUserDataTablesContents), user.ID).Scan(&filtered).Error; err != nil {
			return result, err
		}

		total = filtered
	}

	if order != "" {
		query = fmt.Sprintf("%s %s", query, order)
	}

	query = fmt.Sprintf("%s %s", query, limit)

	rows, err := repo.DB.Raw(helper.ConvertToInLineQuery(query), user.ID).Rows()

	if err != nil {
		return result, err
	}

	defer rows.Close()

	for rows.Next() {
		var (
			totalImage int = 0
			finishedAt sql.NullTime
			status     string
		)

		tmp := Content{}
		err := rows.Scan(
			&tmp.ID,
			&tmp.UserID,
			&tmp.CategoryID,
			&tmp.Title,
			&tmp.Slug,
			&tmp.ShortDescription,
			&tmp.Description,
			&totalImage,
			&status,
			&finishedAt,
			&tmp.CreatedAt,
			&tmp.CreatedBy,
			&tmp.UpdatedAt,
			&tmp.UpdatedBy,
			&tmp.DeletedAt,
			&tmp.DeletedBy,
		)

		if err != nil {
			return result, err
		}

		data = append(data, map[string]any{
			"no":                no,
			"id":                tmp.ID,
			"user_id":           tmp.UserID,
			"category_id":       tmp.CategoryID,
			"title":             tmp.Title,
			"slug":              tmp.Slug,
			"short_description": tmp.ShortDescription,
			"description":       tmp.Description,
			"total_image":       totalImage,
			"status":            status,
			"finished_at":       helper.HNTime(finishedAt),
			"created_at":        helper.HNTime(tmp.CreatedAt),
			"created_by":        helper.HNString(tmp.CreatedBy),
			"updated_at":        helper.HNTime(tmp.UpdatedAt),
			"updated_by":        helper.HNString(tmp.UpdatedBy),
			"deleted_at":        helper.HNTimeGDeletedAt(tmp.DeletedAt),
			"deleted_by":        helper.HNString(tmp.DeletedBy),
		})

		no++
	}

	return helper.BuildDatatTables(data, filtered, total), nil
}
