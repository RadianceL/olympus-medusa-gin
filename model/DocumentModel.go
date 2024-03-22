package model

import (
	"container/list"
	"errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"olympus-medusa/common"
	"olympus-medusa/common/convert"
	"olympus-medusa/common/language"
	"olympus-medusa/config"
	. "olympus-medusa/data/data"
	. "olympus-medusa/data/request"
	"strconv"
	"strings"
	"time"
)

type IDocumentModel interface {
	//ImportDocument(namespaceRequest multipart.File) (data.ExportGlobalDocument, error)
	//
	//getImportDataList(rows [][]string, isSuccess bool) (data.ExportGlobalDocument, error)

	CreateDocument(namespaceRequest *GlobalDocumentRequest) (int64, error)

	SearchDocumentValue(globalDocumentRequest *GlobalDocumentRequest) (arr []interface{})

	QueryDocument(globalDocumentRequest *GlobalDocumentRequest) ([]GlobalDocumentExcel, error)

	SearchDocumentByNamespaceId(globalDocumentRequest *GlobalDocumentRequest) (GlobalDocumentPage, error)

	UpdateDocumentByDocumentId(namespaceRequest *GlobalDocumentRequest) (int64, error)

	DeleteDocumentByDocumentId(namespaceRequest *GlobalDocumentRequest) (int64, error)

	SearchDocumentCode(globalDocumentRequest *GlobalDocumentRequest) (*GlobalDocument, error)

	SearchDocumentById(globalDocumentRequest *GlobalDocumentRequest) (GlobalDocument, error)

	SearchDocumentByCountryIso(globalDocumentIsoQueryRequest *GlobalDocumentIsoQueryRequest) (map[string]string, error)

	SearchApplicationByCountryIso(globalDocumentIsoQueryRequest *GlobalDocumentIsoQueryRequest) (map[string]map[string]string, error)

	SearchOptionByNamespace(globalDocumentRequest *GlobalDocumentRequest) ([]TableGlobalizationDocumentCode, error)
}

// DocumentModel is application model structure.
type DocumentModel struct {
	logger *logrus.Logger
	db     *gorm.DB
}

const (
	// ApplicationModelTableName tb_application
	documentTableName  = "tb_application_globalization_document_code"
	applicationIdField = "application_id"
	namespaceIdField   = "namespace_id"
	documentCodeField  = "document_code"
	documentDescField  = "document_desc"

	deleteFlagField = "delete_flag"

	documentValueTableName = "tb_application_globalization_document_value"
	documentIdField        = "document_id"
	countryIsoField        = "country_iso"
)

func NewDocumentModel() IDocumentModel {
	return DocumentModel{db: common.GetDB(), logger: config.GetLogger()}
}

func (documentModel DocumentModel) CreateDocument(globalDocumentRequest *GlobalDocumentRequest) (int64, error) {
	tx := documentModel.db.Begin()
	documentCode := TableGlobalizationDocumentCode{
		ApplicationID:        globalDocumentRequest.ApplicationId,
		NamespaceID:          globalDocumentRequest.NamespaceId,
		DocumentCode:         globalDocumentRequest.DocumentCode,
		DocumentDesc:         globalDocumentRequest.DocumentValue,
		OnlineTime:           time.Now(),
		OnlineOperatorUserID: 0,
		DeleteFlag:           0,
		IsEnable:             1,
		CreateUserID:         0,
		CreateTime:           time.Now(),
		DeleteUserID:         0,
	}
	insertDocumentCodeTx := documentModel.db.
		Table("tb_application_globalization_document_code").
		Create(&documentCode)

	if insertDocumentCodeTx.Error != nil {
		_ = tx.Rollback()
		return 0, insertDocumentCodeTx.Error
	}

	documents := globalDocumentRequest.Documents
	for _, document := range documents {
		languageCountry := language.FindLanguage(document.CountryIso)
		if languageCountry == nil {
			_ = tx.Rollback()
			return 0, errors.New("未识别的国家编码，请检查后重试")
		}
		documentValue := TableGlobalizationDocumentValue{
			DocumentId:       documentCode.DocumentID,
			ApplicationId:    documentCode.ApplicationID,
			NamespaceId:      documentCode.NamespaceID,
			CountryIso:       document.CountryIso,
			CountryName:      languageCountry.CountryName,
			DocumentCode:     documentCode.DocumentCode,
			DocumentValue:    document.DocumentValue,
			DocumentIsOnline: 1,
			CreateTime:       time.Now(),
		}
		documentValueTx := documentModel.db.
			Create(&documentValue)
		if documentValueTx.Error != nil {
			_ = tx.Rollback()
			return 0, documentValueTx.Error
		}
	}
	commitError := tx.Commit()
	if commitError != nil {
		_ = tx.Rollback()
		return 0, tx.Error
	}
	return tx.RowsAffected, nil
}

func (documentModel DocumentModel) SearchDocumentValue(globalDocumentRequest *GlobalDocumentRequest) (arr []interface{}) {
	var documentValue []TableGlobalizationDocumentValue
	var valueArr []interface{}
	if err := documentModel.db.Where("document_value LIKE ?", "%"+globalDocumentRequest.DocumentValue+"%").
		Find(&documentValue).
		Select("*").Error; err != nil {
		return valueArr
	}
	for _, document := range documentValue {
		valueArr = append(valueArr, document.DocumentId)
	}
	return valueArr
}

func (documentModel DocumentModel) QueryDocument(globalDocumentRequest *GlobalDocumentRequest) ([]GlobalDocumentExcel, error) {
	var resultData []GlobalDocumentExcel
	tx := documentModel.db.Table("tb_application_globalization_document_code").
		Joins("LEFT JOIN "+
			"tb_application_globalization_document_code "+
			"ON tb_application_globalization_document_code.document_id = tb_application_globalization_document_value.document_id").
		Joins("LEFT JOIN "+
			"tb_application "+
			"ON tb_application.id = tb_application_globalization_document_value.document_id").
		Joins("LEFT JOIN "+
			"tb_application_namespace "+
			"ON tb_application_namespace.namespace_id = tb_application_globalization_document_code.namespace_id").
		Where("tb_application_globalization_document_code.delete_flag = ?", 0)
	if globalDocumentRequest.NamespaceId != 0 {
		tx.Where("tb_application_globalization_document_code.namespace_id = ?", globalDocumentRequest.NamespaceId)
	}
	if globalDocumentRequest.ApplicationId != 0 {
		tx.Where("tb_application_globalization_document_code.application_id = ?", globalDocumentRequest.ApplicationId)
	}
	if len(globalDocumentRequest.DocumentIds) > 0 {
		tx.Where("tb_application_globalization_document_code.document_id IN ?", globalDocumentRequest.DocumentIds)
	}
	if globalDocumentRequest.DocumentCode != "" {
		tx.Where("tb_application_globalization_document_code.document_code LIKE", "%"+globalDocumentRequest.DocumentCode+"%")
	}

	if globalDocumentRequest.DocumentValue != "" {
		tx.Where("tb_application_globalization_document_code.document_id LIKE ?", globalDocumentRequest.DocumentValue)
	}
	if err := tx.Find(&resultData).Error; err != nil {
		return []GlobalDocumentExcel{}, err
	}
	return resultData, nil
}

func (documentModel DocumentModel) SearchOptionByNamespace(globalDocumentRequest *GlobalDocumentRequest) ([]TableGlobalizationDocumentCode, error) {
	resultMap := make(map[string]string)
	applicationStatement := documentModel.db.
		Where("application_id = ? AND namespace_id = ?",
			globalDocumentRequest.ApplicationId, globalDocumentRequest.NamespaceId)
	if globalDocumentRequest.DocumentValue != "" {
		applicationStatement.Where("document_desc LIKE ", "%"+convert.ToString(globalDocumentRequest.DocumentValue)+"%")
	}
	if globalDocumentRequest.DocumentCode != "" {
		applicationStatement.Where("document_code LIKE ", "%"+convert.ToString(globalDocumentRequest.DocumentCode)+"%")
	}
	var applicationCodes []TableGlobalizationDocumentCode
	if err := applicationStatement.Find(&applicationCodes).Error; err != nil {
		return applicationCodes, err
	}
	if len(applicationCodes) <= 0 {
		return applicationCodes, errors.New("应用空间数据异常，请检查配置后重试")
	}

	for _, namespaceDataMap := range applicationCodes {
		resultMap[convert.ToString(namespaceDataMap.DocumentCode)] = convert.ToString(namespaceDataMap.DocumentDesc)
	}
	return applicationCodes, nil
}

func (documentModel DocumentModel) SearchDocumentByNamespaceId(globalDocumentRequest *GlobalDocumentRequest) (GlobalDocumentPage, error) {
	tx := documentModel.db.Table(documentTableName).
		Joins("LEFT JOIN "+
			documentTableName+
			"ON "+documentTableName+".document_id = tb_application_globalization_document_value.document_id").
		Joins("LEFT JOIN "+
			"tb_application "+
			"ON tb_application.id = tb_application_globalization_document_value.document_id").
		Joins("LEFT JOIN "+
			"tb_application_namespace "+
			"ON tb_application_namespace.namespace_id = "+documentTableName+".namespace_id").
		Where(documentTableName+".delete_flag = ?", 0)

	if globalDocumentRequest.NamespaceId != 0 {
		tx.Where(documentTableName+"."+namespaceIdField, "=", globalDocumentRequest.NamespaceId)
	}
	if globalDocumentRequest.ApplicationId != 0 {
		tx.Where(documentTableName+"."+applicationIdField, "=", globalDocumentRequest.ApplicationId)
	}
	if globalDocumentRequest.DocumentCode != "" {
		tx.Where(documentCodeField+"LIKE ?", "%"+globalDocumentRequest.DocumentCode+"%")
	}
	var arr = make([]interface{}, 0)
	if globalDocumentRequest.DocumentValue != "" {
		arr = documentModel.SearchDocumentValue(globalDocumentRequest)
		if arr == nil || len(arr) <= 0 {
			arr = make([]interface{}, 1)
			arr[0] = 0
		}
		tx.Where(documentIdField+"IN ?", arr)
	}
	if globalDocumentRequest.PageIndex != 0 && globalDocumentRequest.PageSize != 0 {
		tx.Offset((globalDocumentRequest.PageIndex - 1) * globalDocumentRequest.PageSize)
		tx.Limit(globalDocumentRequest.PageSize)
	}
	var documentList []GlobalDocument
	tx.Find(&documentList)

	countTx := documentModel.db.Table(documentTableName)
	countTx.Where(documentTableName+"."+deleteFlagField, "=", 0)
	if globalDocumentRequest.NamespaceId != 0 {
		countTx.Where(documentTableName+"."+namespaceIdField, "=", globalDocumentRequest.NamespaceId)
	}
	if globalDocumentRequest.ApplicationId != 0 {
		countTx.Where(documentTableName+"."+applicationIdField, "=", globalDocumentRequest.ApplicationId)
	}
	if globalDocumentRequest.DocumentCode != "" {
		countTx.Where(documentCodeField, "LIKE", "%"+globalDocumentRequest.DocumentCode+"%")
	}
	if arr != nil && len(arr) > 0 {
		countTx.Where(documentIdField, arr)
	}

	var count int64
	countTx.Count(&count)

	var resultData GlobalDocumentPage
	if tx.Error != nil {
		resultData.TotalSize = count
		return resultData, tx.Error
	}
	if len(documentList) <= 0 {
		resultData.TotalSize = count
		resultData.GlobalDocument = make([]GlobalDocument, 0)
		return resultData, nil
	}
	var resultList []GlobalDocument
	for _, document := range documentList {
		var result []GlobalDocumentLanguage
		queryDocumentValueStatement := documentModel.db.Table(documentValueTableName)
		queryDocumentValueStatement.Where(documentIdField, "=", document.DocumentId)
		queryDocumentValueStatement.Find(&result)
		if queryDocumentValueStatement.Error != nil {
			resultData.TotalSize = count
			return resultData, nil
		}
		document.Documents = result
		resultList = append(resultList, document)
	}
	resultData.TotalSize = count
	resultData.GlobalDocument = resultList
	return resultData, nil
}

func (documentModel DocumentModel) UpdateDocumentByDocumentId(namespaceRequest *GlobalDocumentRequest) (int64, error) {
	err := documentModel.db.Transaction(func(tx *gorm.DB) error {
		var applicationGlobalizationDocumentCode TableGlobalizationDocumentCode
		tx.Model(&applicationGlobalizationDocumentCode).
			Where(documentIdField+"=", namespaceRequest.DocumentId).
			Update(documentDescField, namespaceRequest.DocumentValue)

		documents := namespaceRequest.Documents
		for _, document := range documents {
			languageCountry := language.FindLanguage(document.CountryIso)
			if languageCountry == nil {
				return errors.New("未识别的国家编码，请检查后重试")
			}

			var tableGlobalDocumentLanguageList []GlobalDocumentLanguage
			queryDocumentValueStatement := tx.Table(documentValueTableName).Select("*")
			queryDocumentValueStatement.Where(documentIdField+"=", namespaceRequest.DocumentId)
			queryDocumentValueStatement.Where(countryIsoField+"=", document.CountryIso)
			queryDocumentValueStatement.Find(&tableGlobalDocumentLanguageList)

			if queryDocumentValueStatement.Error != nil {
				return errors.New("更新多语言文案语言编码查重异常，请稍后重试")
			}
			if len(tableGlobalDocumentLanguageList) <= 0 {
				languageCountry := language.FindLanguage(document.CountryIso)
				if languageCountry == nil {
					_ = tx.Rollback()
					return errors.New("未识别的国家编码，请检查后重试")
				}
				documentValue := TableGlobalizationDocumentValue{
					DocumentId:    namespaceRequest.DocumentId,
					NamespaceId:   namespaceRequest.NamespaceId,
					CountryIso:    document.CountryIso,
					CountryName:   languageCountry.CountryName,
					DocumentValue: document.DocumentValue,
				}
				createDocumentValueStatement := tx.Create(&documentValue)
				if createDocumentValueStatement.Error != nil {
					return createDocumentValueStatement.Error
				}
			} else {
				var tableGlobalDocumentValueResult TableGlobalizationDocumentValue
				if document.DocumentId == 0 {
					queryDocumentValueItemStatement := tx.Table(documentValueTableName)
					queryDocumentValueItemStatement.Where(documentIdField+"=", namespaceRequest.DocumentId)
					queryDocumentValueItemStatement.Where(namespaceIdField+"=", namespaceRequest.NamespaceId)
					queryDocumentValueItemStatement.Where(countryIsoField+"=", document.CountryIso)
					queryDocumentValueItemStatement = queryDocumentValueStatement.First(&applicationGlobalizationDocumentCode)

					if queryDocumentValueItemStatement.Error != nil {
						return queryDocumentValueItemStatement.Error
					}
				} else {
					documentResultDataResult := tx.Find(&applicationGlobalizationDocumentCode, "id=?", document.DocumentId)
					if documentResultDataResult.Error != nil {
						return documentResultDataResult.Error
					}
				}
				updateDocumentValueStatement := tx.Model(TableGlobalizationDocumentValue{}).
					Where("id =", tableGlobalDocumentValueResult.Id).
					Updates(TableGlobalizationDocumentValue{
						DocumentValue:      document.DocumentValue,
						LastUpdateDocument: tableGlobalDocumentValueResult.DocumentValue,
					})
				if updateDocumentValueStatement.Error != nil {
					if updateDocumentValueStatement.Error.Error() != "no affect row" {
						return updateDocumentValueStatement.Error
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	return 1, nil
}

func (documentModel DocumentModel) DeleteDocumentByDocumentId(namespaceRequest *GlobalDocumentRequest) (int64, error) {
	nano := time.Now().Unix()
	documentTableNameStatement := documentModel.db.Table(documentTableName).
		Where(documentIdField, "=", namespaceRequest.DocumentId).
		Updates(TableGlobalizationDocumentCode{
			DeleteFlag:   1,
			Remarks:      namespaceRequest.Remarks,
			DocumentCode: namespaceRequest.DocumentCode + "_@delete_" + strconv.FormatInt(nano, 10),
			DeleteTime:   time.Now(),
		})
	if documentTableNameStatement.Error != nil {
		if documentTableNameStatement.Error.Error() != "no affect row" {
			return 0, documentTableNameStatement.Error
		}
	}
	return 1, nil
}

func (documentModel DocumentModel) SearchDocumentCode(globalDocumentRequest *GlobalDocumentRequest) (*GlobalDocument, error) {
	var globalDocument TableGlobalizationDocumentCode
	statement := documentModel.db.Table(documentTableName)
	if globalDocumentRequest.Id != 0 {
		statement.Where(documentIdField+"=", globalDocumentRequest.Id)
	} else if globalDocumentRequest.DocumentCode != "" {
		statement.Where(namespaceIdField+"=", globalDocumentRequest.NamespaceId)
		statement.Where(documentCodeField+"=", globalDocumentRequest.DocumentCode)
	} else {
		return nil, errors.New("文案ID或文案编码与命名空间ID必传其中一个")
	}
	globalDocumentStatement := statement.Find(&globalDocument)
	if globalDocumentStatement.Error != nil {
		return nil, globalDocumentStatement.Error
	}
	var result = GlobalDocument{
		ApplicationId: globalDocument.ApplicationID,
		DocumentCode:  globalDocument.DocumentCode,
		DocumentDesc:  globalDocument.DocumentDesc,
		NamespaceId:   globalDocument.NamespaceID,
		DocumentId:    globalDocument.DocumentID,
	}
	return &result, nil
}

func (documentModel DocumentModel) SearchDocumentById(globalDocumentRequest *GlobalDocumentRequest) (GlobalDocument, error) {
	var globalDocument TableGlobalizationDocumentCode
	statement := documentModel.db.Table(documentTableName)
	if globalDocumentRequest.Id != 0 {
		statement.Where(documentIdField+"=", globalDocumentRequest.Id)
	} else if globalDocumentRequest.DocumentCode != "" {
		statement.Where(namespaceIdField+"=", globalDocumentRequest.NamespaceId)
		statement.Where(documentCodeField+"=", globalDocumentRequest.DocumentCode)
	} else {
		return GlobalDocument{}, errors.New("文案ID或文案编码与命名空间ID必传其中一个")
	}
	globalDocumentStatement := statement.Find(&globalDocument)
	if globalDocumentStatement.Error != nil {
		return GlobalDocument{}, globalDocumentStatement.Error
	}

	var globalDocumentValue []TableGlobalizationDocumentValue
	queryDocumentValueStatement := documentModel.db.Table(documentValueTableName)
	queryDocumentValueStatement.Where(documentIdField+"=", globalDocument.DocumentID)

	documentValueErr := queryDocumentValueStatement.Find(&globalDocumentValue)
	if documentValueErr.Error != nil {
		return GlobalDocument{}, documentValueErr.Error
	}

	var result GlobalDocument
	var resultItem []GlobalDocumentLanguage
	for _, documentValueResultData := range globalDocumentValue {
		var tableGlobalDocumentLanguageOutputResult = GlobalDocumentLanguage{
			CountryIso:         documentValueResultData.CountryIso,
			DocumentValue:      documentValueResultData.DocumentValue,
			CreateTime:         convert.ToString(documentValueResultData.CreateTime),
			LastUpdateDocument: documentValueResultData.LastUpdateDocument,
			Id:                 documentValueResultData.Id,
			DocumentCode:       documentValueResultData.DocumentCode,
		}

		resultItem = append(resultItem, tableGlobalDocumentLanguageOutputResult)
	}
	result.DocumentId = globalDocument.DocumentID
	result.DocumentCode = globalDocument.DocumentCode
	result.DocumentDesc = globalDocument.DocumentDesc
	result.ApplicationId = globalDocument.ApplicationID
	result.NamespaceId = globalDocument.NamespaceID
	result.CreateTime = convert.ToString(globalDocument.CreateTime)
	result.Documents = resultItem
	return result, nil
}

func (documentModel DocumentModel) SearchDocumentByCountryIso(globalDocumentIsoQueryRequest *GlobalDocumentIsoQueryRequest) (map[string]string, error) {
	resultMap := make(map[string]string)

	var namespaceResultMap []TableApplicationNamespace
	if len(globalDocumentIsoQueryRequest.NamespacePath) >= 5 {
		return resultMap, errors.New("单次查询不能超过4个命名空间")
	}
	namespaceStatement := documentModel.db.Table(namespaceModelTableName)
	namespaceStatement.Where(namespacePath+"=", globalDocumentIsoQueryRequest.NamespacePath)
	namespaceResultStatement := namespaceStatement.Find(&namespaceResultMap)
	if namespaceResultStatement.Error != nil {
		return resultMap, namespaceResultStatement.Error
	}
	if len(namespaceResultMap) >= 5 || len(namespaceResultMap) <= 0 {
		return resultMap, errors.New("命名空间数据异常，请检查配置后重试")
	}
	for _, namespaceResult := range namespaceResultMap {
		var resultList []GlobalDocument
		var documentResult GlobalDocument

		var documentCodes []TableGlobalizationDocumentCode
		queryDocumentCodeStatement := documentModel.db.Table(documentTableName).Select("*")
		queryDocumentCodeStatement.Where(documentTableName+"."+namespaceIdField+"=", namespaceResult.NamespaceId)
		queryDocumentCodeStatementResult := queryDocumentCodeStatement.Find(&documentCodes)
		if queryDocumentCodeStatementResult.Error != nil {
			return resultMap, queryDocumentCodeStatementResult.Error
		}

		var documentValues []TableGlobalizationDocumentValue
		queryDocumentValueStatement := documentModel.db.Table(documentValueTableName).Select("*")
		queryDocumentValueStatement.Where(documentValueTableName+"."+namespaceIdField+"=", namespaceResult.NamespaceId)
		queryDocumentValueStatement.Where(countryIsoField, "=", globalDocumentIsoQueryRequest.CountryIso)
		queryDocumentValueStatement.Joins("LEFT JOIN " +
			documentTableName + "ON " +
			documentValueTableName + "." + documentIdField + " = " + documentTableName + "." + documentIdField)

		queryDocumentValueStatementResult := queryDocumentValueStatement.Find(&documentValues)
		if queryDocumentValueStatementResult.Error != nil {
			return resultMap, queryDocumentValueStatementResult.Error
		}

		documentResultList := list.New()
		for _, documentValue := range documentValues {
			documentResultList.PushBack(convert.ToString(documentValue.DocumentCode))
		}
		documentCodeList := list.New()
		for _, documentCode := range documentCodes {
			documentCodeList.PushBack(convert.ToString(documentCode.DocumentCode))
		}
		for e := documentResultList.Front(); e != nil; e = e.Next() {
			for f := documentCodeList.Front(); f != nil; f = f.Next() {
				if strings.EqualFold(e.Value.(string), f.Value.(string)) {
					documentCodeList.Remove(f)
					break
				}
			}
		}
		if !strings.EqualFold(globalDocumentIsoQueryRequest.CountryIso, "CN") {
			var documentValueResultEnList []TableGlobalizationDocumentValue
			queryDocumentValueEnStatement := documentModel.db.Table(documentValueTableName).Select("*")
			queryDocumentValueEnStatement.Where(documentValueTableName+"."+namespaceIdField+"=", namespaceResult.NamespaceId)
			queryDocumentValueEnStatement.Where(countryIsoField+"=", "EN")
			queryDocumentValueEnStatement.Joins("LEFT JOIN " +
				documentTableName + "ON " +
				documentValueTableName + "." + documentIdField + " = " + documentTableName + "." + documentIdField)
			var arrays []interface{}
			for e := documentCodeList.Front(); e != nil; e = e.Next() {
				arrays = append(arrays, e.Value)
			}
			if len(arrays) != 0 {
				queryDocumentValueEnStatement.Where(documentTableName+"."+documentCodeField+" IN", arrays)
				queryDocumentValueEnStatementResult := queryDocumentValueEnStatement.Find(&documentValueResultEnList)
				if queryDocumentValueEnStatementResult.Error != nil {
					return resultMap, queryDocumentValueEnStatementResult.Error
				}
				var result []GlobalDocumentLanguage
				for _, documentValueResultData := range documentValueResultEnList {
					var tableGlobalDocumentLanguageOutputResult = GlobalDocumentLanguage{
						CountryIso:         documentValueResultData.CountryIso,
						DocumentValue:      documentValueResultData.DocumentValue,
						CreateTime:         convert.ToString(documentValueResultData.CreateTime),
						LastUpdateDocument: documentValueResultData.LastUpdateDocument,
						Id:                 documentValueResultData.Id,
						DocumentCode:       documentValueResultData.DocumentCode,
					}
					result = append(result, tableGlobalDocumentLanguageOutputResult)
				}
				documentResult.Documents = result
				resultList = append(resultList, documentResult)
				for _, value := range resultList {
					if len(value.Documents) > 0 {
						for _, document := range value.Documents {
							resultMap[document.DocumentCode] = document.DocumentValue
						}
					}
				}
			}
		}
		var result []GlobalDocumentLanguage
		for _, documentValueResultData := range documentValues {
			var tableGlobalDocumentLanguageOutputResult = GlobalDocumentLanguage{
				Id:                 documentValueResultData.DocumentId,
				CountryIso:         documentValueResultData.CountryIso,
				DocumentCode:       documentValueResultData.DocumentCode,
				DocumentValue:      documentValueResultData.DocumentValue,
				LastUpdateDocument: documentValueResultData.LastUpdateDocument,
				CreateTime:         convert.ToString(documentValueResultData.CreateTime),
			}
			result = append(result, tableGlobalDocumentLanguageOutputResult)
		}
		documentResult.Documents = result
		resultList = append(resultList, documentResult)
		for _, value := range resultList {
			if len(value.Documents) > 0 {
				for _, document := range value.Documents {
					resultMap[document.DocumentCode] = document.DocumentValue
				}
			}
		}
	}
	return resultMap, nil
}

func (documentModel DocumentModel) SearchApplicationByCountryIso(globalDocumentIsoQueryRequest *GlobalDocumentIsoQueryRequest) (map[string]map[string]string, error) {
	resultMap := make(map[string]map[string]string)
	var application []TableApplication
	applicationStatement := documentModel.db.Table(applicationModelTableName)
	applicationStatement.Where(applicationPath, "=", globalDocumentIsoQueryRequest.ApplicationPath)
	applicationStatementResult := applicationStatement.Find(&application)
	if applicationStatementResult.Error != nil {
		return resultMap, applicationStatementResult.Error
	}
	if len(application) > 1 || len(application) <= 0 {
		return resultMap, errors.New("应用空间数据异常，请检查配置后重试")
	}

	var namespaceResultMap []TableApplicationNamespace
	namespaceStatement := documentModel.db.Table(namespaceModelTableName)
	namespaceStatement.Where(applicationIdField+"=", application[0].Id)
	namespaceStatementResult := namespaceStatement.Find(&namespaceResultMap)
	if namespaceStatementResult.Error != nil {
		return resultMap, namespaceStatementResult.Error
	}
	if len(namespaceResultMap) <= 0 {
		return resultMap, errors.New("命名空间数据异常，请检查配置后重试")
	}

	for _, namespaceDataMap := range namespaceResultMap {
		var resultList []GlobalDocument
		var documentResult GlobalDocument
		var documentValueResultDataLists []TableGlobalizationDocumentValue
		queryDocumentValueStatement := documentModel.db.Table(documentValueTableName).Select("*")
		queryDocumentValueStatement.Joins("LEFT JOIN " +
			documentTableName + "ON " +
			documentValueTableName + "." + documentIdField + " = " + documentTableName + "." + documentIdField)
		queryDocumentValueStatement.Where(documentValueTableName+"."+namespaceIdField+"=", namespaceDataMap.NamespaceId)
		documentValueResultDataResult := queryDocumentValueStatement.Find(&documentValueResultDataLists)
		if documentValueResultDataResult.Error != nil {
			return resultMap, documentValueResultDataResult.Error
		}
		var result []GlobalDocumentLanguage
		for _, documentValueResultData := range documentValueResultDataLists {
			var tableGlobalDocumentLanguageResult = GlobalDocumentLanguage{
				Id:                 documentValueResultData.DocumentId,
				CountryIso:         documentValueResultData.CountryIso,
				DocumentCode:       documentValueResultData.DocumentCode,
				DocumentValue:      documentValueResultData.DocumentValue,
				LastUpdateDocument: documentValueResultData.LastUpdateDocument,
				CreateTime:         convert.ToString(documentValueResultData.CreateTime),
			}
			result = append(result, tableGlobalDocumentLanguageResult)
		}
		documentResult.Documents = result
		resultList = append(resultList, documentResult)

		for _, value := range resultList {
			if len(value.Documents) > 0 {
				for _, value := range value.Documents {
					if resultMap[value.CountryIso] == nil {
						resultMap[value.CountryIso] = make(map[string]string)
					}
					countryIsoMap := resultMap[value.CountryIso]
					countryIsoMap[value.DocumentCode] = value.DocumentValue
				}
			}
		}
	}
	return resultMap, nil
}
