package constant

import (
	"os"
	"strconv"
)

var (
	DB_USER string
	DB_PASS string
	DB_NAME string
	RW_HOST string
	RO_HOST string
	DB_PORT string
)

var (
	DB_MAX_OPEN_CONNECTIONS int
	DB_MAX_IDLE_CONNECTIONS int
)

var (
	DB_CACHING bool
	DELETED_BY bool
	STRING_DSN string
)

func InitDBConstant() {
	if APP_PROD {
		dbCaching, _ := strconv.ParseBool(os.Getenv("DB_CACHING"))
		dbDeletedBy, _ := strconv.ParseBool(os.Getenv("DELETED_BY"))
		dbMaxOpenConns, _ := strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONNECTIONS"))
		dbMaxIdleConns, _ := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNECTIONS"))

		DB_USER = os.Getenv("DB_USER")
		DB_PASS = os.Getenv("DB_PASS")
		DB_NAME = os.Getenv("DB_NAME")
		RW_HOST = os.Getenv("RW_HOST")
		RO_HOST = os.Getenv("RW_HOST")
		DB_PORT = os.Getenv("DB_PORT")

		DB_CACHING = dbCaching
		DB_MAX_OPEN_CONNECTIONS = dbMaxOpenConns
		DB_MAX_IDLE_CONNECTIONS = dbMaxIdleConns
		DELETED_BY = dbDeletedBy
	} else {
		dbCaching, _ := strconv.ParseBool(os.Getenv("DEV_DB_CACHING"))
		dbDeletedBy, _ := strconv.ParseBool(os.Getenv("DEV_DELETED_BY"))
		dbMaxOpenConns, _ := strconv.Atoi(os.Getenv("DEV_DB_MAX_OPEN_CONNECTIONS"))
		dbMaxIdleConns, _ := strconv.Atoi(os.Getenv("DEV_DB_MAX_IDLE_CONNECTIONS"))

		DB_USER = os.Getenv("DEV_DB_USER")
		DB_PASS = os.Getenv("DEV_DB_PASS")
		DB_NAME = os.Getenv("DEV_DB_NAME")
		RW_HOST = os.Getenv("DEV_RW_HOST")
		RO_HOST = os.Getenv("DEV_RW_HOST")
		DB_PORT = os.Getenv("DEV_DB_PORT")

		DB_CACHING = dbCaching
		DB_MAX_OPEN_CONNECTIONS = dbMaxOpenConns
		DB_MAX_IDLE_CONNECTIONS = dbMaxIdleConns
		DELETED_BY = dbDeletedBy
	}

	STRING_DSN = os.Getenv("STRING_DSN")
}
