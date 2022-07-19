package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"gitlab.com/369-engineer/369backend/account/pkg/logging"
	maria "gitlab.com/369-engineer/369backend/account/pkg/mariadb"
	"gitlab.com/369-engineer/369backend/account/pkg/middleware"
	multilanguage "gitlab.com/369-engineer/369backend/account/pkg/multiLanguage"

	"gitlab.com/369-engineer/369backend/account/pkg/redisdb"
	s3gateway "gitlab.com/369-engineer/369backend/account/pkg/s3"
	"gitlab.com/369-engineer/369backend/account/pkg/setting"
	"gitlab.com/369-engineer/369backend/account/routers"
)

func init() {
	setting.Setup()
	multilanguage.Setup()
	logging.Setup()
	maria.Setup()
	redisdb.Setup()
	s3gateway.Setup()
}

// @title Base 369
// @version 1.0
// @description Backend REST API for golang nuryanto2121

// @contact.name Nuryanto
// @contact.url https://www.linkedin.com/in/nuryanto-1b2721156/
// @contact.email nuryantofattih@gmail.com

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	gin.SetMode(setting.ServerSetting.RunMode)

	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)

	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	// r.Use(cors.Default())
	r.Use(middleware.CORSMiddleware())

	R := routers.GinRoutes{G: r}
	R.Init()
	r.Run(endPoint)
	fmt.Println("Server Run at ", setting.ServerSetting.HttpPort)

}
