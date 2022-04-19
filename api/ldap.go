package api

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/rs/xid"
	"log"
	"ops-server/consts"
	"ops-server/db"
	"ops-server/entity"
	"ops-server/logs"
	"ops-server/models"
	"ops-server/service"
	"strings"
	"time"
)

func OperateLdap(g *echo.Group) {
	g.POST("/user/create", CreateLdapUser)
	g.DELETE("/user/delete", DeleteLdapUser)
	g.POST("/user/login", AuthLdapUser)
	g.GET("/health", HealthCheck)
	g.GET("/user/info", QueryUserInfo)
	g.POST("/user/logout", Logout)
	g.GET("/user/list-info", GetLdapUsersListInfo)

}
func HealthCheck(ctx echo.Context) error {
	return SuccessResp(ctx, nil)
}

func CreateLdapUser(ctx echo.Context) error {
	// receive create user http request params
	req := &entity.CreateUserParams{}
	if err := ctx.Bind(req); err != nil {
		logs.GetLogger().Errorf("CreateLdapUser req is failed reqParams is %s  err is %s\n", req, err.Error())
		return ErrorResp(ctx, consts.StatusText[consts.CodeLdapParamIsError], consts.CodeLdapParamIsError)
	}

	log.Printf("req is  %s\n", req)

	if req.Cn == "" || req.DisplayName == "" || req.GivenName == "" || len(req.EmployeeType) == 0 || req.UserPassword == "" || req.Mail == "" {
		logs.GetLogger().Errorf("req is failed reqParams is %s\n", req)
		return ErrorResp(ctx, consts.StatusText[consts.CodeLdapParamIsError], consts.CodeLdapParamIsError)
	}

	if err := db.GetLdap().CreateUser(req); err != nil {
		logs.GetLogger().Errorf("CreateUser is failed err is %s\n", err.Error())
		return ErrorResp(ctx, err.Error(), consts.CodeLdapCreateUserFailed)
	}

	//同步用户到mysql  用户表中
	uid, err := service.AddUser(models.User{UserName: req.Cn, Email: req.Mail, DisplayName: req.DisplayName})
	if err != nil {
		logs.GetLogger().Errorf("api CreateLdapUser AddUser is failed   err is %s\n", err.Error())
		return ErrorResp(ctx, consts.StatusText[consts.CodeInternalServerError], consts.CodeInternalServerError)
	}

	//默认guest 权限
	if err := service.CreateUserRoleRecord(models.UserRole{UserId: uid, RoleId: 5}); err != nil {
		logs.GetLogger().Errorf("api CreateLdapUser CreateUserRoleRecord    err is %s\n", err.Error())
		return ErrorResp(ctx, consts.StatusText[consts.CodeInternalServerError], consts.CodeInternalServerError)
	}

	if len(req.Role) != 0 {
		//todo add user role

		if err := service.AddUserRoles(uid, req.Role); err != nil {
			logs.GetLogger().Errorf("api CreateLdapUser  AddUserRoles   err is %s\n", err.Error())
			return ErrorResp(ctx, consts.StatusText[consts.CodeInternalServerError], consts.CodeInternalServerError)
		}
	}

	return SuccessResp(ctx, nil)

}

func DeleteLdapUser(ctx echo.Context) error {
	//create ldap connection
	req := new(entity.DeleteUserParams)
	if err := ctx.Bind(req); err != nil {
		logs.GetLogger().Errorf("DeleteLdapUser req is failed reqParams is %s  err is %s\n", req, err.Error())
		return ErrorResp(ctx, consts.StatusText[consts.CodeLdapParamIsError], consts.CodeLdapParamIsError)
	}

	if req.Cn == "" {
		logs.GetLogger().Errorf("DeleteLdapUser req.Cn is null   req is %s\n", req)
		return ErrorResp(ctx, consts.StatusText[consts.CodeLdapParamIsError], consts.CodeLdapParamIsError)
	}

	Dn := fmt.Sprintf("cn=%s,ou=person,dc=langzhihe,dc=com", req.Cn)

	if err := db.GetLdap().DeleteUser(Dn); err != nil {
		logs.GetLogger().Errorf("DeleteUser is failed err is %s\n", err.Error())
		return ErrorResp(ctx, consts.StatusText[consts.CodeLdapDeleteUserFailed], consts.CodeLdapDeleteUserFailed)
	}

	if err := service.DeleteUser(req.Cn); err != nil {

	}

	return SuccessResp(ctx, nil)
}

func AuthLdapUser(ctx echo.Context) error {
	req := new(entity.AuthUserParams)

	if err := ctx.Bind(req); err != nil {
		logs.GetLogger().Errorf("api AuthLdapUser req is failed reqParams is %s  err is %s\n", req, err.Error())
		return ErrorResp(ctx, consts.StatusText[consts.CodeLdapParamIsError], consts.CodeLdapParamIsError)
	}

	log.Printf("req is  %s\n", req)

	//根据参数 查找LDAP 用户是否存在
	result, err := db.GetLdap().SearchUser(req)
	if err != nil {
		logs.GetLogger().Errorf("api AuthLdapUser SearchUser ldap  failed  is  err is %s\n", err.Error())
		return ErrorResp(ctx, consts.StatusText[consts.CodeLdapSearchUserFailed], consts.CodeLdapSearchUserFailed)
	}

	fmt.Printf("result is %v\n", result)
	fmt.Printf("result len is %v\n", len(result.Entries))
	//result 为空查询不到ldap用户 返回登录失败
	if len(result.Entries) == 0 {
		logs.GetLogger().Errorf("api AuthLdapUser SearchUser ldap  can't find user    err is", )
		return ErrorResp(ctx, consts.StatusText[consts.CodeLdapUserNotExist], consts.CodeLdapUserNotExist)
	}

	ldapUserInfo := entity.LdapUserInfo{}

	for k, v := range result.Entries[0].Attributes {
		if v.Name == "cn" {
			ldapUserInfo.Cn = result.Entries[0].Attributes[k].Values[0]
		}

		if v.Name == "mail" {
			ldapUserInfo.Mail = result.Entries[0].Attributes[k].Values[0]
		}

		if v.Name == "displayName" {
			ldapUserInfo.DisPlayName = result.Entries[0].Attributes[k].Values[0]

		}
	}

	fmt.Printf("cn is : %v\n", ldapUserInfo)

	user, err := service.QueryUser(ldapUserInfo.Cn)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			uid, err := service.CreateUserRecord(models.User{UserName: ldapUserInfo.Cn, Email: ldapUserInfo.Mail, DisplayName: ldapUserInfo.DisPlayName})
			if err != nil {
				logs.GetLogger().Errorf("api AuthLdapUser CreateUserRecord    err is %s\n", err.Error())
				return ErrorResp(ctx, consts.StatusText[consts.CodeInternalServerError], consts.CodeInternalServerError)
			}
			user.Id = uid
			user.UserName = ldapUserInfo.Cn
			user.Email = ldapUserInfo.Mail
			user.DisplayName = ldapUserInfo.DisPlayName

			//每个用户给默认guest权限
			if err := service.CreateUserRoleRecord(models.UserRole{UserId: user.Id, RoleId: 5}); err != nil {
				logs.GetLogger().Errorf("api AuthLdapUser CreateUserRoleRecord    err is %s\n", err.Error())
				return ErrorResp(ctx, consts.StatusText[consts.CodeInternalServerError], consts.CodeInternalServerError)
			}

		}
	}

	if ldapUserInfo.Mail != user.Email || ldapUserInfo.DisPlayName != user.DisplayName {
		updates := map[string]interface{}{
			"email":        ldapUserInfo.Mail,
			"display_name": ldapUserInfo.DisPlayName,
		}

		if err := service.UpdateUser(user.Id, updates); err != nil {
			logs.GetLogger().Errorf("api AuthLdapUser UpdateUser is failed   err is %s\n", err.Error())
			return ErrorResp(ctx, consts.StatusText[consts.CodeInternalServerError], consts.CodeInternalServerError)
		}
	}

	//生成TOKEN 存入redis
	token := xid.New().String()

	log.Printf("api AuthLdapUser token  is %s   \n", token)

	logs.GetLogger().Info(fmt.Sprintf("api AuthLdapUser token  is %s   \n", token))

	//查出用户角色
	roleIDs, err := service.QueryUserRoleId(user.Id)
	if err != nil {
		logs.GetLogger().Errorf("api AuthLdapUser QueryUserRoleId is failed  err is %s   \n", err.Error())
		return ErrorResp(ctx, consts.StatusText[consts.CodeInternalServerError], consts.CodeInternalServerError)
	}

	roleNames, err := service.QueryUserRoles(roleIDs)
	if err != nil {
		logs.GetLogger().Errorf("api AuthLdapUser QueryUserRoles is failed  err is %s   \n", err.Error())
		return ErrorResp(ctx, consts.StatusText[consts.CodeInternalServerError], consts.CodeInternalServerError)
	}

	tmpMap := make(map[string]interface{})
	tmpMap["user_id"] = user.Id
	tmpMap["user_name"] = user.UserName
	tmpMap["email"] = user.Email

	if len(roleNames) != 0 {
		tmpMap["roles"] = strings.Join(roleNames, ",")
	}

	if err := db.RedisHMSet(token, tmpMap); err != nil {
		logs.GetLogger().Errorf("api AuthLdapUser RedisHMSet is failed   err is %s\n", err.Error())
		return ErrorResp(ctx, consts.StatusText[consts.CodeInternalServerError], consts.CodeInternalServerError)
	}

	if err := db.RedisSetKeyTtl(token, time.Minute*60); err != nil {
		logs.GetLogger().Errorf("api AuthLdapUser RedisSetKeyTtl is failed   err is %s\n", err.Error())
		return ErrorResp(ctx, consts.StatusText[consts.CodeInternalServerError], consts.CodeInternalServerError)
	}

	resp := entity.LoginAuthResp{
		Token:    token,
		UserId:   user.Id,
		UserName: user.UserName,
		Email:    user.Email,
		Role:     roleNames,
	}

	return SuccessResp(ctx, resp)

}

func QueryUserInfo(ctx echo.Context) error {

	token := ctx.Request().Header.Get("Authorization")
	userMapInfo, err := db.RedisHGetAll(token)

	if err != nil {
		logs.GetLogger().Errorf("api QueryUserInfo RedisHGetAll is failed   err is %s\n", err.Error())
		return ErrorResp(ctx, consts.StatusText[consts.CodeLdapParamIsError], consts.CodeLdapParamIsError)
	}

	v, ok := userMapInfo["roles"]

	if ok {
		roles := strings.Split(v, ",")
		resp := entity.UserInfo{
			Token:    token,
			UserId:   userMapInfo["user_id"],
			UserName: userMapInfo["user_name"],
			Email:    userMapInfo["email"],
			Role:     roles,
		}
		return SuccessResp(ctx, resp)

	}

	resp := entity.UserInfo{
		Token:    token,
		UserId:   userMapInfo["user_id"],
		UserName: userMapInfo["user_name"],
		Email:    userMapInfo["email"],
	}

	return SuccessResp(ctx, resp)

}

func Logout(ctx echo.Context) error {

	token := ctx.Request().Header.Get("Authorization")

	err := db.RedisDelKeys(token)
	if err != nil {
		return ErrorResp(ctx, consts.StatusText[consts.CodeUserDelTokenFailed], consts.CodeUserDelTokenFailed)
	}

	return SuccessResp(ctx, nil)

}

func GetLdapUsersListInfo(ctx echo.Context) error {
	req := &entity.UserInfoListRequest{}
	if err := ctx.Bind(req); err != nil {
		logs.GetLogger().Errorf("GetLdapUsersListInfo req is failed reqParams is %s  err is %s\n", req, err.Error())
		return ErrorResp(ctx, consts.StatusText[consts.CodeLdapParamIsError], consts.CodeLdapParamIsError)
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	fmt.Printf("req is %+v\n", req)

	usersListInfo, totalCount, err := service.QueryUserList(req)
	if err != nil {
		logs.GetLogger().Errorf("QueryUserList req is failed  is  err is %s\n", err.Error())
		return ErrorResp(ctx, consts.StatusText[consts.CodeInternalServerError], consts.CodeInternalServerError)
	}

	usersID := make([]uint64, 0, 0)

	for _, v := range usersListInfo {
		usersID = append(usersID, v.Id)
	}

	userMappingRole, err := service.QueryUserListAndRoles(usersID)

	if err != nil {
		logs.GetLogger().Errorf("QueryUserListAndRoles req is failed  is  err is %s\n", err.Error())
		return ErrorResp(ctx, consts.StatusText[consts.CodeInternalServerError], consts.CodeInternalServerError)
	}

	m1 := make(map[uint64]string)
	for _, v := range userMappingRole {
		m1[v.UserId] = v.Role
	}

	resp := &entity.UserInfoListResponse{}

	for _, v := range usersListInfo {
		resp.Items = append(resp.Items, entity.UserList{
			Id:          v.Id,
			UserName:    v.UserName,
			Email:       v.Email,
			DisplayName: v.DisplayName,
			Role:        m1[v.Id],
			CreatedAt:   v.CreatedAt,
		})
	}
	resp.TotalCount = totalCount

	return SuccessResp(ctx, resp)

}

//todo 重置密码
