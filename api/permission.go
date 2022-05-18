package api

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"ops-server/consts"
	"ops-server/entity"
	"ops-server/logs"
	"ops-server/models"
	"ops-server/service"
	"strings"
)

func OperatePermission(g *echo.Group) {
	g.GET("/list-info", GetPermissionListInfo)
	g.POST("/create", CreateApi)
	g.DELETE("/delete", DeleteApi)
	g.POST("/modify", UpdateApiInfo)
	g.GET("/info",QueryPermissionInfo)
}

func QueryPermissionInfo(ctx echo.Context) error {

	req := &entity.PermissionInfoRequest{}

	if err := ctx.Bind(req); err != nil {
		logs.GetLogger().Errorf("UpdateApiInfo req is failed reqParams is %s  err is %s\n", req, err.Error())
		return ErrorResp(ctx, consts.StatusText[consts.CodeParamIsError], consts.CodeParamIsError)
	}
	log.Println(req.ID)

	permissionInfo, err := service.QueryPermission(req.ID)

	if err != nil {
		logs.GetLogger().Errorf("api QueryPermissionInfo QueryPermission is failed   err is %s\n", err.Error())
		return ErrorResp(ctx, consts.StatusText[consts.CodeInternalServerError], consts.CodeInternalServerError)
	}

	roleIds, err := service.QueryPermissionRole(req.ID)
	if err != nil {
		logs.GetLogger().Errorf("api QueryPermissionInfo QueryPermissionRole is failed   err is %s\n", err.Error())
		return ErrorResp(ctx, consts.StatusText[consts.CodeInternalServerError], consts.CodeInternalServerError)
	}

	roleNames, err := service.QueryPermissionRoleName(roleIds)

	if err != nil {
		logs.GetLogger().Errorf("api QueryPermissionInfo QueryPermissionRoleName is failed   err is %s\n", err.Error())
		return ErrorResp(ctx, consts.StatusText[consts.CodeInternalServerError], consts.CodeInternalServerError)
	}

	resp := entity.PermissionInfoList{Id: permissionInfo.Id, Name: permissionInfo.Name, Api: permissionInfo.Api, Role: strings.Join(roleNames, ",")}

	return SuccessResp(ctx, resp)
}

func UpdateApiInfo(ctx echo.Context) error {
	req := &entity.UpdateApiParamsRequest{}

	if err := ctx.Bind(req); err != nil {
		logs.GetLogger().Errorf("UpdateApiInfo req is failed reqParams is %s  err is %s\n", req, err.Error())
		return ErrorResp(ctx, consts.StatusText[consts.CodeParamIsError], consts.CodeParamIsError)
	}
	logs.GetLogger().Info(fmt.Sprintf("api UpdateApiInfo req  is %+v   \n", req))

	updates := map[string]interface{}{
		"api":  req.Api,
		"name": req.Name,
	}

	if err := service.UpdateApi(req.ID, updates); err != nil {
		logs.GetLogger().Errorf("api UpdateApiInfo UpdateApi is failed   err is %s\n", err.Error())
		return ErrorResp(ctx, consts.StatusText[consts.CodeInternalServerError], consts.CodeInternalServerError)
	}

	if err := service.DeleteApiRoleById(req.ID); err != nil {
		logs.GetLogger().Errorf("UpdateApiInfo service DeleteApiRoleById   failed  is  err is %s\n", err.Error())
		return ErrorResp(ctx, consts.StatusText[consts.CodeInternalServerError], consts.CodeInternalServerError)
	}

	roleIDs, err := service.QueryUserRoleIdByRoleName(req.Role)
	if err != nil {
		logs.GetLogger().Errorf("UpdateApiInfo service QueryUserRoleIdByRoleName   failed  is  err is %s\n", err.Error())
		return ErrorResp(ctx, consts.StatusText[consts.CodeInternalServerError], consts.CodeInternalServerError)
	}

	if err := service.AddApiRoles(req.ID, roleIDs); err != nil {
		logs.GetLogger().Errorf("UpdateApiInfo service AddApiRoles   failed  is  err is %s\n", err.Error())
		return ErrorResp(ctx, consts.StatusText[consts.CodeInternalServerError], consts.CodeInternalServerError)
	}

	return SuccessResp(ctx, nil)

}

func DeleteApi(ctx echo.Context) error {

	req := &entity.PermissionInfoList{}

	if err := ctx.Bind(req); err != nil {
		logs.GetLogger().Errorf("DeleteApi req is failed reqParams is %s  err is %s\n", req, err.Error())
		return ErrorResp(ctx, consts.StatusText[consts.CodeParamIsError], consts.CodeParamIsError)
	}
	logs.GetLogger().Info(fmt.Sprintf("api DeleteApi req  is %+v   \n", req))

	if req.Id == 0 {
		return ErrorResp(ctx, consts.StatusText[consts.CodeParamIsError], consts.CodeParamIsError)
	}

	err := service.DeleteApi(req.Id)

	if err != nil {
		logs.GetLogger().Errorf("DeleteApi is failed  is  err is %s\n", err.Error())
		return ErrorResp(ctx, consts.StatusText[consts.CodeInternalServerError], consts.CodeInternalServerError)
	}

	err = service.DeleteApiPower(req.Id)

	if err != nil {
		logs.GetLogger().Errorf("DeleteApiPower is failed  is  err is %s\n", err.Error())
		return ErrorResp(ctx, consts.StatusText[consts.CodeInternalServerError], consts.CodeInternalServerError)
	}

	return SuccessResp(ctx, nil)
}

func CreateApi(ctx echo.Context) error {

	req := &entity.CreateApiParamsRequest{}
	if err := ctx.Bind(req); err != nil {
		logs.GetLogger().Errorf("CreateApi req is failed reqParams is %s  err is %s\n", req, err.Error())
		return ErrorResp(ctx, consts.StatusText[consts.CodeParamIsError], consts.CodeParamIsError)
	}
	logs.GetLogger().Info(fmt.Sprintf("api CreateApi req  is %+v   \n", req))

	permissionId, err := service.CreateApiRecord(models.Permission{Api: req.Api, Name: req.Name})

	if err != nil {
		logs.GetLogger().Errorf("CreateApiRecord req is failed  is  err is %s\n", err.Error())
		return ErrorResp(ctx, consts.StatusText[consts.CodeInternalServerError], consts.CodeInternalServerError)
	}

	err = service.AddApiRoles(permissionId, req.Role)
	if err != nil {
		logs.GetLogger().Errorf("CreateApiRecord req is failed  is  err is %s\n", err.Error())
		return ErrorResp(ctx, consts.StatusText[consts.CodeInternalServerError], consts.CodeInternalServerError)
	}

	return SuccessResp(ctx, nil)
}

func GetPermissionListInfo(ctx echo.Context) error {

	req := &entity.PermissionInfoListRequest{}
	if err := ctx.Bind(req); err != nil {
		logs.GetLogger().Errorf("GetPermissionListInfo req is failed reqParams is %s  err is %s\n", req, err.Error())
		return ErrorResp(ctx, consts.StatusText[consts.CodeParamIsError], consts.CodeParamIsError)
	}

	logs.GetLogger().Info(fmt.Sprintf("api GetPermissionListInfo req  is %+v   \n", req))

	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.Page <= 0 {
		req.Page = 1
	}

	permissionListInfo, totalCount, err := service.QueryPermissionList(req)
	if err != nil {
		logs.GetLogger().Errorf("QueryUserList req is failed  is  err is %s\n", err.Error())
		return ErrorResp(ctx, consts.StatusText[consts.CodeInternalServerError], consts.CodeInternalServerError)
	}

	permissionIds := make([]uint64, 0, 0)

	for _, v := range permissionListInfo {
		permissionIds = append(permissionIds, v.Id)
	}

	permissionMapRole, err := service.QueryPermissionListAndRoles(permissionIds)

	if err != nil {
		logs.GetLogger().Errorf("QueryPermissionListAndRoles req is failed  is  err is %s\n", err.Error())
		return ErrorResp(ctx, consts.StatusText[consts.CodeInternalServerError], consts.CodeInternalServerError)
	}

	m1 := make(map[uint64]string)
	for _, v := range permissionMapRole {
		m1[v.PermissionId] = v.Role
	}

	resp := &entity.PermissionInfoListResponse{}

	for _, v := range permissionListInfo {
		resp.Items = append(resp.Items, entity.PermissionInfoList{
			Id:        v.Id,
			Api:       v.Api,
			Name:      v.Name,
			Role:      m1[v.Id],
			CreatedAt: v.CreatedAt,
		})
	}
	resp.TotalCount = totalCount

	return SuccessResp(ctx, resp)

}
