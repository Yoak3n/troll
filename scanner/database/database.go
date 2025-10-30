package database

import (
	"fmt"
	"path"
	"time"

	"github.com/Yoak3n/gulu/logger"
	"github.com/Yoak3n/gulu/util"
	"github.com/Yoak3n/troll/scanner/model"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func initSqlite(dir string, name string) *gorm.DB {
	_ = util.CreateDirNotExists(dir)
	dsn := fmt.Sprintf("%s.db", name)
	dsn = path.Join(dir, dsn)
	sdb, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Logger.Println(fmt.Sprintf("database connected err:%v", err))
	}
	if err != nil {
		logger.Logger.Println(fmt.Sprintf("database connected err:%v", err))
	}
	return sdb
}

func InitDatabase(path string, name string) *gorm.DB {
	db := initSqlite(path, name)
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&model.VideoTable{}, &model.UserTable{}, &model.ConfigurationTable{}, &model.CommentTable{})
	if err != nil {
		panic(err)
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)
	return db
}
