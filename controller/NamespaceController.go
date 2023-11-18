package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
	"olympus-medusa/config"
	Rquest "olympus-medusa/data/request"
	"olympus-medusa/data/response"
	"olympus-medusa/model"
)

type INamespaceController interface {
	CreateGlobalizationCopyWritingNamespace(context *gin.Context)

	ListGlobalizationCopyWritingStruct(context *gin.Context)

	ListGlobalizationCopyWritingNamespace(context *gin.Context)
}

type NamespaceControllerController struct {
	logger         *logrus.Logger
	namespaceModel model.INamespaceModel
}

func NewNamespaceControllerController() INamespaceController {
	return NamespaceControllerController{
		logger:         config.GetLogger(),
		namespaceModel: model.NewNamespaceModel(),
	}
}

// CreateGlobalizationCopyWritingNamespace 创建应用空间namespace/**
func (controller NamespaceControllerController) CreateGlobalizationCopyWritingNamespace(context *gin.Context) {
	namespaceRequest := &Rquest.NamespaceRequest{}
	shouldBindBodyWithErr := context.ShouldBindBodyWith(&namespaceRequest, binding.JSON)
	if shouldBindBodyWithErr != nil {
		response.ResFail(context, "json解析异常")
		return
	}
	namespace := namespaceRequest.ConvertToTableApplicationNamespace(namespaceRequest)
	_, searchApplicationError := controller.namespaceModel.CreateApplicationNamespace(namespace)
	if searchApplicationError != nil {
		controller.logger.Error(searchApplicationError)
		response.ResErrCli(context, searchApplicationError)
		return
	}
	response.ResSuccessMsg(context)
}

// ListGlobalizationCopyWritingStruct 获取多语言文案结构/**
func (controller NamespaceControllerController) ListGlobalizationCopyWritingStruct(context *gin.Context) {
	namespaceRequest := &Rquest.NamespaceRequest{}
	shouldBindBodyWithErr := context.ShouldBindBodyWith(&namespaceRequest, binding.JSON)
	if shouldBindBodyWithErr != nil {
		response.ResFail(context, "json解析异常")
		return
	}
	namespace := namespaceRequest.ConvertToTableApplicationNamespace(namespaceRequest)
	searchApplicationList, searchApplicationError := controller.namespaceModel.ListApplicationNamespace(namespace)
	if searchApplicationError != nil {
		controller.logger.Error(searchApplicationError)
		response.ResFail(context, "应用处理异常")
		return
	}
	response.ResSuccess(context, searchApplicationList)
}

// ListGlobalizationCopyWritingNamespace 查询应用文案命名空间/**
func (controller NamespaceControllerController) ListGlobalizationCopyWritingNamespace(context *gin.Context) {
	namespaceRequest := &Rquest.NamespaceRequest{}
	shouldBindBodyWithErr := context.ShouldBindBodyWith(&namespaceRequest, binding.JSON)
	if shouldBindBodyWithErr != nil {
		response.ResFail(context, "json解析异常")
		return
	}
	namespace := namespaceRequest.ConvertToTableApplicationNamespace(namespaceRequest)
	searchApplicationList, searchApplicationError :=
		controller.namespaceModel.SearchDocumentByNamespaceId(namespace)
	if searchApplicationError != nil {
		controller.logger.Error(searchApplicationError)
		response.ResFail(context, "应用处理异常")
		return
	}
	response.ResSuccess(context, searchApplicationList)
}
