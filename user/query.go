package user

const (
	QueryGetAllUser = `
		SELECT
			id,
			role,
			name,
			email,
			password,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM
			users
		WHERE
			deleted_at IS NULL
	`

	QueryAdminDataTablesUsers = `
		SELECT
			id,
			role,
			name,
			email,
			password,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM
			users
		WHERE
			deleted_at IS NULL
	`

	QueryCountAllAdminDataTablesUsers = `
		SELECT
			COUNT(id) AS count_id
		FROM
			users
		WHERE
			deleted_at IS NULL
	`

	QueryUserRegistered = QueryCountAllAdminDataTablesUsers
)
