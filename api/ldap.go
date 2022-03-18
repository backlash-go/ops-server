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
	"ops-server/models"
	"ops-server/service"
	"ops-server/utils"
	"time"
)

func OperateLdap(g *echo.Group) {
	g.POST("/user/create", CreateLdapUser)
	g.DELETE("/user/delete", DeleteLdapUser)
	g.POST("/user/auth", AuthLdapUser)

}

func CreateLdapUser(ctx echo.Context) error {
	//create ldap connection
	//DefaultLdap, err := service.CreateLdapConnection()
	//if err != nil {
	//	utils.GetLogger().Errorf("CreateLdapConnection is failed err is %s\n", err.Error())
	//	return ErrorResp(ctx, consts.StatusText[consts.CodeLdapConnectionFailed], consts.CodeLdapConnectionFailed)
	//}

	// receive create user http request params
	req := &entity.CreateUserParams{}
	if err := ctx.Bind(req); err != nil {
		utils.GetLogger().Errorf("CreateLdapUser req is failed reqParams is %s  err is %s\n", req, err.Error())
		return ErrorResp(ctx, consts.StatusText[consts.CodeLdapParamIsError], consts.CodeLdapParamIsError)
	}

	log.Printf("req is  %s\n", req)

	if req.Cn == "" || req.DisplayName == "" || req.GivenName == "" || len(req.EmployeeType) == 0 || req.UserPassword == "" {
		utils.GetLogger().Errorf("req is failed reqParams is %s\n", req)
		return ErrorResp(ctx, consts.StatusText[consts.CodeLdapParamIsError], consts.CodeLdapParamIsError)
	}

	if err := db.GetLdap().CreateUser(req); err != nil {
		utils.GetLogger().Errorf("CreateUser is failed err is %s\n", err.Error())
		return ErrorResp(ctx, consts.StatusText[consts.CodeLdapCreateUserFailed], consts.CodeLdapCreateUserFailed)
	}

	return SuccessResp(ctx, nil)

}

func DeleteLdapUser(ctx echo.Context) error {
	//create ldap connection
	req := new(entity.DeleteUserParams)
	if err := ctx.Bind(req); err != nil {
		utils.GetLogger().Errorf("DeleteLdapUser req is failed reqParams is %s  err is %s\n", req, err.Error())
		return ErrorResp(ctx, consts.StatusText[consts.CodeLdapParamIsError], consts.CodeLdapParamIsError)
	}

	if req.Dn == "" {
		utils.GetLogger().Errorf("DeleteLdapUser req.Dn is null   req is %s\n", req)
		return ErrorResp(ctx, consts.StatusText[consts.CodeLdapParamIsError], consts.CodeLdapParamIsError)
	}

	if err := db.GetLdap().DeleteUser(req); err != nil {
		utils.GetLogger().Errorf("DeleteUser is failed err is %s\n", err.Error())
		return ErrorResp(ctx, consts.StatusText[consts.CodeLdapDeleteUserFailed], consts.CodeLdapDeleteUserFailed)
	}

	return SuccessResp(ctx, nil)
}

func AuthLdapUser(ctx echo.Context) error {

	req := new(entity.AuthUserParams)

	if err := ctx.Bind(req); err != nil {
		utils.GetLogger().Errorf("api AuthLdapUser req is failed reqParams is %s  err is %s\n", req, err.Error())
		return ErrorResp(ctx, consts.StatusText[consts.CodeLdapParamIsError], consts.CodeLdapParamIsError)
	}

	//根据参数 查找LDAP 用户是否存在
	result, err := db.GetLdap().SearchUser(req)
	if err != nil {
		utils.GetLogger().Errorf("api AuthLdapUser SearchUser ldap  failed  is  err is %s\n", err.Error())
		return ErrorResp(ctx, consts.StatusText[consts.CodeLdapSearchUserFailed], consts.CodeLdapSearchUserFailed)
	}

	//result 为空查询不到ldap用户 返回登录失败
	if len(result.Entries) == 0 {
		utils.GetLogger().Errorf("api AuthLdapUser SearchUser ldap  can't find user    err is %s\n", err.Error())
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
			//todo create user record

			userId, err := service.AddUser(models.User{UserName: ldapUserInfo.Cn, Email: ldapUserInfo.Mail})
			if err != nil {
				utils.GetLogger().Errorf("api AuthLdapUser AddUser is failed   err is %s\n", err.Error())
				return ErrorResp(ctx, consts.StatusText[consts.CodeInternalServerError], consts.CodeInternalServerError)
			}
			token := xid.New().String()

			tmpMap := make(map[string]interface{})
			tmpMap["id"] = userId
			tmpMap["user_name"] = ldapUserInfo.Cn
			tmpMap["email"] = ldapUserInfo.Mail

			if err := db.RedisHMSet(token, tmpMap); err != nil {
				utils.GetLogger().Errorf("api AuthLdapUser RedisHMSet is failed   err is %s\n", err.Error())
				return ErrorResp(ctx, consts.StatusText[consts.CodeInternalServerError], consts.CodeInternalServerError)
			}

			if err := db.RedisSetKeyTtl(token, time.Minute*60); err != nil {
				utils.GetLogger().Errorf("api AuthLdapUser RedisSetKeyTtl is failed   err is %s\n", err.Error())
				return ErrorResp(ctx, consts.StatusText[consts.CodeInternalServerError], consts.CodeInternalServerError)
			}

			userInfo := models.User{Id: userId, UserName: ldapUserInfo.Cn, Email: ldapUserInfo.Mail}

			return SuccessResp(ctx, userInfo)

		}
		utils.GetLogger().Errorf("api AuthLdapUser QueryUser is failed   err is %s\n", err.Error())
		return ErrorResp(ctx, consts.StatusText[consts.CodeInternalServerError], consts.CodeInternalServerError)
	}

	if ldapUserInfo.Mail != user.Email {
		updates := map[string]interface{}{
			"email": ldapUserInfo.Mail,
		}

		user.Email = ldapUserInfo.Mail

		if err := service.UpdateUser(user.Id, updates); err != nil {
			utils.GetLogger().Errorf("api AuthLdapUser UpdateUser is failed   err is %s\n", err.Error())
			return ErrorResp(ctx, consts.StatusText[consts.CodeInternalServerError], consts.CodeInternalServerError)
		}
	}

	//生成TOKEN
	token := xid.New().String()

	//存入redis

	tmpMap := make(map[string]interface{})
	tmpMap["id"] = user.Id
	tmpMap["user_name"] = user.UserName
	tmpMap["email"] = user.Email

	if err := db.RedisHMSet(token, tmpMap); err != nil {
		utils.GetLogger().Errorf("api AuthLdapUser RedisHMSet is failed   err is %s\n", err.Error())
		return ErrorResp(ctx, consts.StatusText[consts.CodeInternalServerError], consts.CodeInternalServerError)
	}

	if err := db.RedisSetKeyTtl(token, time.Minute*60); err != nil {
		utils.GetLogger().Errorf("api AuthLdapUser RedisSetKeyTtl is failed   err is %s\n", err.Error())
		return ErrorResp(ctx, consts.StatusText[consts.CodeInternalServerError], consts.CodeInternalServerError)
	}
	fmt.Printf("cn is : %v\n", user)



	resp := entity.LoginAuthResp{
		Token:    token,
		UserId: int64(user.Id),
		UserName: user.UserName,
		Email:    user.Email,
	}

	return SuccessResp(ctx, resp)

}
