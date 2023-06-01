package content

import "time"

type (
	ContentFormatter struct {
		ID               int                                     `json:"id"`
		UserID           int                                     `json:"user_id"`
		CategoryID       int                                     `json:"category_id"`
		Title            string                                  `json:"title"`
		Slug             string                                  `json:"slug"`
		ShortDescription string                                  `json:"short_description"`
		Description      string                                  `json:"description"`
		Status           string                                  `json:"status"`
		FinishedAt       time.Time                               `json:"finished_at"`
		ContentImages    []ContentImageWithoutContentIDFormatter `json:"images"`
	}

	ContentImageFormatter struct {
		ID           int    `json:"id"`
		ContentID    int    `json:"content_id"`
		FileLocation string `json:"file_location"`
		IsPrimary    int    `json:"is_primary"`
	}

	ContentImageWithoutContentIDFormatter struct {
		ID           int    `json:"id"`
		FileLocation string `json:"file_location"`
		IsPrimary    int    `json:"is_primary"`
	}

	ContentCategoryFormatter struct {
		ID       int    `json:"id"`
		Category string `json:"category"`
	}
)

func FormatContentData(content Content) (response ContentFormatter) {
	response = ContentFormatter{
		ID:               content.ID,
		UserID:           content.UserID,
		CategoryID:       content.CategoryID,
		Title:            content.Title,
		Slug:             content.Slug,
		ShortDescription: content.ShortDescription,
		Description:      content.Description,
		Status:           content.Status,
		FinishedAt:       content.FinishedAt,
	}

	images := []ContentImageWithoutContentIDFormatter{}
	tmpImages := ContentImageWithoutContentIDFormatter{}

	for _, img := range content.ContentImages {
		tmpImages.ID = img.ID
		tmpImages.FileLocation = img.FileLocation
		tmpImages.IsPrimary = img.IsPrimary

		images = append(images, tmpImages)
	}

	response.ContentImages = images

	return response
}

func FormatMultipleContentData(contents []Content) (response []ContentFormatter) {
	tmp := ContentFormatter{}
	tmpImages := ContentImageWithoutContentIDFormatter{}

	for _, val := range contents {
		tmp.ID = val.ID
		tmp.UserID = val.UserID
		tmp.CategoryID = val.CategoryID
		tmp.Title = val.Title
		tmp.Slug = val.Slug
		tmp.ShortDescription = val.ShortDescription
		tmp.Description = val.Description
		tmp.Status = val.Status
		tmp.FinishedAt = val.FinishedAt

		images := []ContentImageWithoutContentIDFormatter{}

		for _, img := range val.ContentImages {
			tmpImages.ID = img.ID
			tmpImages.FileLocation = img.FileLocation
			tmpImages.IsPrimary = img.IsPrimary

			images = append(images, tmpImages)
		}

		tmp.ContentImages = images

		response = append(response, tmp)
	}

	if len(response) == 0 {
		return []ContentFormatter{}
	}

	return response
}

func FormatContentImageData(content ContentImage) (response ContentImageFormatter) {
	response = ContentImageFormatter{
		ID:           content.ID,
		ContentID:    content.ContentID,
		FileLocation: content.FileLocation,
		IsPrimary:    content.IsPrimary,
	}

	return response
}

func FormatMultipleContentImageData(contentImages []ContentImage) (response []ContentImageFormatter) {
	for _, val := range contentImages {
		tmp := ContentImageFormatter{}
		tmp.ID = val.ID
		tmp.ContentID = val.ContentID
		tmp.FileLocation = val.FileLocation
		tmp.IsPrimary = val.IsPrimary

		response = append(response, tmp)
	}

	if len(response) == 0 {
		return []ContentImageFormatter{}
	}

	return response
}

func FormatContentCategoryData(category ContentCategory) (response ContentCategoryFormatter) {
	response = ContentCategoryFormatter{
		ID:       category.ID,
		Category: category.Category,
	}

	return response
}

func FormatMultipleContentCategoryData(contentCategories []ContentCategory) (response []ContentCategoryFormatter) {
	for _, val := range contentCategories {
		tmp := ContentCategoryFormatter{}
		tmp.ID = val.ID
		tmp.Category = val.Category

		response = append(response, tmp)
	}

	if len(response) == 0 {
		return []ContentCategoryFormatter{}
	}

	return response
}
