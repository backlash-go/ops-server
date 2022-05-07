package service

import (
	"ops-server/db"
	"ops-server/entity"
	"ops-server/models"
)


func DeleteApiRoleById(Id uint64) error {
	err := db.GetDB().Where("permission_id = ?", Id).Delete(&models.RolePermission{}).Error
	return err
}

func UpdateApi(permissionId uint64, updates map[string]interface{}) error {
	return db.GetDB().Model(&models.Permission{}).Where("id = ?", permissionId).Updates(updates).Error
}

func DeleteApi(id uint64) error {
	err := db.GetDB().Where("id = ?", id).Delete(&models.Permission{}).Error
	return err
}

func DeleteApiPower(id uint64) error {
	err := db.GetDB().Where("id = ?", id).Delete(&models.RolePermission{}).Error
	return err
}

func CreateApiRecord(permission models.Permission) (uint64, error) {
	err := db.GetDB().Model(&models.Permission{}).Create(&permission).Error
	return permission.Id, err
}


func AddApiRoles(permissionId uint64, roles []uint64) error {
	permissionRoles := make([]models.RolePermission, 0, 0)
	for _, v := range roles {
		permissionRole := models.RolePermission{RoleId: v, PermissionId: permissionId}
		permissionRoles = append(permissionRoles, permissionRole)
	}
	err := db.GetDB().Model(&models.RolePermission{}).Create(&permissionRoles).Error
	return err
}



func QueryPermissionList(req *entity.PermissionInfoListRequest) (permissionInfo []entity.PermissionInfoList, totalCount int64, err error) {

	if len(req.SearchName) > 0 {
		err = db.GetDB().Model(&models.Permission{}).Offset(req.PageSize*(req.Page-1)).Limit(req.PageSize).Where("api LIKE ?", "%"+req.SearchName+"%").Scan(&permissionInfo).Error
		if err != nil {
			return
		}
		err = db.GetDB().Model(&models.Permission{}).Where("api LIKE ?", "%"+req.SearchName+"%").Count(&totalCount).Error
		if err != nil {
			return
		}
		return
	}

	err = db.GetDB().Model(&models.Permission{}).Count(&totalCount).Error
	if err != nil {
		return
	}

	err = db.GetDB().Model(&models.Permission{}).Offset(req.PageSize * (req.Page - 1)).Limit(req.PageSize).Scan(&permissionInfo).Error
	if err != nil {
		return
	}
	return
}

func QueryPermissionListAndRoles(permissionId []uint64) (permissionRole []entity.PermissionIdRoleContact, err error) {
	//SELECT permission_id, group_concat(distinct role_name) as role FROM `role_permission` join role on role_permission.role_id = role.id WHERE `role_permission`.`deleted_at` IS NULL GROUP BY `permission_id` HAVING permission_id in (10);


	err = db.GetDB().Model(&models.RolePermission{}).Select("permission_id, group_concat(distinct role_name) as role").
		Joins("join role on role_permission.role_id = role.id").Group("permission_id").Having("permission_id in (?)", permissionId).Scan(&permissionRole).Error
	return
}
