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

// GetNode gets node
func GetNode(db *gorm.DB, id string, n *models.Node) error {
	return db.Where("deleted_at IS NULL").Where("id = ?", id).First(n).Error
}

// UpdateNode updates node
func UpdateNode(db *gorm.DB, n *models.Node, data *models.Node) error {
	return db.Model(n).Updates(data).Error
}

// DeleteNode deletes node
func DeleteNode(db *gorm.DB, n *models.Node) error {
	return db.Delete(n).Error
}
