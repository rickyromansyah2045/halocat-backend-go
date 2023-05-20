package user

type (
	UserFormatter struct {
		ID    int    `json:"id"`
		Role  string `json:"role"`
		Name  string `json:"name"`
		Email string `json:"email"`
		Token string `json:"token"`
	}

	UserListFormatter struct {
		ID    int    `json:"id"`
		Role  string `json:"role"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
)

func FormatUserData(user User, token string) UserFormatter {
	formatData := UserFormatter{
		ID:    user.ID,
		Role:  user.Role,
		Name:  user.Name,
		Email: user.Email,
		Token: token,
	}

	return formatData
}

func FormatUserFullData(user User) UserListFormatter {
	formatData := UserListFormatter{
		ID:    user.ID,
		Role:  user.Role,
		Name:  user.Name,
		Email: user.Email,
	}

	return formatData
}

func FormatListUserData(users []User) (response []UserListFormatter) {
	for _, val := range users {
		tmp := UserListFormatter{}
		tmp.ID = val.ID
		tmp.Role = val.Role
		tmp.Name = val.Name
		tmp.Email = val.Email

		response = append(response, tmp)
	}

	if len(response) == 0 {
		return []UserListFormatter{}
	}

	return response
}
