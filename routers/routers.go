package routers

import (
	"time"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	_ "gitlab.com/369-engineer/369backend/account/docs"
	maria "gitlab.com/369-engineer/369backend/account/pkg/mariadb"
	"gitlab.com/369-engineer/369backend/account/pkg/setting"

	_contAuth "gitlab.com/369-engineer/369backend/account/controllers/auth"

	_repoUser "gitlab.com/369-engineer/369backend/account/repository/user"
	_useAuth "gitlab.com/369-engineer/369backend/account/usecase/auth"
	_useUser "gitlab.com/369-engineer/369backend/account/usecase/user"

	_contFileUpload "gitlab.com/369-engineer/369backend/account/controllers/fileupload"
	_repoFileUpload "gitlab.com/369-engineer/369backend/account/repository/fileupload"
	_useFileUpload "gitlab.com/369-engineer/369backend/account/usecase/fileupload"
)

type GinRoutes struct {
	G *gin.Engine
}

func (g *GinRoutes) Init() {
	timeoutContext := time.Duration(setting.ServerSetting.ReadTimeout) * time.Second

	r := g.G
	r.GET("/v1/account//swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	repoFileUpload := _repoFileUpload.NewRepoFileUpload(maria.Conn)
	useFileUpload := _useFileUpload.NewSaFileUpload(repoFileUpload, timeoutContext)

	repoUser := _repoUser.NewRepoSysUser(maria.Conn)
	useAuth := _useAuth.NewUserAuth(repoUser, repoFileUpload, timeoutContext)
	_ = _useUser.NewUserSysUser(repoUser, timeoutContext)

	_contAuth.NewContAuth(g.G, useAuth)

	_contFileUpload.NewContFileUpload(g.G, useFileUpload)

}
