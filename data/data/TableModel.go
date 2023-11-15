package data

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
	ApplicationLanguage []string `json:"applicationLanguage,omitempty"`
	// 应用环境
	ApplicationEnvironment string `json:"applicationEnvironment,omitempty"`
}

type TableApplicationNamespacePage struct {
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
