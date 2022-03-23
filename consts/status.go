package consts

const (
	CodeSuccess              = 20000
	CodeLdapCreateUserFailed = 30001
	CodeLdapDeleteUserFailed = 30002
	CodeLdapParamIsError     = 30003
	CodeLdapConnectionFailed = 30004
	CodeLdapUserNotExist     = 30005
	CodeLdapSearchUserFailed = 30006
	CodeUserNoAssignRole     = 30007

	CodeUserNoApiPermission  = 40001

	CodeNeedLogin = 400003

	CodeInternalServerError = 50000
)

var StatusText = map[int]string{
	CodeLdapCreateUserFailed: "创建LDAP用户失败",
	CodeLdapDeleteUserFailed: "删除LDAP用户失败",
	CodeLdapParamIsError:     "Ldap请求参数错误",
	CodeLdapConnectionFailed: "Ldap服务器连接失败",
	CodeLdapUserNotExist:     "Ldap用户不存在",
	CodeLdapSearchUserFailed: "搜索Ldap用户失败",
	CodeInternalServerError:  "服务器内部错误，请联系开发人员",
	CodeNeedLogin:            "invalid authorization",
	CodeUserNoAssignRole:     "用户没有分配角色",
	CodeUserNoApiPermission:  "接口没有权限",
}
