package constant

import (
	"os"
	"strconv"
)

var APP_PROD bool

func InitGeneralConstant() {
	APP_PROD, _ = strconv.ParseBool(os.Getenv("APP_PROD"))
}
