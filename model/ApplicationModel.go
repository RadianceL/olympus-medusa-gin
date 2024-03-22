package model

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"olympus-medusa/common"
	"olympus-medusa/config"
	. "olympus-medusa/data/data"
	. "olympus-medusa/data/request"
)

type IApplicationModel interface {
	AddApplication(applicationAddRequest *ApplicationRequest) (int64, error)
	SearchApplicationList(applicationAddRequest *ApplicationRequest) ([]TableApplication, error)
	SearchApplicationById(applicationId int) (TableApplication, error)
}

const (
	applicationModelTableName = "tb_application"
	id                        = "id"
	applicationName           = "application_name"
	applicationAdministrators = "application_administrators"
	applicationType           = "application_type"
	applicationPath           = "application_path"
	mustContainLanguage       = "must_contain_language"
	applicationEnvironment    = "application_environment"
)

// ApplicationModel is application model structure.
type ApplicationModel struct {
	logger *logrus.Logger
	db     *gorm.DB
}

func NewApplicationModel() IApplicationModel {
	return ApplicationModel{db: common.GetDB(), logger: config.GetLogger()}
}

// AddApplication add a role to the menu.
func (applicationModel ApplicationModel) AddApplication(applicationAddRequest *ApplicationRequest) (int64, error) {
	containLanguageList, err := json.Marshal(applicationAddRequest.ApplicationLanguage)
	if err != nil {
		applicationModel.logger.Panic(err)
	}
	if applicationAddRequest.ApplicationPath == "" {
		applicationAddRequest.ApplicationPath = "/" + applicationAddRequest.ApplicationName
	}
	application := TableApplication{
		ApplicationName:           applicationAddRequest.ApplicationName,
		ApplicationAdministrators: applicationAddRequest.ApplicationAdministrators,
		ApplicationType:           applicationAddRequest.ApplicationType,
		ApplicationPath:           applicationAddRequest.ApplicationPath,
		MustContainLanguage:       string(containLanguageList),
		ApplicationEnvironment:    applicationAddRequest.ApplicationEnvironment,
	}
	tx := applicationModel.db.Create(&application)
	return tx.RowsAffected, tx.Error
}

func (applicationModel ApplicationModel) SearchApplicationList(applicationAddRequest *ApplicationRequest) ([]TableApplication, error) {
	var applications []TableApplication
	if err := applicationModel.db.Table("tb_application").
		Where("application_name LIKE ?", "%"+applicationAddRequest.ApplicationName+"%").Find(&applications).Error; err != nil {
		return []TableApplication{}, err
	}
	if applications == nil {
		return []TableApplication{}, nil
	}
	return applications, nil
}

func (applicationModel ApplicationModel) SearchApplicationById(applicationId int) (TableApplication, error) {
	var application Application
	if err := applicationModel.db.Where("id = ?", applicationId).Find(&application).Error; err != nil {
		return TableApplication{}, err
	}
	var tableApplication = &TableApplication{}
	tableApplication.Id = application.ID
	tableApplication.ApplicationAdministrators = application.ApplicationAdministrators
	tableApplication.ApplicationEnvironment = application.ApplicationEnvironment
	tableApplication.ApplicationName = application.ApplicationName
	tableApplication.ApplicationType = application.ApplicationType
	tableApplication.ApplicationPath = application.ApplicationPath
	return *tableApplication, nil
}
