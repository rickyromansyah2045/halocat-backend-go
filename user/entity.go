package user

import (
	"time"

	"github.com/rickyromansyah2045/halocat-backend-go/constant"
)

type (
	User struct {
		ID       int
		Role     string
		Name     string
		Email    string
		Password string
		constant.CreatedUpdatedDeleted
	}

	UserForgotPasswordToken struct {
		ID        int       `json:"id"`
		UserID    int       `json:"user_id"`
		Token     string    `json:"token"`
		CreatedAt time.Time `json:"created_at"`
	}
)
