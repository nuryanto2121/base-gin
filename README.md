# Go Gin Auth base 

### Required
- Docker
    - [Install for Windows](https://docs.docker.com/desktop/windows/install/)
    - [Install for Mac](https://docs.docker.com/desktop/mac/install/)

## How to run
```
$ mkdir -p proto-gen/auth/
$ protoc --proto_path=proto/auth --go_out=proto-gen/auth --go_opt=paths=source_relative --go-grpc_out=proto-gen/auth --go-grpc_opt=paths=source_relative proto/auth/*.proto

$ go mod tidy

#run docker-compose up for DB 
$ docker-compose up -d

#
```
### Conf

You should modify `config.ini`

```
[database]
Type = postgres
User = postgres
Password = postgres_dev
Host = 127.0.0.1
Port = 54320
Name = postgres
TablePrefix = 

[redis]
Host = 127.0.0.1
Port = 6379
Password = redis_dev
Db = 0
...
```

### Run
```
$ go run main.go 
```

Project information and existing API

```
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /swagger/*any             --> github.com/swaggo/gin-swagger.CustomWrapHandler.func1 (6 handlers)
[GIN-debug] GET    /health_check             --> gitlab.com/369-engineer/369backend/account/controllers/auth.(*ContAuth).Health-fm (6 handlers)
[GIN-debug] POST   /api/auth/login           --> gitlab.com/369-engineer/369backend/account/controllers/auth.(*ContAuth).Login-fm (6 handlers)
[GIN-debug] POST   /api/auth/forgot          --> gitlab.com/369-engineer/369backend/account/controllers/auth.(*ContAuth).ForgotPassword-fm (6 handlers)
[GIN-debug] POST   /api/auth/change_password --> gitlab.com/369-engineer/369backend/account/controllers/auth.(*ContAuth).ChangePassword-fm (6 handlers)
[GIN-debug] POST   /api/auth/register        --> gitlab.com/369-engineer/369backend/account/controllers/auth.(*ContAuth).Register-fm (6 handlers)
[GIN-debug] POST   /api/auth/verify_otp      --> gitlab.com/369-engineer/369backend/account/controllers/auth.(*ContAuth).VerifyOTP-fm (6 handlers)
[GIN-debug] POST   /api/auth/logout          --> gitlab.com/369-engineer/369backend/account/controllers/auth.(*ContAuth).Logout-fm (6 handlers)
[GIN-debug] POST   /fileupload               --> gitlab.com/369-engineer/369backend/account/controllers/fileupload.(*ContFileUpload).CreateImage-fm (6 handlers)


Listening port is 8000
Actual pid is 4393
```
Swagger doc

![image](https://i.ibb.co/DVGZ5rW/swagger.png)

## Features

- RESTful API
- Gorm
- Swagger
- logging
- Jwt-go
- Gin
- Redis
- Email Smtp