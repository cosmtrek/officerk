package node

import (
	"github.com/jinzhu/gorm"
)

type xMigrate interface {
	migrate(db *gorm.DB) error
}

var tables = []interface{}{
	&Job{},
	&Task{},
	&JobLog{},
	&TaskLog{},
}

func (r *Controller) autoMigrate() error {
	return r.db.AutoMigrate(tables...).Error
}
