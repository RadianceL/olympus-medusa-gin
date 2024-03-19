package controller

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
	"net/http"
	"olympus-medusa/config"
	"olympus-medusa/data/data"
	Request "olympus-medusa/data/request"
	"olympus-medusa/data/response"
	"olympus-medusa/model"
	"strconv"
	"time"
)

type IDocumentController interface {
	ImportGlobalizationCopyWriting(context *gin.Context)

	ExportGlobalizationCopyWriting(context *gin.Context)

	CreateGlobalizationCopyWriting(context *gin.Context)

	QueryGlobalizationCopyWritingDetail(context *gin.Context)

	UpdateGlobalizationCopyWriting(context *gin.Context)

	DeleteGlobalizationCopyWriting(context *gin.Context)

	ListGlobalizationCopyWritingHistory(context *gin.Context)

	CommitGlobalizationCopyWriting(context *gin.Context)
}

type DocumentController struct {
	logger        *logrus.Logger
	documentModel model.IDocumentModel
}

func NewDocumentController() IDocumentController {
	return DocumentController{
		logger:        config.GetLogger(),
		documentModel: model.NewDocumentModel(),
	}
}

// ImportGlobalizationCopyWriting 创建多语言文案/**
func (controller DocumentController) ImportGlobalizationCopyWriting(context *gin.Context) {
	_, _, err := context.Request.FormFile("excel")
	if err != nil {
		response.ResFail(context, "json解析异常")
		return
	}
	if err != nil {
		response.ResErrCliData(context, err, nil)
		return
	}
	response.ResSuccess(context, nil)
}

func (controller DocumentController) ExportGlobalizationCopyWriting(context *gin.Context) {
	globalDocumentRequest := &Request.GlobalDocumentRequest{}
	shouldBindBodyWithErr := context.ShouldBindBodyWith(&globalDocumentRequest, binding.JSON)
	if shouldBindBodyWithErr != nil {
		response.ResFail(context, "json解析异常")
		return
	}
	searchApplicationList, searchApplicationError := controller.documentModel.QueryDocument(globalDocumentRequest)
	if searchApplicationError != nil {
		controller.logger.Error(searchApplicationError)
		response.ResFail(context, "应用处理异常")
		return
	}

	fileName := "多语言数据导出列表.xlsx"
	context.Writer.Header().Add("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
	context.Writer.Header().Add("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	//创建文件
	f := excelize.NewFile()
	//创建excel表头信息
	f = controller.ExcelHeader(f)
	for i, document := range searchApplicationList {
		curIndex := i + 2
		var documentResult data.GlobalDocumentExcel
		_ = mapstructure.Decode(document, &documentResult)
		strIndex := strconv.Itoa(curIndex)
		err := f.SetCellValue("Sheet1", "A"+strIndex, document.ApplicationName+"("+strconv.Itoa(document.ApplicationId)+")")
		if err != nil {
			return
		}
		err = f.SetCellValue("Sheet1", "B"+strIndex, document.NamespaceName+"("+strconv.Itoa(document.NamespaceId)+")")
		if err != nil {
			return
		}
		err = f.SetCellValue("Sheet1", "C"+strIndex, document.DocumentCode)
		if err != nil {
			return
		}
		err = f.SetCellValue("Sheet1", "D"+strIndex, document.DocumentDesc)
		if err != nil {
			return
		}
		err = f.SetCellValue("Sheet1", "E"+strIndex, document.CountryName+"("+document.CountryIso+")")
		if err != nil {
			return
		}
		err = f.SetCellValue("Sheet1", "F"+strIndex, document.DocumentValue)
		if err != nil {
			return
		}
	}
	var b bytes.Buffer
	err := f.Write(&b)
	if err != nil {
		return
	}
	content := bytes.NewReader(b.Bytes())
	http.ServeContent(context.Writer, context.Request, fileName, time.Now(), content)
}

// CreateGlobalizationCopyWriting 创建多语言文案/**
func (controller DocumentController) CreateGlobalizationCopyWriting(context *gin.Context) {
	json := &Request.GlobalDocumentRequest{}
	err := context.ShouldBindBodyWith(&json, binding.JSON)
	if err != nil {
		response.ResFail(context, "json解析异常")
		return
	}
	resultId, err := controller.documentModel.CreateDocument(json)
	if err != nil {
		response.ResErrCli(context, err)
		return
	}
	println(resultId)
	response.ResSuccessMsg(context)
}

// QueryGlobalizationCopyWritingDetail 创建多语言文案/**
func (controller DocumentController) QueryGlobalizationCopyWritingDetail(context *gin.Context) {
	globalDocumentRequest := &Request.GlobalDocumentRequest{}
	err := context.ShouldBindBodyWith(&globalDocumentRequest, binding.JSON)
	if err != nil {
		response.ResFail(context, "json解析异常")
		return
	}
	result, err := controller.documentModel.SearchDocumentById(globalDocumentRequest)
	if err != nil {
		response.ResErrCli(context, err)
		return
	}
	response.ResSuccess(context, result)
}

// UpdateGlobalizationCopyWriting 更新多语言文案/**
func (controller DocumentController) UpdateGlobalizationCopyWriting(context *gin.Context) {
	json := &Request.GlobalDocumentRequest{}
	err := context.ShouldBindBodyWith(&json, binding.JSON)
	if err != nil {
		response.ResFail(context, "json解析异常")
		return
	}
	documentByDocumentUpdateResult, err := controller.documentModel.UpdateDocumentByDocumentId(json)
	if err != nil {
		response.ResErrCli(context, err)
		return
	}
	if documentByDocumentUpdateResult == 0 {
		response.ResFail(context, "更新文档失败，请稍后重试")
		return
	}
	response.ResSuccessMsg(context)
}

func (controller DocumentController) DeleteGlobalizationCopyWriting(context *gin.Context) {
	json := &Request.GlobalDocumentRequest{}
	err := context.ShouldBindBodyWith(&json, binding.JSON)
	if err != nil {
		response.ResFail(context, "json解析异常")
		return
	}
	documentByDocumentUpdateResult, err := controller.documentModel.DeleteDocumentByDocumentId(json)
	if err != nil {
		response.ResErrCli(context, err)
		return
	}
	if documentByDocumentUpdateResult == 0 {
		response.ResFail(context, "删除文档失败，请稍后重试")
		return
	}
	response.ResSuccessMsg(context)
}

// ListGlobalizationCopyWritingHistory 查询多语言文案历史/**
func (controller DocumentController) ListGlobalizationCopyWritingHistory(context *gin.Context) {
	json := &Request.GlobalDocumentRequest{}
	err := context.ShouldBindBodyWith(&json, binding.JSON)
	if err != nil {
		response.ResFail(context, "json解析异常")
		return
	}
	response.ResSuccessMsg(context)
}

// CommitGlobalizationCopyWriting 提交多语言文案更新/**
func (controller DocumentController) CommitGlobalizationCopyWriting(context *gin.Context) {
	json := &Request.GlobalDocumentRequest{}
	err := context.ShouldBindBodyWith(&json, binding.JSON)
	if err != nil {
		response.ResFail(context, "json解析异常")
		return
	}
	response.ResSuccessMsg(context)
}

func (controller DocumentController) ExcelHeader(f *excelize.File) *excelize.File {
	err := f.SetCellValue("Sheet1", "A1", "所属应用")
	if err != nil {
		return nil
	}
	err = f.SetCellValue("Sheet1", "B1", "命名空间")
	if err != nil {
		return nil
	}
	err = f.SetCellValue("Sheet1", "C1", "字段CODE")
	if err != nil {
		return nil
	}
	err = f.SetCellValue("Sheet1", "D1", "描述")
	if err != nil {
		return nil
	}
	err = f.SetCellValue("Sheet1", "E1", "语言")
	if err != nil {
		return nil
	}
	err = f.SetCellValue("Sheet1", "F1", "翻译")
	if err != nil {
		return nil
	}
	return f
}
