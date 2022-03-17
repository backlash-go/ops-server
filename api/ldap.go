package api

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"log"
	"ops-server/consts"
	"ops-server/db"
	"ops-server/entity"
	"ops-server/service"
	"ops-server/utils"
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

	//result 为空查询不到用户
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
		}
		utils.GetLogger().Errorf("api AuthLdapUser QueryUser is failed   err is %s\n", err.Error())
		return ErrorResp(ctx, consts.StatusText[consts.CodeInternalServerError], consts.CodeInternalServerError)
	}

	if ldapUserInfo.Mail != user.Email {
		updates := map[string]interface{}{
			"email": ldapUserInfo.Mail,
		}

		if err := service.UpdateUser(user.Id, updates); err != nil {
			utils.GetLogger().Errorf("api AuthLdapUser UpdateUser is failed   err is %s\n", err.Error())
			return ErrorResp(ctx, consts.StatusText[consts.CodeInternalServerError], consts.CodeInternalServerError)
		}
	}

	fmt.Printf("cn is : %v\n", user)

	return SuccessResp(ctx, nil)

}
