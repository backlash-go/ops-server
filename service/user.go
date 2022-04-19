package service

import (
	"fmt"
	"ops-server/db"
	"ops-server/entity"
	"ops-server/models"
)

//判断用户是否存在
//func QueryUser(cn string) (user models.User, err error) {
//	err = db.GetDB().Where("user_name = ?", cn).First(&user).Error
//	return
//}

func DeleteUser(cn string) error {
	err := db.GetDB().Where("user_name = ?", cn).Delete(&models.User{}).Error
	return err
}

func QueryUserListAndRoles(userID []uint64) (userRole []entity.UserIDRoleContact, err error) {

	//sql := `SELECT user_id, group_concat(distinct role_name) as role FROM ops.user_role join ops.role on ops.user_role.role_id = ops.role.id   group by user_id`

	err = db.GetDB().Model(&models.UserRole{}).Select("user_id, group_concat(distinct role_name) as role").
		Joins("join role on user_role.role_id = role.id").Group("user_id").Having("user_id in (?)", userID).Scan(&userRole).Error
	return
}

func QueryUserList(req *entity.UserInfoListRequest) (users []entity.UserList, totalCount int64, err error) {

	fmt.Println(req.SearchName)

	if len(req.SearchName) > 0 {
		err = db.GetDB().Model(&models.User{}).Offset(req.PageSize*(req.Page-1)).Limit(req.PageSize).Where("user_name LIKE ?", "%"+req.SearchName+"%").Scan(&users).Error
		if err != nil {
			return
		}

		err = db.GetDB().Model(&models.User{}).Where("user_name LIKE ?", "%"+req.SearchName+"%").Count(&totalCount).Error
		if err != nil {
			return
		}
		return
	}

	err = db.GetDB().Model(&models.User{}).Count(&totalCount).Error
	if err != nil {
		return
	}

	err = db.GetDB().Model(&models.User{}).Offset(req.PageSize * (req.Page - 1)).Limit(req.PageSize).Scan(&users).Error
	if err != nil {
		return
	}
	return
}

func QueryAllUser(userName string) (user []models.User, err error) {
	err = db.GetDB().Model(&models.User{}).Where("user_name = ?", userName).Find(&user).Error
	return
}

func QueryUser(cn string) (user models.User, err error) {
	err = db.GetDB().Model(&models.User{}).Where("user_name = ?", cn).First(&user).Error
	return
}

func UpdateUser(userId uint64, updates map[string]interface{}) error {
	return db.GetDB().Model(&models.User{}).Where("id = ?", userId).Updates(updates).Error
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
	err := db.GetDB().Create(&userRole).Error
	return err
}

func AddUserRoles(uid uint64, roles []uint64) error {
	userRoles := make([]models.UserRole, 0, 0)
	for _, v := range roles {
		userRole := models.UserRole{RoleId: v, UserId: uid}
		userRoles = append(userRoles, userRole)
	}
	err := db.GetDB().Model(&models.UserRole{}).Create(&userRoles).Error
	return err
}
