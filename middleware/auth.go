package middleware

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/rickyromansyah2045/halocat-backend-go/auth"
	"github.com/rickyromansyah2045/halocat-backend-go/helper"
	"github.com/rickyromansyah2045/halocat-backend-go/user"
)

func Auth(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, helper.BasicAPIResponseError(http.StatusUnauthorized, "Unauthorized, token not found!"))
			return
		}

		tokenString, arrayToken := "", strings.Split(authHeader, " ")

		if len(arrayToken) != 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, helper.BasicAPIResponseError(http.StatusUnauthorized, "Unauthorized, invalid token!"))
			return
		}

		tokenString = arrayToken[1]
		token, err := authService.ValidateToken(tokenString)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, helper.BasicAPIResponseError(http.StatusUnauthorized, "Unauthorized, invalid token!"))
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, helper.BasicAPIResponseError(http.StatusUnauthorized, "Unauthorized, invalid token!"))
			return
		}

		userID := int(claim["the_cloud_donation_user_id"].(float64))
		user, err := userService.GetUserByID(userID)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, helper.BasicAPIResponseError(http.StatusUnauthorized, "Unauthorized, invalid token!"))
			return
		}

		ctx.Set("userData", user)
	}
}

func AdminAuth(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, helper.BasicAPIResponseError(http.StatusUnauthorized, "Unauthorized, token not found!"))
			return
		}

		tokenString, arrayToken := "", strings.Split(authHeader, " ")

		if len(arrayToken) != 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, helper.BasicAPIResponseError(http.StatusUnauthorized, "Unauthorized, invalid token!"))
			return
		}

		tokenString = arrayToken[1]
		token, err := authService.ValidateToken(tokenString)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, helper.BasicAPIResponseError(http.StatusUnauthorized, "Unauthorized, invalid token!"))
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, helper.BasicAPIResponseError(http.StatusUnauthorized, "Unauthorized, invalid token!"))
			return
		}

		userID := int(claim["the_cloud_donation_user_id"].(float64))
		user, err := userService.GetUserByID(userID)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, helper.BasicAPIResponseError(http.StatusUnauthorized, "Unauthorized, invalid token!"))
			return
		}

		if user.Role == "user" {
			ctx.AbortWithStatusJSON(http.StatusForbidden, helper.BasicAPIResponseError(http.StatusForbidden, "You cannot access this endpoint!"))
			return
		}

		ctx.Set("userData", user)
	}
}
