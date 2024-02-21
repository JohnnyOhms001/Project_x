package config

import "github.com/JohnnyOhms/projectx/model"

func SyncDB() {
	DB.AutoMigrate(&model.User{}, &model.User_Info{}, &model.Avatar{})
}
