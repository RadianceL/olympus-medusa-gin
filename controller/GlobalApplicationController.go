package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"olympus-medusa/common"
	"olympus-medusa/common/language"
	"olympus-medusa/config"
	"olympus-medusa/controller/basic"
	Data "olympus-medusa/data/request"
	"olympus-medusa/data/response"
	"olympus-medusa/model"
)

type IGlobalApplicationController interface {
	basic.RestController
	CreateApplication(context *gin.Context)
	ListApplication(context *gin.Context)
	ListSupportLanguages(context *gin.Context)
}

type RestHandler struct {
	DB               *gorm.DB
	logger           *logrus.Logger
	applicationModel model.IApplicationModel
}

func NewGlobalApplicationController() IGlobalApplicationController {
	return RestHandler{
		DB:               common.GetDB(),
		logger:           config.GetLogger(),
		applicationModel: model.NewApplicationModel(),
	}
}

// CreateApplication 创建多语言应用/**
func (restHandler RestHandler) CreateApplication(context *gin.Context) {
	applicationAddRequest := &Data.ApplicationRequest{}
	err := context.ShouldBindBodyWith(&applicationAddRequest, binding.JSON)
	if applicationAddRequest.ApplicationName == "" {
		response.ResFail(context, "应用名称不能为空")
		return
	}
	if err != nil {
		response.ResErrCli(context, err)
		return
	}
	_, err = restHandler.applicationModel.AddApplication(applicationAddRequest)
	if err != nil {
		response.ResErrCli(context, err)
		return
	}
	response.ResSuccessMsg(context)
}

// ListApplication 查询应用列表/**
func (restHandler RestHandler) ListApplication(context *gin.Context) {
	applicationAddRequest := &Data.ApplicationRequest{}
	shouldBindBodyWithErr := context.ShouldBindBodyWith(&applicationAddRequest, binding.JSON)
	if shouldBindBodyWithErr != nil {
		response.ResFail(context, "json解析异常")
		return
	}
	searchApplicationList, searchApplicationError := restHandler.applicationModel.SearchApplicationList(applicationAddRequest)
	if searchApplicationError != nil {
		restHandler.logger.Error(errors.New("文件类型不合法"), "上传错误", false)
		response.ResFail(context, "应用处理异常")
		return
	}
	response.ResSuccess(context, searchApplicationList)
}

// ListSupportLanguages 查询应用支持的语言列表/**
func (restHandler RestHandler) ListSupportLanguages(context *gin.Context) {
	response.ResSuccess(context, language.Values())
}

//
//// CreateGlobalizationCopyWritingNamespace 创建应用空间namespace/**
//func (restHandler RestHandler) CreateGlobalizationCopyWritingNamespace(context *gin.Context) {
//	namespaceRequest := &Data.NamespaceRequest{}
//	shouldBindBodyWithErr := context.ShouldBindBodyWith(&namespaceRequest, binding.JSON)
//	if shouldBindBodyWithErr != nil {
//		response.ResFail(context, "json解析异常")
//		return
//	}
//	_, searchApplicationError := model.NamespaceHandler.CreateApplicationNamespace(namespaceRequest)
//	if searchApplicationError != nil {
//		restHandler.logger.Error(searchApplicationError)
//		response.ResErrCli(context, searchApplicationError)
//		return
//	}
//	response.ResSuccessMsg(context)
//}
//
//// ListGlobalizationCopyWritingStruct 获取多语言文案结构/**
//func (restHandler RestHandler) ListGlobalizationCopyWritingStruct(context *gin.Context) {
//	namespaceRequest := &Data.NamespaceRequest{}
//	shouldBindBodyWithErr := context.ShouldBindBodyWith(&namespaceRequest, binding.JSON)
//	if shouldBindBodyWithErr != nil {
//		response.ResFail(context, "json解析异常")
//		return
//	}
//	searchApplicationList, searchApplicationError := model.NamespaceHandler.ListApplicationNamespace(namespaceRequest)
//	if searchApplicationError != nil {
//		restHandler.logger.Error(searchApplicationError)
//		response.ResFail(context, "应用处理异常")
//		return
//	}
//	response.ResSuccess(context, searchApplicationList)
//}
//
//// ListGlobalizationCopyWritingNamespace 查询应用文案命名空间/**
//func (restHandler RestHandler) ListGlobalizationCopyWritingNamespace(context *gin.Context) {
//	globalDocumentRequest := &Data.GlobalDocumentRequest{}
//	shouldBindBodyWithErr := context.ShouldBindBodyWith(&globalDocumentRequest, binding.JSON)
//	if shouldBindBodyWithErr != nil {
//		response.ResFail(context, "json解析异常")
//		return
//	}
//	searchApplicationList, searchApplicationError :=
//		model.DocumentHandler.SearchDocumentByNamespaceId(globalDocumentRequest)
//	if searchApplicationError != nil {
//		restHandler.logger.Error(searchApplicationError)
//		response.ResFail(context, "应用处理异常")
//		return
//	}
//	response.ResSuccess(context, searchApplicationList)
//}
//
//// ImportGlobalizationCopyWriting 创建多语言文案/**
//func (restHandler RestHandler) ImportGlobalizationCopyWriting(context *gin.Context) {
//	file, _, err := context.Request.FormFile("excel")
//	if err != nil {
//		response.ResFail(context, "json解析异常")
//		return
//	}
//	resultData, err := model.DocumentHandler.ImportDocument(file)
//	if err != nil {
//		response.ResErrCliData(context, err, resultData)
//		return
//	}
//	response.ResSuccess(context, resultData)
//}
//
//func (restHandler RestHandler) ExportGlobalizationCopyWriting(context *gin.Context) {
//	globalDocumentRequest := &Data.GlobalDocumentRequest{}
//	shouldBindBodyWithErr := context.ShouldBindBodyWith(&globalDocumentRequest, binding.JSON)
//	if shouldBindBodyWithErr != nil {
//		response.ResFail(context, "json解析异常")
//		return
//	}
//	searchApplicationList, searchApplicationError := model.DocumentHandler.QueryDocument(globalDocumentRequest)
//	if searchApplicationError != nil {
//		restHandler.logger.Error(searchApplicationError)
//		response.ResFail(context, "应用处理异常")
//		return
//	}
//
//	fileName := "多语言数据导出列表.xlsx"
//	context.Writer.Header().Add("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
//	context.Writer.Header().Add("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
//	//创建文件
//	f := excelize.NewFile()
//	//创建excel表头信息
//	f = restHandler.ExcelHeader(f)
//	for i, document := range searchApplicationList {
//		curIndex := i + 2
//		var documentResult Data.TableGlobalDocumentExcel
//		_ = mapstructure.Decode(document, &documentResult)
//		strIndex := strconv.Itoa(curIndex)
//		f.SetCellValue("Sheet1", "A"+strIndex, document.ApplicationName+"("+strconv.Itoa(document.ApplicationId)+")")
//		f.SetCellValue("Sheet1", "B"+strIndex, document.NamespaceName+"("+strconv.Itoa(document.NamespaceId)+")")
//		f.SetCellValue("Sheet1", "C"+strIndex, document.DocumentCode)
//		f.SetCellValue("Sheet1", "D"+strIndex, document.DocumentDesc)
//		f.SetCellValue("Sheet1", "E"+strIndex, document.CountryName+"("+document.CountryIso+")")
//		f.SetCellValue("Sheet1", "F"+strIndex, document.DocumentValue)
//	}
//
//	var b bytes.Buffer
//	f.Write(&b)
//	content := bytes.NewReader(b.Bytes())
//	http.ServeContent(context.Writer, context.Request, fileName, time.Now(), content)
//}
//
//func (restHandler RestHandler) ExcelHeader(f *excelize.File) *excelize.File {
//	f.SetCellValue("Sheet1", "A1", "所属应用")
//	f.SetCellValue("Sheet1", "B1", "命名空间")
//	f.SetCellValue("Sheet1", "C1", "字段CODE")
//	f.SetCellValue("Sheet1", "D1", "描述")
//	f.SetCellValue("Sheet1", "E1", "语言")
//	f.SetCellValue("Sheet1", "F1", "翻译")
//	return f
//}
//
//// CreateGlobalizationCopyWriting 创建多语言文案/**
//func (restHandler RestHandler) CreateGlobalizationCopyWriting(context *gin.Context) {
//	json := &Data.GlobalDocumentRequest{}
//	err := context.ShouldBindBodyWith(&json, binding.JSON)
//	if err != nil {
//		response.ResFail(context, "json解析异常")
//		return
//	}
//	resultId, err := model.DocumentHandler.CreateDocument(json)
//	if err != nil {
//		response.ResErrCli(context, err)
//		return
//	}
//	println(resultId)
//	response.ResSuccessMsg(context)
//}
//
//// QueryGlobalizationCopyWritingDetail 创建多语言文案/**
//func (restHandler RestHandler) QueryGlobalizationCopyWritingDetail(context *gin.Context) {
//	globalDocumentRequest := &Data.GlobalDocumentRequest{}
//	err := context.ShouldBindBodyWith(&globalDocumentRequest, binding.JSON)
//	if err != nil {
//		response.ResFail(context, "json解析异常")
//		return
//	}
//	result, err := model.DocumentHandler.SearchDocumentById(globalDocumentRequest)
//	if err != nil {
//		response.ResErrCli(context, err)
//		return
//	}
//	response.ResSuccess(context, result)
//}
//
//// UpdateGlobalizationCopyWriting 更新多语言文案/**
//func (restHandler RestHandler) UpdateGlobalizationCopyWriting(context *gin.Context) {
//	json := &Data.GlobalDocumentRequest{}
//	err := context.ShouldBindBodyWith(&json, binding.JSON)
//	if err != nil {
//		response.ResFail(context, "json解析异常")
//		return
//	}
//	documentByDocumentUpdateResult, err := model.DocumentHandler.UpdateDocumentByDocumentId(json)
//	if err != nil {
//		response.ResErrCli(context, err)
//		return
//	}
//	if documentByDocumentUpdateResult == 0 {
//		response.ResFail(context, "更新文档失败，请稍后重试")
//		return
//	}
//	response.ResSuccessMsg(context)
//}
//
//func (restHandler RestHandler) DeleteGlobalizationCopyWriting(context *gin.Context) {
//	json := &Data.GlobalDocumentRequest{}
//	err := context.ShouldBindBodyWith(&json, binding.JSON)
//	if err != nil {
//		response.ResFail(context, "json解析异常")
//		return
//	}
//	documentByDocumentUpdateResult, err := model.DocumentHandler.DeleteDocumentByDocumentId(json)
//	if err != nil {
//		response.ResErrCli(context, err)
//		return
//	}
//	if documentByDocumentUpdateResult == 0 {
//		response.ResFail(context, "删除文档失败，请稍后重试")
//		return
//	}
//	response.ResSuccessMsg(context)
//}

// ListGlobalizationCopyWritingHistory 查询多语言文案历史/**
func (restHandler RestHandler) ListGlobalizationCopyWritingHistory(context *gin.Context) {
	json := &Data.GlobalDocumentRequest{}
	err := context.ShouldBindBodyWith(&json, binding.JSON)
	if err != nil {
		response.ResFail(context, "json解析异常")
		return
	}
	response.ResSuccessMsg(context)
}

// CommitGlobalizationCopyWriting 提交多语言文案更新/**
func (restHandler RestHandler) CommitGlobalizationCopyWriting(context *gin.Context) {
	json := &Data.GlobalDocumentRequest{}
	err := context.ShouldBindBodyWith(&json, binding.JSON)
	if err != nil {
		response.ResFail(context, "json解析异常")
		return
	}
	response.ResSuccessMsg(context)
}
