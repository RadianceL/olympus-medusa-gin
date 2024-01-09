package model

import (
	"container/list"
	"errors"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"olympus-medusa/common"
	"olympus-medusa/common/convert"
	"olympus-medusa/common/language"
	"olympus-medusa/config"
	"olympus-medusa/data/data"
	"olympus-medusa/data/request"
	"strconv"
	"strings"
	"time"
)

type IDocumentModel interface {
	//ImportDocument(namespaceRequest multipart.File) (data.ExportGlobalDocument, error)
	//
	//getImportDataList(rows [][]string, isSuccess bool) (data.ExportGlobalDocument, error)

	CreateDocument(namespaceRequest *request.GlobalDocumentRequest) (int64, error)

	SearchDocumentValue(globalDocumentRequest *request.GlobalDocumentRequest) (arr []interface{})

	QueryDocument(globalDocumentRequest *request.GlobalDocumentRequest) ([]data.TableGlobalDocumentExcel, error)

	SearchDocumentByNamespaceId(globalDocumentRequest *request.GlobalDocumentRequest) (data.TableGlobalDocumentPage, error)

	UpdateDocumentByDocumentId(namespaceRequest *request.GlobalDocumentRequest) (int64, error)

	DeleteDocumentByDocumentId(namespaceRequest *request.GlobalDocumentRequest) (int64, error)

	SearchDocumentCode(globalDocumentRequest *request.GlobalDocumentRequest) (*data.TableGlobalDocument, error)

	SearchDocumentById(globalDocumentRequest *request.GlobalDocumentRequest) (data.TableGlobalDocument, error)

	SearchDocumentByCountryIso(globalDocumentIsoQueryRequest *request.GlobalDocumentIsoQueryRequest) (map[string]string, error)

	SearchApplicationByCountryIso(globalDocumentIsoQueryRequest *request.GlobalDocumentIsoQueryRequest) (map[string]map[string]string, error)

	SearchNamespaceById(applicationId int, appNamespaceId int) (data.TableApplicationNamespace, error)
}

// DocumentModel is application model structure.
type DocumentModel struct {
	logger *logrus.Logger
	db     *gorm.DB
}

func NewDocumentModel() IDocumentModel {
	return DocumentModel{db: common.GetDB(), logger: config.GetLogger()}
}

func (documentModel DocumentModel) CreateDocument(globalDocumentRequest *request.GlobalDocumentRequest) (int64, error) {
	tx := documentModel.db.Begin()
	documentCode := data.ApplicationGlobalizationDocumentCode{
		ApplicationID:        globalDocumentRequest.ApplicationId,
		NamespaceID:          globalDocumentRequest.NamespaceId,
		DocumentCode:         globalDocumentRequest.DocumentCode,
		DocumentDesc:         globalDocumentRequest.DocumentDesc,
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
		documentValue := data.TableGlobalDocumentValue{
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

func (documentModel DocumentModel) SearchDocumentValue(globalDocumentRequest *request.GlobalDocumentRequest) (arr []interface{}) {
	statementValue := documentModel.Table(documentValueTableName).Select("*")
	statementValue.Where("document_value", "LIKE", "%"+globalDocumentRequest.DocumentDesc+"%")
	documentValueMaps, _ := statementValue.All()
	var valueArr []interface{}
	for _, document := range documentValueMaps {
		var documentResult data.TableGlobalDocument
		_ = mapstructure.Decode(document, &documentResult)
		valueArr = append(valueArr, documentResult.DocumentId)
	}
	return valueArr
}

func (documentModel DocumentModel) QueryDocument(globalDocumentRequest *request.GlobalDocumentRequest) ([]data.TableGlobalDocumentExcel, error) {
	statement := documentModel.Table(documentValueTableName).Select("*")
	statement.LeftJoin(documentTableName, documentTableName+"."+documentIdField, "=", documentValueTableName+"."+documentIdField)
	statement.LeftJoin(applicationModelTableName, applicationModelTableName+"."+id, "=", documentTableName+"."+applicationIdField)
	statement.LeftJoin(namespaceModelTableName, namespaceModelTableName+"."+namespaceIdField, "=", documentTableName+"."+namespaceIdField)
	statement.Where(documentTableName+"."+deleteFlagField, "=", 0)
	if globalDocumentRequest.NamespaceId != 0 {
		statement.Where(documentTableName+"."+namespaceIdField, "=", globalDocumentRequest.NamespaceId)
	}
	if globalDocumentRequest.ApplicationId != 0 {
		statement.Where(documentTableName+"."+applicationIdField, "=", globalDocumentRequest.ApplicationId)
	}
	if len(globalDocumentRequest.DocumentIds) > 0 {
		statement.WhereIn(documentTableName+"."+documentIdField, globalDocumentRequest.DocumentIds)
	}
	if globalDocumentRequest.DocumentCode != "" {
		statement.Where(documentCodeField, "LIKE", "%"+globalDocumentRequest.DocumentCode+"%")
	}
	if globalDocumentRequest.DocumentDesc != "" {
		var arr = documentModel.SearchDocumentValue(globalDocumentRequest)
		if arr == nil || len(arr) <= 0 {
			arr = make([]interface{}, 1)
			arr[0] = 0
		}
		statement.WhereIn(documentTableName+"."+documentIdField, arr)
	}

	documentMaps, err := statement.All()

	var resultData []data.TableGlobalDocumentExcel
	if err != nil {
		return resultData, err
	}
	if len(documentMaps) <= 0 {
		return resultData, nil
	}
	for _, document := range documentMaps {
		var documentResult data.TableGlobalDocumentExcel
		_ = mapstructure.Decode(document, &documentResult)
		resultData = append(resultData, documentResult)
	}
	return resultData, nil
}

func (documentModel DocumentModel) SearchDocumentByNamespaceId(globalDocumentRequest *request.GlobalDocumentRequest) (data.TableGlobalDocumentPage, error) {
	statement := documentModel.Table(documentTableName).Select("*")
	statement.LeftJoin(applicationModelTableName, applicationModelTableName+"."+id, "=", documentTableName+"."+applicationIdField)
	statement.LeftJoin(namespaceModelTableName, namespaceModelTableName+"."+namespaceIdField, "=", documentTableName+"."+namespaceIdField)
	statement.Where(documentTableName+"."+deleteFlagField, "=", 0)
	if globalDocumentRequest.NamespaceId != 0 {
		statement.Where(documentTableName+"."+namespaceIdField, "=", globalDocumentRequest.NamespaceId)
	}
	if globalDocumentRequest.ApplicationId != 0 {
		statement.Where(documentTableName+"."+applicationIdField, "=", globalDocumentRequest.ApplicationId)
	}
	if globalDocumentRequest.DocumentCode != "" {
		statement.Where(documentCodeField, "LIKE", "%"+globalDocumentRequest.DocumentCode+"%")
	}
	var arr = make([]interface{}, 0)
	if globalDocumentRequest.DocumentDesc != "" {
		arr = documentModel.SearchDocumentValue(globalDocumentRequest)
		if arr == nil || len(arr) <= 0 {
			arr = make([]interface{}, 1)
			arr[0] = 0
		}
		statement.WhereIn(documentIdField, arr)
	}
	if globalDocumentRequest.PageIndex != 0 && globalDocumentRequest.PageSize != 0 {
		statement.Skip((globalDocumentRequest.PageIndex - 1) * globalDocumentRequest.PageSize)
		statement.Take(globalDocumentRequest.PageSize)
	}

	documentMaps, err := statement.All()
	statementCount := documentModel.Table(documentTableName)
	statementCount.Where(documentTableName+"."+deleteFlagField, "=", 0)
	if globalDocumentRequest.NamespaceId != 0 {
		statementCount.Where(documentTableName+"."+namespaceIdField, "=", globalDocumentRequest.NamespaceId)
	}
	if globalDocumentRequest.ApplicationId != 0 {
		statementCount.Where(documentTableName+"."+applicationIdField, "=", globalDocumentRequest.ApplicationId)
	}
	if globalDocumentRequest.DocumentCode != "" {
		statementCount.Where(documentCodeField, "LIKE", "%"+globalDocumentRequest.DocumentCode+"%")
	}
	if arr != nil && len(arr) > 0 {
		statementCount.WhereIn(documentIdField, arr)
	}
	count, _ := statementCount.Count()
	var resultData data.TableGlobalDocumentPage
	if err != nil {
		resultData.TotalSize = count
		return resultData, err
	}
	if len(documentMaps) <= 0 {
		resultData.TotalSize = count
		resultData.GlobalDocument = make([]data.TableGlobalDocument, 0)
		return resultData, nil
	}

	var resultList []data.TableGlobalDocument
	for _, document := range documentMaps {
		var documentResult data.TableGlobalDocument
		_ = mapstructure.Decode(document, &documentResult)

		queryDocumentValueStatement := documentModel.Table(documentValueTableName).Select("*")
		queryDocumentValueStatement.Where(documentIdField, "=", documentResult.DocumentId)
		documentValueResultDataMaps, documentValueErr := queryDocumentValueStatement.All()
		if documentValueErr != nil {
			resultData.TotalSize = count
			return resultData, nil
		}
		var result []data.TableGlobalDocumentLanguage
		for _, documentValueResultData := range documentValueResultDataMaps {
			var tableGlobalDocumentLanguageOutputResult data.TableGlobalDocumentLanguage
			_ = mapstructure.Decode(documentValueResultData, &tableGlobalDocumentLanguageOutputResult)
			result = append(result, tableGlobalDocumentLanguageOutputResult)
		}
		documentResult.Documents = result
		resultList = append(resultList, documentResult)
	}
	resultData.TotalSize = count
	resultData.GlobalDocument = resultList
	return resultData, nil
}

func (documentModel DocumentModel) UpdateDocumentByDocumentId(namespaceRequest *request.GlobalDocumentRequest) (int64, error) {
	tx := documentModel.db.Begin()
	if namespaceRequest.DocumentDesc != "" {
		_, err := documentModel.Table(documentTableName).
			WithTx(tx).
			Where(documentIdField, "=", namespaceRequest.DocumentId).
			Update(dialect.H{
				documentDescField: namespaceRequest.DocumentDesc,
			})
		if err != nil {
			if err.Error() != "no affect row" {
				_ = tx.Rollback()
				return 0, err
			}
		}
	}
	documents := namespaceRequest.Documents
	for _, document := range documents {
		languageCountry := language.FindLanguage(document.CountryIso)
		if languageCountry == nil {
			_ = tx.Rollback()
			return 0, errors.New("未识别的国家编码，请检查后重试")
		}
		queryDocumentValueStatement := documentModel.Table(documentValueTableName).Select("*")
		queryDocumentValueStatement.Where(documentIdField, "=", namespaceRequest.DocumentId)
		queryDocumentValueStatement.Where(countryIsoField, "=", document.CountryIso)
		documentValueResultDataMaps, documentValueErr := queryDocumentValueStatement.All()
		if documentValueErr != nil {
			return 0, errors.New("更新多语言文案语言编码查重异常，请稍后重试")
		}
		if len(documentValueResultDataMaps) <= 0 {
			languageCountry := language.FindLanguage(document.CountryIso)
			if languageCountry == nil {
				_ = tx.Rollback()
				return 0, errors.New("未识别的国家编码，请检查后重试")
			}
			_, err := documentModel.Table(documentValueTableName).
				WithTx(tx).
				Insert(dialect.H{
					documentIdField:    namespaceRequest.DocumentId,
					namespaceIdField:   namespaceRequest.NamespaceId,
					countryIsoField:    document.CountryIso,
					countryNameField:   languageCountry.CountryName,
					documentValueField: document.DocumentValue,
				})
			if err != nil {
				_ = tx.Rollback()
				return 0, err
			}
		} else {
			var documentResultDataMaps map[string]interface{}
			if document.DocumentId == 0 {
				queryDocumentValueStatement := documentModel.Table(documentValueTableName).Select("*")
				queryDocumentValueStatement.Where(documentIdField, "=", namespaceRequest.DocumentId)
				queryDocumentValueStatement.Where(namespaceIdField, "=", namespaceRequest.NamespaceId)
				queryDocumentValueStatement.Where(countryIsoField, "=", document.CountryIso)

				documentResultDataMaps, documentValueErr = queryDocumentValueStatement.First()
				if documentValueErr != nil {
					_ = tx.Rollback()
					return 0, documentValueErr
				}

			} else {
				queryDocumentValueStatement := documentModel.Table(documentValueTableName).Select("*")
				queryDocumentValueStatement.Where(id, "=", document.DocumentId)
				documentResultDataMaps, documentValueErr = queryDocumentValueStatement.First()
				if documentValueErr != nil {
					_ = tx.Rollback()
					return 0, documentValueErr
				}
			}
			var tableGlobalDocumentLanguageOutputResult data.TableGlobalDocumentLanguage
			_ = mapstructure.Decode(documentResultDataMaps, &tableGlobalDocumentLanguageOutputResult)
			_, err := documentModel.Table(documentValueTableName).
				WithTx(tx).
				Where(id, "=", tableGlobalDocumentLanguageOutputResult.Id).
				Update(dialect.H{
					documentValueField:      document.DocumentValue,
					lastUpdateDocumentField: tableGlobalDocumentLanguageOutputResult.DocumentValue,
				})
			if err != nil {
				if err.Error() != "no affect row" {
					_ = tx.Rollback()
					return 0, err
				}
			}
		}
	}
	commitError := tx.Commit()
	if commitError != nil {
		_ = tx.Rollback()
		return 0, tx.Error
	}
	return 1, nil
}

func (documentModel DocumentModel) DeleteDocumentByDocumentId(namespaceRequest *request.GlobalDocumentRequest) (int64, error) {
	nano := time.Now().Unix()
	tx := documentModel.db.Begin()
	_, err := documentModel.Table(documentTableName).
		WithTx(tx).
		Where(documentIdField, "=", namespaceRequest.DocumentId).
		Update(dialect.H{
			deleteFlagField:   1,
			remarksField:      namespaceRequest.Remarks,
			documentCodeField: namespaceRequest.DocumentCode + "_@delete_" + strconv.FormatInt(nano, 10),
			deleteTimeField:   time.Now(),
		})
	if err != nil {
		if err.Error() != "no affect row" {
			_ = tx.Rollback()
			return 0, err
		}
	}
	commitError := tx.Commit()
	if commitError != nil {
		_ = tx.Rollback()
		return 0, tx.Error
	}
	return 1, nil
}

func (documentModel DocumentModel) SearchDocumentCode(globalDocumentRequest *request.GlobalDocumentRequest) (*data.TableGlobalDocument, error) {
	statement := documentModel.Table(documentTableName).Select("*")
	if globalDocumentRequest.Id != 0 {
		statement.Where(documentIdField, "=", globalDocumentRequest.Id)
	} else if globalDocumentRequest.DocumentCode != "" {
		statement.Where(namespaceIdField, "=", globalDocumentRequest.NamespaceId)
		statement.Where(documentCodeField, "=", globalDocumentRequest.DocumentCode)
	} else {
		return nil, errors.New("文案ID或文案编码与命名空间ID必传其中一个")
	}

	resultData, err := statement.All()
	if err != nil {
		return nil, err
	}
	if len(resultData) <= 0 {
		return nil, nil
	}
	var outputResult data.TableGlobalDocument
	_ = mapstructure.Decode(resultData[0], &outputResult)
	return &outputResult, nil
}

func (documentModel DocumentModel) SearchDocumentById(globalDocumentRequest *request.GlobalDocumentRequest) (data.TableGlobalDocument, error) {
	statement := documentModel.Table(documentTableName).Select("*")
	if globalDocumentRequest.Id != 0 {
		statement.Where(documentIdField, "=", globalDocumentRequest.Id)
	} else if globalDocumentRequest.DocumentCode != "" {
		statement.Where(namespaceIdField, "=", globalDocumentRequest.NamespaceId)
		statement.Where(documentCodeField, "=", globalDocumentRequest.DocumentCode)
	} else {
		return data.TableGlobalDocument{}, errors.New("文案ID或文案编码与命名空间ID必传其中一个")
	}

	resultData, err := statement.All()
	if err != nil {
		return data.TableGlobalDocument{}, err
	}
	if len(resultData) <= 0 {
		return data.TableGlobalDocument{}, errors.New("未查询到编码信息，请确认后重试")
	}
	var outputResult data.TableGlobalDocument
	_ = mapstructure.Decode(resultData[0], &outputResult)

	queryDocumentValueStatement := documentModel.Table(documentValueTableName).Select("*")
	queryDocumentValueStatement.Where(documentIdField, "=", outputResult.DocumentId)
	documentValueResultDataMaps, documentValueErr := queryDocumentValueStatement.All()
	if documentValueErr != nil {
		return data.TableGlobalDocument{}, documentValueErr
	}
	var result []data.TableGlobalDocumentLanguage
	for _, documentValueResultData := range documentValueResultDataMaps {
		var tableGlobalDocumentLanguageOutputResult data.TableGlobalDocumentLanguage
		_ = mapstructure.Decode(documentValueResultData, &tableGlobalDocumentLanguageOutputResult)
		result = append(result, tableGlobalDocumentLanguageOutputResult)
	}
	outputResult.Documents = result
	return outputResult, nil
}

func (documentModel DocumentModel) SearchDocumentByCountryIso(globalDocumentIsoQueryRequest *request.GlobalDocumentIsoQueryRequest) (map[string]string, error) {
	resultMap := make(map[string]string)
	if len(globalDocumentIsoQueryRequest.NamespacePath) >= 5 {
		return resultMap, errors.New("单次查询不能超过4个命名空间")
	}
	namespaceStatement := documentModel.Table(namespaceModelTableName)
	namespaceStatement.WhereIn(namespacePath, globalDocumentIsoQueryRequest.NamespacePath)
	namespaceResultMap, err := namespaceStatement.All()
	if err != nil {
		return resultMap, err
	}
	if len(namespaceResultMap) >= 5 || len(namespaceResultMap) <= 0 {
		return resultMap, errors.New("命名空间数据异常，请检查配置后重试")
	}
	for _, value := range namespaceResultMap {
		var resultList []data.TableGlobalDocument
		var documentResult data.TableGlobalDocument

		queryDocumentCodeStatement := documentModel.Table(documentTableName).Select("*")
		queryDocumentCodeStatement.Where(documentTableName+"."+namespaceIdField, "=", value["NamespaceId"])
		queryDocumentCodeStatementResult, _ := queryDocumentCodeStatement.All()

		queryDocumentValueStatement := documentModel.Table(documentValueTableName).Select("*")
		queryDocumentValueStatement.Where(documentValueTableName+"."+namespaceIdField, "=", value["NamespaceId"])
		queryDocumentValueStatement.Where(countryIsoField, "=", globalDocumentIsoQueryRequest.CountryIso)
		queryDocumentValueStatement.LeftJoin(documentTableName, documentValueTableName+"."+documentIdField, "=", documentTableName+"."+documentIdField)
		documentValueResultDataMaps, documentValueErr := queryDocumentValueStatement.All()

		documentResultList := list.New()
		for _, documentResultMap := range documentValueResultDataMaps {
			documentResultList.PushBack(convert.ToString(documentResultMap["DocumentCode"]))
		}
		documentCodeList := list.New()
		for _, documentCode := range queryDocumentCodeStatementResult {
			documentCodeList.PushBack(convert.ToString(documentCode["DocumentCode"]))
		}
		for e := documentResultList.Front(); e != nil; e = e.Next() {
			for f := documentCodeList.Front(); f != nil; f = f.Next() {
				if strings.EqualFold(e.Value.(string), f.Value.(string)) {
					documentCodeList.Remove(f)
					break
				}
			}
		}
		if documentValueErr != nil {
			return resultMap, documentValueErr
		}
		if !strings.EqualFold(globalDocumentIsoQueryRequest.CountryIso, "CN") {
			queryDocumentValueEnStatement := documentModel.Table(documentValueTableName).Select("*")
			queryDocumentValueEnStatement.Where(documentValueTableName+"."+namespaceIdField, "=", value["NamespaceId"])
			queryDocumentValueEnStatement.Where(countryIsoField, "=", "EN")
			queryDocumentValueEnStatement.LeftJoin(documentTableName, documentValueTableName+"."+documentIdField, "=", documentTableName+"."+documentIdField)
			var arrays []interface{}
			for e := documentCodeList.Front(); e != nil; e = e.Next() {
				arrays = append(arrays, e.Value)
			}
			if len(arrays) != 0 {
				queryDocumentValueEnStatement.WhereIn(documentTableName+"."+documentCodeField, arrays)
				documentValueResultEnDataMaps, documentEnValueErr := queryDocumentValueEnStatement.All()
				if documentEnValueErr != nil {
					return resultMap, documentEnValueErr
				}
				var result []data.TableGlobalDocumentLanguage
				for _, documentValueResultData := range documentValueResultEnDataMaps {
					var tableGlobalDocumentLanguageOutputResult data.TableGlobalDocumentLanguage
					_ = mapstructure.Decode(documentValueResultData, &tableGlobalDocumentLanguageOutputResult)
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
		var result []data.TableGlobalDocumentLanguage
		for _, documentValueResultData := range documentValueResultDataMaps {
			var tableGlobalDocumentLanguageOutputResult data.TableGlobalDocumentLanguage
			_ = mapstructure.Decode(documentValueResultData, &tableGlobalDocumentLanguageOutputResult)
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

func (documentModel DocumentModel) SearchApplicationByCountryIso(globalDocumentIsoQueryRequest *request.GlobalDocumentIsoQueryRequest) (map[string]map[string]string, error) {
	resultMap := make(map[string]map[string]string)
	applicationStatement := documentModel.Table(applicationModelTableName)
	applicationStatement.Where(applicationPath, "=", globalDocumentIsoQueryRequest.ApplicationPath)
	applicationResultMap, err := applicationStatement.All()
	if err != nil {
		return resultMap, err
	}
	if len(applicationResultMap) > 1 || len(applicationResultMap) <= 0 {
		return resultMap, errors.New("应用空间数据异常，请检查配置后重试")
	}

	namespaceStatement := documentModel.Table(namespaceModelTableName)
	namespaceStatement.Where(applicationIdField, "=", applicationResultMap[0]["Id"])
	namespaceResultMap, err := namespaceStatement.All()
	if err != nil {
		return resultMap, err
	}
	if len(namespaceResultMap) <= 0 {
		return resultMap, errors.New("命名空间数据异常，请检查配置后重试")
	}

	for _, namespaceDataMap := range namespaceResultMap {
		var resultList []data.TableGlobalDocument
		var documentResult data.TableGlobalDocument
		queryDocumentValueStatement := documentModel.Table(documentValueTableName).Select("*")
		queryDocumentValueStatement.LeftJoin(documentTableName, documentValueTableName+"."+documentIdField, "=", documentTableName+"."+documentIdField)
		queryDocumentValueStatement.Where(documentValueTableName+"."+namespaceIdField, "=", namespaceDataMap["NamespaceId"])
		documentValueResultDataMaps, documentValueErr := queryDocumentValueStatement.All()
		if documentValueErr != nil {
			return resultMap, documentValueErr
		}
		var result []data.TableGlobalDocumentLanguage
		for _, documentValueResultData := range documentValueResultDataMaps {
			var tableGlobalDocumentLanguageOutputResult data.TableGlobalDocumentLanguage
			_ = mapstructure.Decode(documentValueResultData, &tableGlobalDocumentLanguageOutputResult)
			result = append(result, tableGlobalDocumentLanguageOutputResult)
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

func (documentModel DocumentModel) SearchOptionByNamespace(globalDocumentRequest *request.GlobalDocumentRequest) ([]data.ApplicationGlobalizationDocumentCode, error) {
	resultMap := make(map[string]string)
	applicationStatement := documentModel.db.
		Where("application_id = ? AND namespace_id = ?",
			globalDocumentRequest.ApplicationId, globalDocumentRequest.NamespaceId)
	if globalDocumentRequest.DocumentDesc != "" {
		applicationStatement.Where("document_desc LIKE ", "%"+convert.ToString(globalDocumentRequest.DocumentDesc)+"%")
	}
	if globalDocumentRequest.DocumentCode != "" {
		applicationStatement.Where("document_code LIKE ", "%"+convert.ToString(globalDocumentRequest.DocumentCode)+"%")
	}
	var applicationCodes []data.ApplicationGlobalizationDocumentCode
	if err := applicationStatement.Find(&applicationCodes).Error; err != nil {
		return applicationCodes, err
	}
	if len(applicationCodes) <= 0 {
		return applicationCodes, errors.New("应用空间数据异常，请检查配置后重试")
	}

	for _, namespaceDataMap := range applicationCodes {
		resultMap[convert.ToString(namespaceDataMap["DocumentCode"])] = convert.ToString(namespaceDataMap["DocumentDesc"])
	}
	return applicationCodes, nil
}

func (documentModel DocumentModel) SearchNamespaceById(applicationId int, appNamespaceId int) (data.TableApplicationNamespace, error) {
	var application data.TableApplicationNamespace
	if err := documentModel.db.Where("application_id = ? AND namespace_id = ?", applicationId, appNamespaceId).Find(&application).Error; err != nil {
		return data.TableApplicationNamespace{}, err
	}
	return application, nil
}
