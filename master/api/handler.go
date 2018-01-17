package api

import (
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

// Handler ...
type Handler struct{}

// NewHandler ...
func NewHandler(d *gorm.DB) *Handler {
	// global db handler
	db = d
	return &Handler{}
}
