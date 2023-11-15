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
