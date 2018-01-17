package utils

import (
	"github.com/jinzhu/gorm"
)

// OpenDB connects database
func OpenDB(c *Config) (*gorm.DB, error) {
	return gorm.Open("mysql", c.MySQL())
}
