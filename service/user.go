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

func QueryAllUser(userName string) (user []models.User, err error) {

	err = db.GetDB().Model(&models.User{}).Where("user_name = ?", userName).Find(&user).Error
	return
}

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

func QueryUserRoleId(userId uint64) (roleIDs []uint64, err error) {
	err = db.GetDB().Model(&models.UserRole{}).Where("user_id = ?", userId).Pluck("role_id", &roleIDs).Error
	return
}

func QueryUserRoles(roleId []uint64) (roleNames []string, err error) {
	err = db.GetDB().Model(&models.Role{}).Where("id in (?)", roleId).Pluck("role_name", &roleNames).Error
	return
}

func QueryUserRoleIdByRoleName(roleName []string) (roleIDs []uint64, err error) {
	err = db.GetDB().Model(&models.Role{}).Where("role_name in (?)", roleName).Pluck("id", &roleIDs).Error
	return
}

func QueryPermissionIdByRoleId(roleId []uint64) (permissionId []uint64, err error) {
	err = db.GetDB().Model(&models.RolePermission{}).Where("role_id in (?)", roleId).Pluck("permission_id", &permissionId).Error
	return
}

func QueryApi(permissionId []uint64) (apis []string, err error) {
	err = db.GetDB().Model(&models.Permission{}).Where("id in (?)", permissionId).Pluck("api", &apis).Error
	return
}

func CreateUserRecord(user models.User) (uint64, error) {
	err := db.GetDB().Model(&models.User{}).Create(&user).Error
	return user.Id, err
}

func CreateUserRoleRecord(userRole models.UserRole) error {
	err := db.GetDB().Model(&models.User{}).Create(&userRole).Error
	return err

}
