package user

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rickyromansyah2045/halocat-backend-go/constant"
	"github.com/rickyromansyah2045/halocat-backend-go/helper"
)

func (repo *repository) SaveUser(user User) (User, error) {
	if err := repo.DB.Create(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (repo *repository) GetUserByEmail(email string) (user User, err error) {
	if err := repo.DB.Where("email = ?", email).Find(&user).Error; err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("sql: no rows in result set")
	}

	return user, nil
}

func (repo *repository) GetUserByID(id int) (user User, err error) {
	if err := repo.DB.Where("id = ?", id).Find(&user).Error; err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("sql: no rows in result set")
	}

	return user, nil
}

func (repo *repository) GetAllUser() (users []User, err error) {
	rows, err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryGetAllUser)).Rows()

	if err != nil {
		return users, err
	}

	defer rows.Close()

	for rows.Next() {
		tmp := User{}
		err := rows.Scan(
			&tmp.ID,
			&tmp.Role,
			&tmp.Name,
			&tmp.Email,
			&tmp.Password,
			&tmp.CreatedAt,
			&tmp.CreatedBy,
			&tmp.UpdatedAt,
			&tmp.UpdatedBy,
			&tmp.DeletedAt,
			&tmp.DeletedBy,
		)

		if err != nil {
			return users, err
		}

		users = append(users, tmp)
	}

	return users, nil
}

func (repo *repository) UpdateUser(user User) (User, error) {
	if err := repo.DB.Save(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (repo *repository) DeleteUser(user User) (bool, error) {
	if constant.DELETED_BY {
		if err := repo.DB.Save(&user).Error; err != nil {
			return false, err
		}
		return true, nil
	}

	if err := repo.DB.Delete(&user).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (repo *repository) AdminDataTablesUsers(ctx *gin.Context) (result helper.DataTables, err error) {
	var (
		query string = QueryAdminDataTablesUsers
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

	listOrder := []string{"", "name", "email", "role", "e_money", ""}

	searchValue := ctx.Query("search[value]")
	orderColumn := ctx.Query("order[0][column]")
	starting, _ := strconv.Atoi(ctx.Query("start"))

	if searchValue != "" {
		likes = fmt.Sprintf(`(
			name LIKE '%%%s%%' OR email LIKE '%%%s%%' OR role LIKE '%%%s%%' OR e_money LIKE '%%%s%%'
		)`, searchValue, searchValue, searchValue, searchValue)
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

		if err := repo.DB.Raw(fmt.Sprintf("%s AND %s", helper.ConvertToInLineQuery(QueryCountAllAdminDataTablesUsers), likes)).Scan(&filtered).Error; err != nil {
			return result, err
		}

		if err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryCountAllAdminDataTablesUsers)).Scan(&total).Error; err != nil {
			return result, err
		}
	} else {
		if err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryCountAllAdminDataTablesUsers)).Scan(&filtered).Error; err != nil {
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
		tmp := User{}
		err := rows.Scan(
			&tmp.ID,
			&tmp.Role,
			&tmp.Name,
			&tmp.Email,
			&tmp.Password,
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
			"role":       tmp.Role,
			"name":       tmp.Name,
			"email":      tmp.Email,
			"password":   tmp.Password,
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

func (repo *repository) GetUserRegistered(condition string) (res int, err error) {
	if err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryUserRegistered) + condition).Scan(&res).Error; err != nil {
		return res, err
	}
	return res, nil
}

func (repo *repository) CreateForgotPasswordToken(userForgotPasswordToken UserForgotPasswordToken) (UserForgotPasswordToken, error) {
	if err := repo.DB.Create(&userForgotPasswordToken).Error; err != nil {
		return userForgotPasswordToken, err
	}
	return userForgotPasswordToken, nil
}

func (repo *repository) GetDataForgotPasswordByToken(token string) (userForgotPasswordToken UserForgotPasswordToken, err error) {
	if err := repo.DB.Where("token = ?", token).Find(&userForgotPasswordToken).Error; err != nil {
		return userForgotPasswordToken, err
	}
	return userForgotPasswordToken, nil
}

func (repo *repository) DeleteForgotPasswordToken(userForgotPasswordToken UserForgotPasswordToken) (bool, error) {
	if err := repo.DB.Delete(&userForgotPasswordToken).Error; err != nil {
		return false, err
	}
	return true, nil
}
