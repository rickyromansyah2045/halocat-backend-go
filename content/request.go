package content

import (
	"github.com/rickyromansyah2045/halocat-backend-go/user"
)

type (
	RequestCreateContent struct {
		UserID           int    `json:"user_id"`
		CategoryID       int    `json:"category_id" binding:"required"`
		Title            string `json:"title" binding:"required"`
		ShortDescription string `json:"short_description" binding:"required"`
		Description      string `json:"description" binding:"required"`
		FinishedAt       string `json:"finished_at"`
		Status           string `json:"status" binding:"required"`
		User             user.User
	}

	RequestUpdateContent struct {
		RequestCreateContent
	}

	RequestDeleteContent struct {
		User user.User
	}

	RequestGetContentByID struct {
		ID int `uri:"id" binding:"required"`
	}

	RequestGetContentImageByID struct {
		RequestGetContentByID
	}

	RequestGetContentCategoryByID struct {
		RequestGetContentByID
	}

	RequestCreateContentImage struct {
		ContentID int  `form:"content_id" binding:"required"`
		IsPrimary bool `form:"is_primary"`
		User      user.User
	}

	RequestDeleteContentImage struct {
		User user.User
	}

	RequestDeleteContentCategory struct {
		User user.User
	}

	RequestCreateContentCategory struct {
		Category string `json:"category" binding:"required"`
		User     user.User
	}

	RequestUpdateContentCategory struct {
		RequestCreateContentCategory
	}
)
