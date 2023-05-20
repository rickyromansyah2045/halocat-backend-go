package user

type (
	RequestRegister struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	RequestLogin struct {
		Email    string `json:"email" form:"email" binding:"required,email"`
		Password string `json:"password" form:"password" binding:"required"`
	}
)

type (
	RequestGetUserByID struct {
		ID int `uri:"id" binding:"required"`
	}

	RequestCreateUser struct {
		Role     string `json:"role" binding:"required"`
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
		User     User
	}

	RequestUpdateUser struct {
		Role     string `json:"role" binding:"required"`
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password"`
		User     User
	}

	RequestSelfUpdateUser struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password"`
		User     User
	}

	RequestDeleteUser struct {
		User User
	}

	RequestCreateForgotPasswordToken struct {
		Email string `json:"email" binding:"required"`
		User  User
	}

	RequestProcessForgotPasswordToken struct {
		Token string `uri:"token" binding:"required"`
		User  User
	}
)
