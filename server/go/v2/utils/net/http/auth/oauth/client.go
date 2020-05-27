package oauth

import (
	"github.com/jinzhu/gorm"
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/models"
)

// NewClientStore create client store
func NewClientStore(db *gorm.DB) *ClientStore {
	return (*ClientStore)(db)
}

// ClientStore client information store
type ClientStore gorm.DB

// GetByID according to the ID for the client information
func (cs *ClientStore) GetByID(id string) (oauth2.ClientInfo, error) {
	db := (*gorm.DB)(cs)
	var client models.Client
	if err := db.Table("oauth_client").Find(&client, id).Error; err != nil {
		return nil, err
	}
	return &client, nil
}

// Set set client information
func (cs *ClientStore) Set(cli oauth2.ClientInfo) (err error) {
	db := (*gorm.DB)(cs)
	db.Table("oauth_client").Create(cli)
	return
}
