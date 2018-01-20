package api

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/jinzhu/gorm"

	"github.com/cosmtrek/officerk/master/property"
	"github.com/cosmtrek/officerk/models"
	"github.com/cosmtrek/officerk/utils/api"
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

// GetEnum ...
func (h *Handler) GetEnum(w http.ResponseWriter, r *http.Request) {
	render.Render(w, r, api.OK(models.Enum))
}
