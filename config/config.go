package config

import (
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() error {
	var err error
	dsn := os.Getenv("DB_String")
	dsn += "?parseTime=true&loc=Asia%2FShanghai" // Set the time zone to Asia/Shanghai
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to DB")
	}
	return nil
}
