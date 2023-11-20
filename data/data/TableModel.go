package data

import (
	"time"
)

type Tabler interface {
	TableName() string
}

type TableApplication struct {
	Id int `json:"id,omitempty"`
	// 应用名称
	ApplicationName string `json:"applicationName,omitempty"`
	// 应用类型 WEB & APPLICATION
	ApplicationType string `json:"applicationType,omitempty"`
	// 应用管理员
	ApplicationAdministrators int `json:"applicationAdministrators,omitempty"`
	// 应用路径 默认应用路径
	ApplicationPath string `json:"applicationPath,omitempty"`
	// 包含的语言范围
	MustContainLanguage string `json:"applicationLanguage,omitempty"`
	// 应用环境
	ApplicationEnvironment string `json:"applicationEnvironment,omitempty"`
}

func (TableApplication) TableName() string {
	return "tb_application"
}

type TablePage struct {
	TotalSize            int64                       `json:"totalSize,omitempty"`
	ApplicationNamespace []TableApplicationNamespace `json:"dataList"`
}

type TableApplicationNamespace struct {
	ApplicationId          int    `json:"applicationId,omitempty"`
	NamespaceId            int    `json:"namespaceId,omitempty"`
	NamespaceCode          string `json:"namespaceCode,omitempty"`
	NamespaceName          string `json:"namespaceName,omitempty"`
	NamespacePath          string `json:"namespacePath,omitempty"`
	NamespaceParentId      int    `json:"namespaceParentId,omitempty"`
	NamespaceApplicationId int    `json:"namespaceApplicationId,omitempty"`
	CreateUserId           int    `json:"createUserId,omitempty"`
}

type TableGlobalDocumentPage struct {
	TotalSize      int64                 `json:"totalSize,omitempty"`
	GlobalDocument []TableGlobalDocument `json:"dataList"`
}

type ExportGlobalDocument struct {
	ImportSuccessList []TableGlobalDocumentExcel `json:"importSuccessList"`
	ImportFailureList []TableGlobalDocumentExcel `json:"importFailureList"`
	Success           bool                       `json:"success,omitempty"`
}

type TableGlobalDocumentExcel struct {
	ApplicationId   int    `json:"applicationId,omitempty"`
	ApplicationName string `json:"applicationName,omitempty"`
	NamespaceId     int    `json:"namespaceId,omitempty"`
	NamespaceName   string `json:"namespaceName,omitempty"`
	DocumentCode    string `json:"documentCode,omitempty"`
	CountryIso      string `json:"countryIso,omitempty"`
	CountryName     string `json:"countryName,omitempty"`
	DocumentValue   string `json:"documentValue,omitempty"`
	DocumentDesc    string `json:"documentDesc,omitempty"`
}
type ApplicationGlobalizationDocumentCode struct {
	DocumentID            int       `gorm:"column:document_id;primaryKey;autoIncrement" json:"document_id"`
	ApplicationID         int       `gorm:"column:application_id;not null" json:"application_id"`
	NamespaceID           int       `gorm:"column:namespace_id;not null" json:"namespace_id"`
	DocumentCode          string    `gorm:"column:document_code;not null" json:"document_code"`
	DocumentDesc          string    `gorm:"column:document_desc" json:"document_desc"`
	IsEnable              int       `gorm:"column:is_enable;default:1;not null" json:"is_enable"`
	OnlineTime            time.Time `gorm:"column:online_time" json:"online_time"`
	OnlineOperatorUserID  int       `gorm:"column:online_operator_user_id" json:"online_operator_user_id"`
	OfflineTime           time.Time `gorm:"column:offline_time" json:"offline_time"`
	OfflineOperatorUserID int       `gorm:"column:offline_operator_user_id" json:"offline_operator_user_id"`
	OfflineAccessUserID   int       `gorm:"column:offline_access_user_id" json:"offline_access_user_id"`
	CreateTime            time.Time `gorm:"column:create_time;default:CURRENT_TIMESTAMP;not null" json:"create_time"`
	UpdateTime            time.Time `gorm:"column:update_time" json:"update_time"`
	CreateUserID          int       `gorm:"column:create_user_id" json:"create_user_id"`
	DeleteFlag            int       `gorm:"column:delete_flag;default:0" json:"delete_flag"`
	DeleteTime            time.Time `gorm:"column:delete_time" json:"delete_time"`
	DeleteUserID          int       `gorm:"column:delete_user_id" json:"delete_user_id"`
	Remarks               string    `gorm:"column:remarks" json:"remarks"`
}

// TableName sets the table name for the struct
func (ApplicationGlobalizationDocumentCode) TableName() string {
	return "tb_application_globalization_document_code"
}

type TableGlobalDocument struct {
	DocumentId      int                           `json:"documentId,omitempty"`
	ApplicationId   int                           `json:"applicationId,omitempty"`
	ApplicationName string                        `json:"applicationName,omitempty"`
	NamespaceId     int                           `json:"namespaceId,omitempty"`
	NamespaceName   string                        `json:"namespaceName,omitempty"`
	DocumentDesc    string                        `json:"documentDesc,omitempty"`
	DocumentCode    string                        `json:"documentCode,omitempty"`
	Documents       []TableGlobalDocumentLanguage `json:"documents,omitempty"`
	CreateTime      string                        `json:"createTime,omitempty"`
}

type TableGlobalDocumentLanguage struct {
	Id                 int    `json:"documentId,omitempty"`
	CountryIso         string `json:"countryIso,omitempty"`
	DocumentCode       string `json:"documentCode,omitempty"`
	DocumentValue      string `json:"documentValue,omitempty"`
	LastUpdateDocument string `json:"lastUpdateDocument,omitempty"`
	CreateTime         string `json:"createTime,omitempty"`
}
