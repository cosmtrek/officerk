package node

var tables = []interface{}{
	&Job{},
	&Task{},
	&JobLog{},
	&TaskLog{},
}

func (r *Controller) autoMigrate() error {
	return r.db.AutoMigrate(tables...).Error
}
