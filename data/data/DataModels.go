package data

import (
	"time"
)

type Application struct {
	ID                        int       `gorm:"primaryKey;autoIncrement"`
	ApplicationName           string    `gorm:"size:64;not null;unique;comment:应用名称"`
	ApplicationAdministrators int       `gorm:"comment:app管理员"`
	ApplicationType           string    `gorm:"size:8;not null;comment:app类型 WEB | APPLICATION"`
	ApplicationPath           string    `gorm:"size:128;not null;unique;comment:应用路径（默认应用名称）"`
	MustContainLanguage       string    `gorm:"type:json;comment:必须包含的语言"`
	ApplicationEnvironment    string    `gorm:"size:16;not null;comment:系统环境 STG & DEV & PROD"`
	DualAuthentication        int       `gorm:"default:0;comment:是否开启双重认证 0关闭，1开启"`
	CreateTime                time.Time `gorm:"default:CURRENT_TIMESTAMP;comment:创建时间"`
	CreateUserID              int       `gorm:"comment:创建人用户ID"`
	UpdateTime                time.Time `gorm:"default:CURRENT_TIMESTAMP;comment:更新时间"`
	UpdateUserID              int       `gorm:"comment:更新用户ID"`
}

type TablePage struct {
	TotalSize            int64                  `json:"totalSize,omitempty"`
	ApplicationNamespace []ApplicationNamespace `json:"dataList"`
}

type ApplicationNamespace struct {
	ApplicationId          int    `json:"applicationId,omitempty"`
	NamespaceId            int    `json:"namespaceId,omitempty"`
	NamespaceCode          string `json:"namespaceCode,omitempty"`
	NamespaceName          string `json:"namespaceName,omitempty"`
	NamespacePath          string `json:"namespacePath,omitempty"`
	NamespaceParentId      int    `json:"namespaceParentId,omitempty"`
	NamespaceApplicationId int    `json:"namespaceApplicationId,omitempty"`
	CreateUserId           int    `json:"createUserId,omitempty"`
}

type GlobalDocumentPage struct {
	TotalSize      int64            `json:"totalSize,omitempty"`
	GlobalDocument []GlobalDocument `json:"dataList"`
}

type GlobalDocument struct {
	DocumentId      int                      `json:"documentId,omitempty"`
	ApplicationId   int                      `json:"applicationId,omitempty"`
	ApplicationName string                   `json:"applicationName,omitempty"`
	NamespaceId     int                      `json:"namespaceId,omitempty"`
	NamespaceName   string                   `json:"namespaceName,omitempty"`
	DocumentDesc    string                   `json:"documentDesc,omitempty"`
	DocumentCode    string                   `json:"documentCode,omitempty"`
	CreateTime      string                   `json:"createTime,omitempty"`
	Documents       []GlobalDocumentLanguage `json:"documents,omitempty"`
}

type GlobalDocumentLanguage struct {
	Id                 int    `json:"documentId,omitempty"`
	CountryIso         string `json:"countryIso,omitempty"`
	DocumentCode       string `json:"documentCode,omitempty"`
	DocumentValue      string `json:"documentValue,omitempty"`
	LastUpdateDocument string `json:"lastUpdateDocument,omitempty"`
	CreateTime         string `json:"createTime,omitempty"`
}

type ExportGlobalDocument struct {
	ImportSuccessList []GlobalDocumentExcel `json:"importSuccessList"`
	ImportFailureList []GlobalDocumentExcel `json:"importFailureList"`
	Success           bool                  `json:"success,omitempty"`
}

type GlobalDocumentExcel struct {
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
