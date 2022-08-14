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

	_contRoles "app/controllers/roles"
	_repoRoles "app/repository/roles"
	_useRoles "app/usecase/roles"

	_contSkumanagement "app/controllers/sku_management"
	_repoSkumanagement "app/repository/sku_management"
	_useSkumanagement "app/usecase/sku_management"

	_repoRoleOutlet "app/repository/group_outlet"
	_useRoleOutlet "app/usecase/group_outlet"

	_repoUserRoles "app/repository/user_role"
	_useUserRoles "app/usecase/user_role"

	_contOutlets "app/controllers/outlets"
	_repoOutlets "app/repository/outlets"
	_useOutlets "app/usecase/outlets"

	_contInventory "app/controllers/inventory"
	_repoInventory "app/repository/inventory"
	_useInventory "app/usecase/inventory"

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

	repoRoles := _repoRoles.NewRepoRoles(postgres.Conn)
	useRoles := _useRoles.NewRoles(repoRoles, timeoutContext)
	_contRoles.NewContRole(g.G, useRoles)

	repoUserRole := _repoUserRoles.NewRepoUserRole(postgres.Conn)
	_ = _useUserRoles.NewUseUserRole(repoUserRole, timeoutContext)

	repoRoleOutlet := _repoRoleOutlet.NewRepoRoleOutlet(postgres.Conn)
	useRoleOutlet := _useRoleOutlet.NewUseRoleOutlet(repoRoleOutlet, timeoutContext)

	repoUser := _repoUser.NewRepoSysUser(postgres.Conn)
	useAuth := _useAuth.NewUserAuth(repoUser, repoFileUpload, repoUserSession, repoUserRole, timeoutContext)
	useUser := _useUser.NewUserSysUser(repoUser, repoUserRole, useRoleOutlet, timeoutContext)

	_contUser.NewContUsers(g.G, useUser)
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
	useOutlet := _useOutlets.NewUseOutlets(repoOutlet, repoOutletDetail, repoRoleOutlet, timeoutContext)
	_contOutlets.NewContOutlets(g.G, useOutlet)

	repoInventory := _repoInventory.NewRepoInventory(postgres.Conn)
	useInventory := _useInventory.NewUseInventory(repoInventory, timeoutContext)
	_contInventory.NewContInventory(g.G, useInventory)
}
