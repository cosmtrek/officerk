package node

var tables = []interface{}{
	&Job{},
	&Task{},
	&JobLog{},
	&TaskLog{},
}

func (r *Controller) autoMigrate() error {
	var err error
	if err = r.db.AutoMigrate(tables...).Error; err != nil {
		return err
	}
	err = r.db.Model(&Task{}).AddUniqueIndex("task_name_and_job_id_unique", "name", "job_id").Error
	if err != nil {
		return err
	}
	return nil
}
