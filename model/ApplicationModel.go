package model

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"olympus-medusa/common"
	"olympus-medusa/config"
	"olympus-medusa/data/data"
	"olympus-medusa/data/request"
	"olympus-medusa/model/basic"
)

type IApplicationModel interface {
	basic.BaseModel

	AddApplication(applicationAddRequest *request.ApplicationRequest) (int64, error)
	SearchApplicationList(applicationAddRequest *request.ApplicationRequest) ([]data.TableApplication, error)
	SearchApplicationById(applicationId int) (data.TableApplication, error)
}

// ApplicationModel is application model structure.
type ApplicationModel struct {
	logger *logrus.Logger
	db     *gorm.DB
}

func NewApplicationModel() IApplicationModel {
	return ApplicationModel{db: common.GetDB(), logger: config.GetLogger()}
}

// Application return a default application model.
func Application() ApplicationModel {
	return ApplicationModel{db: common.GetDB()}
}

// AddApplication add a role to the menu.
func (applicationModel ApplicationModel) AddApplication(applicationAddRequest *request.ApplicationRequest) (int64, error) {
	containLanguageList, err := json.Marshal(applicationAddRequest.ApplicationLanguage)
	if err != nil {
		applicationModel.logger.Panic(err)
	}
	if applicationAddRequest.ApplicationPath == "" {
		applicationAddRequest.ApplicationPath = "/" + applicationAddRequest.ApplicationName
	}
	application := data.Application{
		ApplicationName:           applicationAddRequest.ApplicationName,
		ApplicationAdministrators: applicationAddRequest.ApplicationAdministrators,
		ApplicationType:           applicationAddRequest.ApplicationType,
		ApplicationPath:           applicationAddRequest.ApplicationPath,
		MustContainLanguage:       string(containLanguageList),
		ApplicationEnvironment:    applicationAddRequest.ApplicationEnvironment,
		DualAuthentication:        0,
	}
	tx := applicationModel.db.Create(application)
	return tx.RowsAffected, tx.Error
}

func (applicationModel ApplicationModel) SearchApplicationList(applicationAddRequest *request.ApplicationRequest) ([]data.TableApplication, error) {
	var applications []data.Application
	if err := applicationModel.db.Where("application_name LIKE ?", applicationAddRequest.ApplicationName).Find(&applications).Error; err != nil {
		return []data.TableApplication{}, err
	}

	var result []data.TableApplication
	for _, value := range applications {
		var tableApplication = &data.TableApplication{}
		tableApplication.ApplicationAdministrators = value.ApplicationAdministrators
		tableApplication.ApplicationEnvironment = value.ApplicationEnvironment
		tableApplication.ApplicationName = value.ApplicationName
		tableApplication.ApplicationType = value.ApplicationType
		tableApplication.ApplicationPath = value.ApplicationPath
		tableApplication.Id = value.ID
		result = append(result, *tableApplication)
	}
	if result == nil {
		return []data.TableApplication{}, nil
	}
	return result, nil
}

func (applicationModel ApplicationModel) SearchApplicationById(applicationId int) (data.TableApplication, error) {
	var application data.Application
	if err := applicationModel.db.Where("id = ?", applicationId).Find(&application).Error; err != nil {
		return data.TableApplication{}, err
	}
	var tableApplication = &data.TableApplication{}
	tableApplication.Id = application.ID
	tableApplication.ApplicationAdministrators = application.ApplicationAdministrators
	tableApplication.ApplicationEnvironment = application.ApplicationEnvironment
	tableApplication.ApplicationName = application.ApplicationName
	tableApplication.ApplicationType = application.ApplicationType
	tableApplication.ApplicationPath = application.ApplicationPath
	return *tableApplication, nil
}
