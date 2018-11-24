package config

import (
	"os"
)

func SetEnv() {
	os.Setenv("MYSQL_URI", `h@tcp(127.0.0.1:3306)/freeze_app`)
	os.Setenv("PORT", `:8080`)
	os.Setenv("SECRET_KEY", `a452gaaagasdfakl4rq5j1lk45j1lkjrajfaoiwj45lk45`)
}
