# ops-server

> 这是一个基于LDAP登陆验证 以及管理LDAP用户的HTTP服务

# 项目环境

```cassandraql
go 1.14


```

```bash
# 克隆项目
git clone https://github.com/backlash-go/ops-server.git

# 进入项目目录
cd ops-server

# 安装依赖
go mod download 


# 启动服务
go run main.go
```



# build on mac  run in  linux  
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build