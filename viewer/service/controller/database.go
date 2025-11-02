package controller

import (
	"sync"

	"github.com/Yoak3n/troll/viewer/service/database"
	"gorm.io/gorm"
)

type Database struct {
	db *gorm.DB
}

var DB *Database
var once sync.Once

func NewDatabase(dir string, name string) *Database {
	return &Database{
		db: database.InitDatabase(dir, name),
	}
}
func GlobalDatabase(args ...string) *Database {
	once.Do(func() {
		if len(args) == 2 {
			DB = NewDatabase(args[0], args[1])
		} else {
			panic("database init args length is not 2")
		}
	})
	return DB
}
