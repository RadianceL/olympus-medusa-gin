package request

import "olympus-medusa/data/data"

type ApplicationRequest struct {
	// 应用名称
	ApplicationName string `json:"applicationName,omitempty"`
	// 应用类型 WEB & APPLICATION
	ApplicationType string `json:"applicationType,omitempty"`
	// 应用管理员
	ApplicationAdministrators int `json:"applicationAdministrators,omitempty"`
	// 应用路径 默认应用路径
	ApplicationPath string `json:"applicationPath,omitempty"`
	// 包含的语言范围
	ApplicationLanguage []string `json:"applicationLanguage,omitempty"`
	// 应用环境
	ApplicationEnvironment string `json:"applicationEnvironment,omitempty"`
}

type NamespaceRequest struct {
	ApplicationId          int    `json:"applicationId,omitempty"`
	NamespaceId            int    `json:"namespaceId,omitempty"`
	NamespaceCode          string `json:"namespaceCode,omitempty"`
	NamespaceName          string `json:"namespaceName,omitempty"`
	NamespacePath          string `json:"namespacePath,omitempty"`
	NamespaceParentId      int    `json:"namespaceParentId,omitempty"`
	NamespaceApplicationId int    `json:"namespaceApplicationId,omitempty"`
	CreateUserId           int    `json:"createUserId,omitempty"`
	PageIndex              int    `json:"pageIndex,omitempty"`
	PageSize               int    `json:"pageSize,omitempty"`
}

func (namespaceRequest NamespaceRequest) ConvertToTableApplicationNamespace(namespace *NamespaceRequest) *data.ApplicationNamespace {
	return &data.ApplicationNamespace{
		ApplicationId:          namespace.ApplicationId,
		NamespaceId:            namespace.NamespaceId,
		NamespaceCode:          namespace.NamespaceCode,
		NamespaceName:          namespace.NamespaceName,
		NamespacePath:          namespace.NamespacePath,
		NamespaceParentId:      namespace.NamespaceParentId,
		NamespaceApplicationId: namespace.NamespaceApplicationId,
		CreateUserId:           namespace.CreateUserId,
	}
}

type GlobalDocumentRequest struct {
	Id            int                             `json:"id,omitempty"`
	ApplicationId int                             `json:"applicationId,omitempty"`
	DocumentId    int                             `json:"documentId,omitempty"`
	DocumentPath  string                          `json:"documentPath,omitempty"`
	NamespaceId   int                             `json:"namespaceId,omitempty"`
	DocumentCode  string                          `json:"documentCode,omitempty"`
	DocumentValue string                          `json:"documentDesc,omitempty"`
	PageIndex     int                             `json:"pageIndex,omitempty"`
	PageSize      int                             `json:"pageSize,omitempty"`
	Documents     []GlobalDocumentLanguageRequest `json:"documents,omitempty"`
	Remarks       string                          `json:"remarks,omitempty"`
	DocumentIds   []interface{}                   `json:"documentIds,omitempty"`
}

type GlobalDocumentLanguageRequest struct {
	DocumentId    int    `json:"documentId,omitempty"`
	CountryIso    string `json:"countryIso,omitempty"`
	DocumentValue string `json:"documentValue,omitempty"`
}

type GlobalDocumentIsoQueryRequest struct {
	ApplicationPath string        `json:"applicationPath,omitempty"`
	NamespacePath   []interface{} `json:"namespacePath,omitempty"`
	CountryIso      string        `json:"countryIso,omitempty"`
}
