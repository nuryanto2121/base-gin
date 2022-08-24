package routers

import (
	"time"

	_ "app/docs"
	"app/pkg/db"
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

	_repoRoleOutlet "app/repository/role_outlet"
	_useRoleOutlet "app/usecase/role_outlet"

	_repoUserRoles "app/repository/user_role"
	_useUserRoles "app/usecase/user_role"

	_contOutlets "app/controllers/outlets"
	_repoOutlets "app/repository/outlets"
	_useOutlets "app/usecase/outlets"

	_contInventory "app/controllers/inventory"
	_repoInventory "app/repository/inventory"
	_useInventory "app/usecase/inventory"

	_contOrder "app/controllers/order"
	_repoOrder "app/repository/order"
	_useOrder "app/usecase/order"

	_repoOutletDetail "app/repository/outlet_detail"

	_contTermAndConditional "app/controllers/term_and_conditional"
	_repoTermAndConditional "app/repository/term_and_conditional"
	_useTermAndConditional "app/usecase/term_and_conditional"

	_contAuditLogs "app/controllers/audit_logs"
	_repoAuditLogs "app/repository/audit_logs"
	_useAuditLogs "app/usecase/audit_logs"
)

type GinRoutes struct {
	G *gin.Engine
}

func (g *GinRoutes) Init() {
	timeoutContext := time.Duration(setting.ServerSetting.ReadTimeout) * time.Second
	dbConn := db.NewDBdelegate(setting.DatabaseSetting.Debug)
	dbConn.Init()

	r := g.G
	r.GET("/v1//swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	repoUserSession := _repoUserSession.NewRepoUserSession(dbConn)

	repoFileUpload := _repoFileUpload.NewRepoFileUpload(dbConn)
	useFileUpload := _useFileUpload.NewSaFileUpload(repoFileUpload, timeoutContext)

	repoRoles := _repoRoles.NewRepoRoles(dbConn)
	useRoles := _useRoles.NewRoles(repoRoles, timeoutContext)
	_contRoles.NewContRole(g.G, useRoles)

	repoUserRole := _repoUserRoles.NewRepoUserRole(dbConn)
	_ = _useUserRoles.NewUseUserRole(repoUserRole, timeoutContext)

	repoRoleOutlet := _repoRoleOutlet.NewRepoRoleOutlet(dbConn)
	_ = _useRoleOutlet.NewUseRoleOutlet(repoRoleOutlet, timeoutContext)

	repoUser := _repoUser.NewRepoSysUser(dbConn)
	useAuth := _useAuth.NewUserAuth(repoUser, repoFileUpload, repoUserSession, repoUserRole, repoRoleOutlet, timeoutContext)
	useUser := _useUser.NewUserSysUser(repoUser, repoUserRole, repoRoleOutlet, timeoutContext)

	_contUser.NewContUsers(g.G, useUser)
	_contAuth.NewContAuth(g.G, useAuth)

	_contFileUpload.NewContFileUpload(g.G, useFileUpload)

	repoHolidays := _repoHolidays.NewRepoHolidays(dbConn)
	userHolidays := _useHolidays.NewHolidaysHolidays(repoHolidays, timeoutContext)
	_contHolidays.NewContHolidays(g.G, userHolidays)

	repoSkuManagement := _repoSkumanagement.NewRepoSkuManagement(dbConn)
	useSkuManagement := _useSkumanagement.NewSkuManagement(repoSkuManagement, timeoutContext)
	_contSkumanagement.NewContSkuManagement(g.G, useSkuManagement)
	repoOutlet := _repoOutlets.NewRepoOutlets(dbConn)
	repoOutletDetail := _repoOutletDetail.NewRepoOutletDetail(dbConn)
	useOutlet := _useOutlets.NewUseOutlets(repoOutlet, repoOutletDetail, repoRoleOutlet, timeoutContext)
	_contOutlets.NewContOutlets(g.G, useOutlet)

	repoInventory := _repoInventory.NewRepoInventory(dbConn)
	useInventory := _useInventory.NewUseInventory(repoInventory, timeoutContext)
	_contInventory.NewContInventory(g.G, useInventory)

	repoTermAndConditional := _repoTermAndConditional.NewRepoTermAndConditioinal(dbConn)
	useTermAndConditional := _useTermAndConditional.NewTermAndConditional(repoTermAndConditional, timeoutContext)
	_contTermAndConditional.NewContTermAndConditional(g.G, useTermAndConditional)

	repoOrder := _repoOrder.NewRepoOrder(dbConn)
	useOrder := _useOrder.NewUseOrder(repoOrder, repoOutlet, repoSkuManagement, timeoutContext)
	_contOrder.NewContOrder(g.G, useOrder)

	repoAuditLogs := _repoAuditLogs.NewRepoAuditLogs(dbConn)
	useAuditLogs := _useAuditLogs.NewUseAuditLogs(repoAuditLogs, time.Nanosecond)
	_contAuditLogs.NewContAuditLogs(g.G, useAuditLogs)
}
