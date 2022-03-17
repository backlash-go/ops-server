package service

import (
	"ops-server/db"
	"ops-server/models"
)

//判断用户是否存在
//func QueryUser(cn string) (user models.User, err error) {
//	err = db.GetDB().Where("user_name = ?", cn).First(&user).Error
//	return
//}

func QueryUser(cn string) (user models.User, err error) {
	err = db.GetDB().Model(&models.User{}).Where("user_name = ?", cn).First(&user).Error
	return

}

func UpdateUser(userId uint64, updates interface{}) error {
	return db.GetDB().Model(&models.User{}).Where("id = ?", userId).Update(updates).Error
}

func AddUser(req models.User) (uint64, error) {
	 err := db.GetDB().Model(&models.User{}).Create(&req).Error
	return req.Id, err
}
