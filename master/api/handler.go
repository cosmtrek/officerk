package api

import (
	"github.com/jinzhu/gorm"

	"github.com/cosmtrek/officerk/master/property"
)

var db *gorm.DB

// Handler ...
type Handler struct {
	runtime *property.Runtime
}

// NewHandler ...
func NewHandler(d *gorm.DB, runtime *property.Runtime) *Handler {
	// global db handler
	db = d
	return &Handler{
		runtime: runtime,
	}
}
