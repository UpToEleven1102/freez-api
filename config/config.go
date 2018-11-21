package config

import "os"

func ConfigEnv() {
	os.Setenv("MYSQL_URI", `root@tcp(127.0.0.1:3306)/freeze_app`)
	os.Setenv("PORT", `:8080`)
}
