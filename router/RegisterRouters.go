package routers

import (
	"github.com/gin-gonic/gin"
	"olympus-medusa/config"
	"olympus-medusa/controller"
)

func RegisterRouterSys(app *gin.RouterGroup) {
	app.Use(config.LoggerToFile())
	menu := controller.GlobalApplicationController{}
	globalApplicationController := controller.NewGlobalApplicationController()

	// 创建应用
	app.POST("/application/create", globalApplicationController.CreateApplication)
	// 查询应用列表
	app.POST("/application/list", globalApplicationController.ListApplication)
	// 支持的语言列表
	app.POST("/support/language/list", globalApplicationController.ListSupportLanguages)
	//// 创建应用空间namespace
	//app.POST("/application/namespace/create", menu.CreateGlobalizationCopyWritingNamespace)
	//// 查询应用文案结构（返回该应用空间下的namespace -tree）
	//app.POST("/application/namespace/list", menu.ListGlobalizationCopyWritingStruct)
	//// 查询应用文案命名空间
	//app.POST("/application/namespace/document/list", menu.ListGlobalizationCopyWritingNamespace)
	//// 导入应用文案
	//app.POST("/application/document/import", menu.ImportGlobalizationCopyWriting)
	//// 导出应用文案
	//app.POST("/application/document/export", menu.ExportGlobalizationCopyWriting)
	//// 创建应用文案
	//app.POST("/application/document/create", menu.CreateGlobalizationCopyWriting)
	//// 应用文案详情
	//app.POST("/application/document/query", menu.QueryGlobalizationCopyWritingDetail)
	//// 更新应用文案
	//app.POST("/application/document/update", menu.UpdateGlobalizationCopyWriting)
	//// 删除应用文案
	//app.POST("/application/document/delete", menu.DeleteGlobalizationCopyWriting)
	// 审核应用文案
	app.POST("/application/document/commit", menu.CommitGlobalizationCopyWriting)
	// 查询应用文案历史记录
	app.POST("/application/document/query/history", menu.ListGlobalizationCopyWritingHistory)
}
