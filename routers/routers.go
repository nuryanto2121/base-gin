package routers

import (
	"time"

	_ "app/docs"
	postgres "app/pkg/postgres"
	"app/pkg/setting"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"

	swaggerFiles "github.com/swaggo/files"

	_repoUserSession "app/repository/user_session"

	_contAuth "app/controllers/auth"

	_contUser "app/controllers/users"
	_repoUser "app/repository/user"
	_useAuth "app/usecase/auth"
	_useUser "app/usecase/user"

	_contFileUpload "app/controllers/fileupload"
	_repoFileUpload "app/repository/fileupload"
	_useFileUpload "app/usecase/fileupload"

	_contHolidays "app/controllers/holidays"
	_repoHolidays "app/repository/holidays"
	_useHolidays "app/usecase/holidays"

	_contGroups "app/controllers/group"
	_repoGroups "app/repository/group"
	_useGroups "app/usecase/group"

	_contSkumanagement "app/controllers/sku_management"
	_repoGroupOutlet "app/repository/group_outlet"
	_repoSkumanagement "app/repository/sku_management"
	_useGroupOutlet "app/usecase/group_outlet"
	_useSkumanagement "app/usecase/sku_management"

	_repoUserGroups "app/repository/user_group"
	_useUserGroups "app/usecase/user_group"

	_contOutlets "app/controllers/outlets"
	_repoOutlets "app/repository/outlets"
	_useOutlets "app/usecase/outlets"

	_repoOutletDetail "app/repository/outlet_detail"
)

type GinRoutes struct {
	G *gin.Engine
}

func (g *GinRoutes) Init() {
	timeoutContext := time.Duration(setting.ServerSetting.ReadTimeout) * time.Second

	r := g.G
	r.GET("/v1//swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	repoUserSession := _repoUserSession.NewRepoUserSession(postgres.Conn)

	repoFileUpload := _repoFileUpload.NewRepoFileUpload(postgres.Conn)
	useFileUpload := _useFileUpload.NewSaFileUpload(repoFileUpload, timeoutContext)

	repoGroups := _repoGroups.NewRepoGroups(postgres.Conn)
	useGroups := _useGroups.NewGroups(repoGroups, timeoutContext)
	_contGroups.NewContGroup(g.G, useGroups)

	repoUserGroup := _repoUserGroups.NewRepoUserGroup(postgres.Conn)
	_ = _useUserGroups.NewUseUserGroup(repoUserGroup, timeoutContext)

	repoGroupOutlet := _repoGroupOutlet.NewRepoGroupOutlet(postgres.Conn)
	useGroupOutlet := _useGroupOutlet.NewUseGroupOutlet(repoGroupOutlet, timeoutContext)

	repoUser := _repoUser.NewRepoSysUser(postgres.Conn)
	useAuth := _useAuth.NewUserAuth(repoUser, repoFileUpload, repoUserSession, repoUserGroup, timeoutContext)
	useUser := _useUser.NewUserSysUser(repoUser, repoUserGroup, useGroupOutlet, timeoutContext)

	_contUser.NewContGroup(g.G, useUser)
	_contAuth.NewContAuth(g.G, useAuth)

	_contFileUpload.NewContFileUpload(g.G, useFileUpload)

	repoHolidays := _repoHolidays.NewRepoHolidays(postgres.Conn)
	userHolidays := _useHolidays.NewHolidaysHolidays(repoHolidays, timeoutContext)
	_contHolidays.NewContHolidays(g.G, userHolidays)

	reposkumanagement := _repoSkumanagement.NewRepoSkuManagement(postgres.Conn)
	useskumanagement := _useSkumanagement.NewSkuManagement(reposkumanagement, timeoutContext)
	_contSkumanagement.NewContSkuManagement(g.G, useskumanagement)
	repoOutlet := _repoOutlets.NewRepoOutlets(postgres.Conn)
	repoOutletDetail := _repoOutletDetail.NewRepoOutletDetail(postgres.Conn)
	useOutlet := _useOutlets.NewUseOutlets(repoOutlet, repoOutletDetail, timeoutContext)
	_contOutlets.NewContOutlets(g.G, useOutlet)
}
