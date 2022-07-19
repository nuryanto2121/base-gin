package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"app/pkg/logging"
	"app/pkg/middleware"
	multilanguage "app/pkg/multiLanguage"
	"app/pkg/postgres"

	"app/pkg/redisdb"
	s3gateway "app/pkg/s3"
	"app/pkg/setting"
	"app/routers"
)

func init() {
	setting.Setup()
	multilanguage.Setup()
	logging.Setup()
	postgres.Setup()
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
