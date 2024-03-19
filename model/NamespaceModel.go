package model

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"olympus-medusa/common"
	"olympus-medusa/config"
	"olympus-medusa/data/data"
)

type INamespaceModel interface {
	CreateApplicationNamespace(namespaceRequest *data.ApplicationNamespace) (int64, error)

	ListApplicationNamespace(namespaceRequest *data.ApplicationNamespace) (data.TablePage, error)

	SearchDocumentByNamespaceId(globalDocumentRequest *data.ApplicationNamespace) (data.GlobalDocumentPage, error)

	SearchNamespaceById(applicationId int, appNamespaceId int) (data.ApplicationNamespace, error)
}

const (
	// NamespaceModelTableName tb_application_namespace
	namespaceModelTableName = "tb_application_namespace"
	namespaceId             = "namespace_id"
	namespaceCode           = "namespace_code"
	namespaceName           = "namespace_name"
	namespacePath           = "namespace_path"
	namespaceParentId       = "namespace_parent_id"
	namespaceApplicationId  = "application_id"
	createUser              = "create_user"
)

// NamespaceModel is application model structure.
type NamespaceModel struct {
	logger *logrus.Logger
	db     *gorm.DB
}

func NewNamespaceModel() INamespaceModel {
	return NamespaceModel{db: common.GetDB(), logger: config.GetLogger()}
}

func (namespace NamespaceModel) CreateApplicationNamespace(namespaceRequest *data.ApplicationNamespace) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (namespace NamespaceModel) ListApplicationNamespace(namespaceRequest *data.ApplicationNamespace) (data.TablePage, error) {
	//TODO implement me
	panic("implement me")
}

func (namespace NamespaceModel) SearchDocumentByNamespaceId(globalDocumentRequest *data.ApplicationNamespace) (data.GlobalDocumentPage, error) {
	//TODO implement me
	panic("implement me")
}

func (namespace NamespaceModel) SearchNamespaceById(applicationId int, appNamespaceId int) (data.ApplicationNamespace, error) {
	var application data.ApplicationNamespace
	if err := namespace.db.Where("application_id = ? AND namespace_id = ?", applicationId, appNamespaceId).Find(&application).Error; err != nil {
		return data.ApplicationNamespace{}, err
	}
	return application, nil
}
