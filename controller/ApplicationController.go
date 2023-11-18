package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
	"olympus-medusa/common/language"
	"olympus-medusa/config"
	Request "olympus-medusa/data/request"
	"olympus-medusa/data/response"
	"olympus-medusa/model"
)

type IApplicationController interface {
	CreateApplication(context *gin.Context)
	ListApplication(context *gin.Context)
	ListSupportLanguages(context *gin.Context)
}

type ApplicationController struct {
	logger           *logrus.Logger
	applicationModel model.IApplicationModel
}

func NewGlobalApplicationController() IApplicationController {
	return ApplicationController{
		logger:           config.GetLogger(),
		applicationModel: model.NewApplicationModel(),
	}
}

// CreateApplication 创建多语言应用/**
func (controller ApplicationController) CreateApplication(context *gin.Context) {
	applicationAddRequest := &Request.ApplicationRequest{}
	err := context.ShouldBindBodyWith(&applicationAddRequest, binding.JSON)
	if err != nil {
		response.ResErrCli(context, err)
		return
	}
	if applicationAddRequest.ApplicationName == "" {
		response.ResFail(context, "应用名称不能为空")
		return
	}
	if applicationAddRequest.ApplicationEnvironment == "" {
		response.ResFail(context, "环境变量不能为空")
		return
	}
	_, err = controller.applicationModel.AddApplication(applicationAddRequest)
	if err != nil {
		response.ResErrCli(context, err)
		return
	}
	response.ResSuccessMsg(context)
}

// ListApplication 查询应用列表/**
func (controller ApplicationController) ListApplication(context *gin.Context) {
	applicationAddRequest := &Request.ApplicationRequest{}
	err := context.ShouldBindBodyWith(&applicationAddRequest, binding.JSON)
	if err != nil {
		response.ResErrCli(context, err)
		return
	}
	if applicationAddRequest.ApplicationEnvironment == "" {
		response.ResFail(context, "环境变量不能为空")
		return
	}
	searchApplicationList, err := controller.applicationModel.
		SearchApplicationList(applicationAddRequest)
	if err != nil {
		controller.logger.Error("查询列表报错：", err)
		response.ResFail(context, "系统异常，请稍后重试")
		return
	}
	response.ResSuccess(context, searchApplicationList)
}

// ListSupportLanguages 查询应用支持的语言列表/**
func (controller ApplicationController) ListSupportLanguages(context *gin.Context) {
	response.ResSuccess(context, language.Values())
}
