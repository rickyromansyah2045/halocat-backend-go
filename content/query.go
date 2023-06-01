package content

const (
	QueryGetAll = `
		SELECT
			id,
			user_id,
			category_id,
			title,
			slug,
			short_description,
			description,
			status,
			finished_at,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM
			content
		WHERE
			deleted_at IS NULL
	`

	QueryGetContentByID = `
		SELECT
			id,
			user_id,
			category_id,
			title,
			slug,
			short_description,
			description,
			status,
			finished_at,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM
			content
		WHERE
			deleted_at IS NULL
		AND
			id = ?
		LIMIT
			1
	`

	QueryGetAllImage = `
		SELECT
			id,
			content_id,
			file_location,
			is_primary,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM
			content_images
		WHERE
			deleted_at IS NULL
	`

	QueryGetContentImageByID = `
		SELECT
			id,
			content_id,
			file_location,
			is_primary,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM
			content_images
		WHERE
			deleted_at IS NULL
		AND
			id = ?
		LIMIT
			1
	`

	QueryGetContentImages = `
		SELECT
			id,
			content_id,
			file_location,
			is_primary,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM
			content_images
		WHERE
			deleted_at IS NULL
		AND
			content_id = ?
		ORDER BY
			is_primary DESC
	`

	QueryGetAllCategory = `
		SELECT
			id,
			category,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM
			content_categories
		WHERE
			deleted_at IS NULL
	`

	QueryGetContentCategoryByID = `
		SELECT
			id,
			category,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM
			content_categories
		WHERE
			deleted_at IS NULL
		AND
			id = ?
		LIMIT
			1
	`

	QueryAdminDataTablesContents = `
		SELECT
			id,
			user_id,
			category_id,
			title,
			slug,
			short_description,
			description,
			(SELECT COUNT(id) FROM content_images WHERE deleted_at IS NULL AND content_id = content.id) AS total_image,
			status,
			finished_at,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM
			content
		WHERE
			deleted_at IS NULL
	`

	QueryCountAllAdminDataTablesContents = `
		SELECT
			COUNT(id) AS count_id
		FROM
			content
		WHERE
			deleted_at IS NULL
	`

	QueryAdminDataTablesCategories = `
		SELECT
			id,
			category,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM
			content_categories
		WHERE
			deleted_at IS NULL
	`

	QueryCountAllAdminDataTablesCategories = `
		SELECT
			COUNT(id) AS count_id
		FROM
			content_categories
		WHERE
			deleted_at IS NULL
	`

	QueryUserDataTablesContents = `
		SELECT
			id,
			user_id,
			category_id,
			title,
			slug,
			short_description,
			description,
			(SELECT COUNT(id) FROM content_images WHERE deleted_at IS NULL AND content_id = content.id) AS total_image,
			status,
			finished_at,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM
			content
		WHERE
			deleted_at IS NULL
		AND
			user_id = ?
	`

	QueryCountAllUserDataTablesContents = `
		SELECT
			COUNT(id) AS count_id
		FROM
			content
		WHERE
			deleted_at IS NULL
		AND
			user_id = ?
	`

	QueryGetTotalContent = QueryCountAllAdminDataTablesContents
)
