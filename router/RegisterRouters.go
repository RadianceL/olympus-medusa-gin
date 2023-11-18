package routers

import (
	"github.com/gin-gonic/gin"
	"olympus-medusa/config"
	"olympus-medusa/controller"
)

func RegisterRouterSys(app *gin.RouterGroup) {
	app.Use(config.LoggerToFile())
	applicationController := controller.NewGlobalApplicationController()
	// 创建应用
	app.POST("/application/create", applicationController.CreateApplication)
	// 查询应用列表
	app.POST("/application/list", applicationController.ListApplication)
	// 支持的语言列表
	app.POST("/support/language/list", applicationController.ListSupportLanguages)

	namespaceController := controller.NewNamespaceControllerController()
	// 创建应用空间namespace
	app.POST("/application/namespace/create", namespaceController.CreateGlobalizationCopyWritingNamespace)
	// 查询应用文案结构（返回该应用空间下的namespace -tree）
	app.POST("/application/namespace/list", namespaceController.ListGlobalizationCopyWritingStruct)
	// 查询应用文案命名空间
	app.POST("/application/namespace/document/list", namespaceController.ListGlobalizationCopyWritingNamespace)

	documentController := controller.NewDocumentController()
	// 导入应用文案
	app.POST("/application/document/import", documentController.ImportGlobalizationCopyWriting)
	// 导出应用文案
	app.POST("/application/document/export", documentController.ExportGlobalizationCopyWriting)
	// 创建应用文案
	app.POST("/application/document/create", documentController.CreateGlobalizationCopyWriting)
	// 应用文案详情
	app.POST("/application/document/query", documentController.QueryGlobalizationCopyWritingDetail)
	// 更新应用文案
	app.POST("/application/document/update", documentController.UpdateGlobalizationCopyWriting)
	// 删除应用文案
	app.POST("/application/document/delete", documentController.DeleteGlobalizationCopyWriting)
	// 审核应用文案
	app.POST("/application/document/commit", documentController.CommitGlobalizationCopyWriting)
	// 查询应用文案历史记录
	app.POST("/application/document/query/history", documentController.ListGlobalizationCopyWritingHistory)
}
