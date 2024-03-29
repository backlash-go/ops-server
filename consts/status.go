package consts

const (
	CodeSuccess = 20000

	CodeLdapCreateUserFailed = 30001
	CodeLdapDeleteUserFailed = 30002
	CodeParamIsError         = 30003
	CodeLdapConnectionFailed = 30004
	CodeLdapUserNotExist     = 30005
	CodeLdapSearchUserFailed = 30006

	CodeUserNoAssignRole = 30007

	CodeUserDelTokenFailed = 30008

	CodeUserNoApiPermission = 30009

	CodeUserPasswordModifyFailed = 30010

	CodeQueryUserNotExist = 30011

	CodeUserInfoModifyFailed = 30012

	CodeNeedLogin = 40003

	CodeInternalServerError = 50000
)

var StatusText = map[int]string{
	CodeLdapCreateUserFailed:     "创建LDAP用户失败",
	CodeLdapDeleteUserFailed:     "删除LDAP用户失败",
	CodeParamIsError:             "Ldap请求参数错误",
	CodeLdapConnectionFailed:     "Ldap服务器连接失败",
	CodeLdapUserNotExist:         "Ldap用户不存在或者密码错误",
	CodeLdapSearchUserFailed:     "搜索Ldap用户失败",
	CodeInternalServerError:      "服务器内部错误，请联系开发人员",
	CodeNeedLogin:                "登录过期,请重新登录",
	CodeUserNoAssignRole:         "用户没有分配角色",
	CodeUserNoApiPermission:      "接口没有权限",
	CodeUserDelTokenFailed:       "删除TOKEN失败",
	CodeUserPasswordModifyFailed: "修改密码失败",
	CodeQueryUserNotExist:        "用户不存在",
	CodeUserInfoModifyFailed:     "修改用户信息失败",
}
