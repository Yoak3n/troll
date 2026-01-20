package controller

import "github.com/Yoak3n/troll/scanner/model"

func (d *Database) QueryConfigurationCookie() (*model.ConfigurationTable, error) {
	c := &model.ConfigurationTable{
		Type: "cookie",
	}
	err := d.db.First(c).Error
	return c, err
}

func (d *Database) QueryConfigurationProxy() (*model.ConfigurationTable, error) {
	c := &model.ConfigurationTable{
		Type: "proxy",
	}
	err := d.db.First(c).Error
	return c, err
}

func (d *Database) QueryConfiguration() ([]model.ConfigurationTable, error) {
	confs := make([]model.ConfigurationTable, 0)
	err := d.db.Where("invalid = ?", false).Find(&confs).Error
	return confs, err
}

func (d *Database) UpdateConfigurationRecord(c *model.ConfigurationTable) error {
	return d.db.Save(c).Error
}
