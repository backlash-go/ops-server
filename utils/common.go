package utils

import (
	"ops-server/logs"
	"ops-server/service"
)

func FilterPermission(roles []string, url string) (bool, error) {

	roleIDs, err := service.QueryUserRoleIdByRoleName(roles)

	if err != nil {
		logs.GetLogger().Errorf("service.QueryUserRoleIdByRoleName   err is %s\n", err.Error())
		return false, err
	}

	if len(roleIDs) == 0 {
		return false, nil
	}

	permissionId, err := service.QueryPermissionIdByRoleId(roleIDs)

	if err != nil {
		logs.GetLogger().Errorf("service.QueryPermissionIdByRoleId  err is %s \n", err.Error())
		return false, err
	}

	if len(permissionId) == 0 {
		return false, nil
	}

	apis, err := service.QueryApi(permissionId)

	if err != nil {
		logs.GetLogger().Errorf("service QueryApi err is %s \n", err.Error())
		return false, err
	}

	if len(apis) == 0 {
		return false, nil
	}

	for _, v := range apis {
		if v == url {
			return true, nil
		}
	}

	return false, nil

}
