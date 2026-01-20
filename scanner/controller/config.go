package controller

import (
	"log"

	"github.com/Yoak3n/troll/scanner/model"
)

func (d *Database) QueryConfigurations() []model.ConfigurationTable {
	confs := make([]model.ConfigurationTable, 0)
	d.db.Find(&confs)
	return confs
}

func (d *Database) UpdateConfiguration(confs []model.ConfigurationTable) {
	d.db.Save(confs).Omit("type")
}

func (d *Database) DeleteConfiguration(ids []uint) {
	if len(ids) == 0 {
		return
	}
	log.Println("delete", ids)
	c := &model.ConfigurationTable{}
	d.db.Delete(c, ids)
}
