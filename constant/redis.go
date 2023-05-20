package constant

import "os"

var (
	REDIS_HOST string
	REDIS_PASS string
	REDIS_PORT string
)

func InitRedisConstant() {
	if APP_PROD {
		REDIS_HOST = os.Getenv("REDIS_HOST")
		REDIS_PASS = os.Getenv("REDIS_PASS")
		REDIS_PORT = os.Getenv("REDIS_PORT")
	} else {
		REDIS_HOST = os.Getenv("DEV_REDIS_HOST")
		REDIS_PASS = os.Getenv("DEV_REDIS_PASS")
		REDIS_PORT = os.Getenv("DEV_REDIS_PORT")
	}
}
