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
	g.POST("/user/auth", AuthLdapUser)
	g.GET("/health", HealthCheck)
	g.GET("/user/info", QueryUserInfo)

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
		return ErrorResp(ctx, consts.StatusText[consts.CodeLdapCreateUserFailed], consts.CodeLdapCreateUserFailed)
	}

	//同步用户到mysql  用户表中
	_, err := service.AddUser(models.User{UserName: req.Cn, Email: req.Mail})
	if err != nil {
		logs.GetLogger().Errorf("api AuthLdapUser AddUser is failed   err is %s\n", err.Error())
		return ErrorResp(ctx, consts.StatusText[consts.CodeInternalServerError], consts.CodeInternalServerError)
	}

	if len(req.Role) != 0 {
		//todo add user role
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

	if req.Dn == "" {
		logs.GetLogger().Errorf("DeleteLdapUser req.Dn is null   req is %s\n", req)
		return ErrorResp(ctx, consts.StatusText[consts.CodeLdapParamIsError], consts.CodeLdapParamIsError)
	}

	if err := db.GetLdap().DeleteUser(req); err != nil {
		logs.GetLogger().Errorf("DeleteUser is failed err is %s\n", err.Error())
		return ErrorResp(ctx, consts.StatusText[consts.CodeLdapDeleteUserFailed], consts.CodeLdapDeleteUserFailed)
	}

	return SuccessResp(ctx, nil)
}

func AuthLdapUser(ctx echo.Context) error {
	req := new(entity.AuthUserParams)

	if err := ctx.Bind(req); err != nil {
		logs.GetLogger().Errorf("api AuthLdapUser req is failed reqParams is %s  err is %s\n", req, err.Error())
		return ErrorResp(ctx, consts.StatusText[consts.CodeLdapParamIsError], consts.CodeLdapParamIsError)
	}

	//根据参数 查找LDAP 用户是否存在
	result, err := db.GetLdap().SearchUser(req)
	if err != nil {
		logs.GetLogger().Errorf("api AuthLdapUser SearchUser ldap  failed  is  err is %s\n", err.Error())
		return ErrorResp(ctx, consts.StatusText[consts.CodeLdapSearchUserFailed], consts.CodeLdapSearchUserFailed)
	}

	//result 为空查询不到ldap用户 返回登录失败
	if len(result.Entries) == 0 {
		logs.GetLogger().Errorf("api AuthLdapUser SearchUser ldap  can't find user    err is %s\n", err.Error())
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
	}

	fmt.Printf("cn is : %v\n", ldapUserInfo)

	user, err := service.QueryUser(ldapUserInfo.Cn)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			if err != nil {
				logs.GetLogger().Errorf("api AuthLdapUser QueryUser user ErrRecordNotFound   err is %s\n", err.Error())
				return ErrorResp(ctx, consts.StatusText[consts.CodeInternalServerError], consts.CodeInternalServerError)
			}

		}
		logs.GetLogger().Errorf("api AuthLdapUser QueryUser is failed   err is %s\n", err.Error())
		return ErrorResp(ctx, consts.StatusText[consts.CodeInternalServerError], consts.CodeInternalServerError)
	}

	if ldapUserInfo.Mail != user.Email {
		updates := map[string]interface{}{
			"email": ldapUserInfo.Mail,
		}

		user.Email = ldapUserInfo.Mail

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
