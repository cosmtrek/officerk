package services

import (
	"github.com/jinzhu/gorm"

	"github.com/cosmtrek/officerk/models"
)

// CreateNode creates node
func CreateNode(db *gorm.DB, n *models.Node) error {
	return db.Create(n).Error
}

// GetNodes fetch a list of nodes
func GetNodes(db *gorm.DB) ([]*models.Node, error) {
	var err error
	nodes := make([]*models.Node, 0)
	err = db.Where("deleted_at IS NULL").Find(&nodes).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nodes, nil
		}
		return nil, err
	}
	return nodes, nil
}
