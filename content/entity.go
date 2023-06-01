package content

import (
	"time"

	"github.com/rickyromansyah2045/halocat-backend-go/constant"
)

type (
	ContentImage struct {
		ID           int    `json:"id"`
		ContentID    int    `json:"content_id"`
		FileLocation string `json:"file_location"`
		IsPrimary    int    `json:"is_primary"`
		constant.CreatedUpdatedDeleted
	}

	ContentCategory struct {
		ID       int    `json:"id"`
		Category string `json:"category"`
		constant.CreatedUpdatedDeleted
	}
	Content struct {
		ID               int       `json:"id"`
		UserID           int       `json:"user_id"`
		CategoryID       int       `json:"category_id"`
		Title            string    `json:"title"`
		Slug             string    `json:"slug"`
		ShortDescription string    `json:"short_description"`
		Description      string    `json:"description"`
		Status           string    `json:"status"`
		FinishedAt       time.Time `json:"finished_at"`
		ContentImages    []ContentImage
		constant.CreatedUpdatedDeleted
	}
)
