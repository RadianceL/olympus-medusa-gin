package common

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/url"
	"time"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	driverName := viper.GetString("gorm.driverName")
	host := viper.GetString("gorm.mysql.host")
	port := viper.GetString("gorm.mysql.port")
	database := viper.GetString("gorm.mysql.database")
	username := viper.GetString("gorm.mysql.username")
	password := viper.GetString("gorm.mysql.password")
	charset := viper.GetString("gorm.mysql.charset")
	loc := viper.GetString("gorm.mysql.loc")
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=%s",
		username, password, host, port,
		database, charset, url.QueryEscape(loc))
	fmt.Print(driverName)
	db, err := gorm.Open(mysql.Open(args), &gorm.Config{})
	if err != nil {
		panic("fail to connect database, err: " + err.Error())
	}
	//db.AutoMigrate(&model.User{})
	//db.AutoMigrate(&model.Article{})
	//db.AutoMigrate(&model.InviteCode{})
	//db.AutoMigrate(&model.Tags{})
	//db.AutoMigrate(&model.TagMapArticle{})

	sqlDB, err := db.DB()
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(viper.GetInt("gorm.max_idle_conns"))

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(viper.GetInt("gorm.max_open_conns"))

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)
	DB = db
	return db
}

func GetDB() *gorm.DB {
	return DB
}
