package model

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"olympus-medusa/common"
	"olympus-medusa/config"
	"olympus-medusa/data/data"
)

type INamespaceModel interface {
	CreateApplicationNamespace(namespaceRequest *data.TableApplicationNamespace) (int64, error)

	ListApplicationNamespace(namespaceRequest *data.TableApplicationNamespace) (data.TableApplicationNamespacePage, error)

	SearchDocumentByNamespaceId(globalDocumentRequest *data.TableApplicationNamespace) (data.TableGlobalDocumentPage, error)
}

// NamespaceModel is application model structure.
type NamespaceModel struct {
	logger *logrus.Logger
	db     *gorm.DB
}

func NewNamespaceModel() INamespaceModel {
	return NamespaceModel{db: common.GetDB(), logger: config.GetLogger()}
}

func (namespace NamespaceModel) CreateApplicationNamespace(namespaceRequest *data.TableApplicationNamespace) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (namespace NamespaceModel) ListApplicationNamespace(namespaceRequest *data.TableApplicationNamespace) (data.TableApplicationNamespacePage, error) {
	//TODO implement me
	panic("implement me")
}

func (namespace NamespaceModel) SearchDocumentByNamespaceId(globalDocumentRequest *data.TableApplicationNamespace) (data.TableGlobalDocumentPage, error) {
	//TODO implement me
	panic("implement me")
}
