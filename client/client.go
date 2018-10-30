package client

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gopkg.in/oauth2.v3"
)

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	db.AutoMigrate(new(Client))
	return &Store{db: db}
}

func (cs *Store) GetByID(uid string) (oauth2.ClientInfo, error) {
	var client Client
	err := cs.db.Where("uid = ?", uid).First(&client).Error
	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (cs *Store) Set(cli *Client) error {
	return cs.db.Save(cli).Error
}

type Client struct {
	gorm.Model
	UID     string `json:"uid" gorm:"type:varchar(100);unique_index"`
	Secret  string `json:"secret" gorm:"type:varchar(100)"`
	Domain  string `json:"domain" gorm:"type:varchar(100)"`
	UserID  string `json:"user_id" gorm:"type:varchar(50)"`
	AppName string `json:"app_name" gorm:"type:varchar(50);unique"`
}

func (c *Client) TableName() string {
	return "oauth2_client"
}

func (c *Client) GetID() string {
	return c.UID
}

func (c *Client) GetSecret() string {
	return c.Secret
}

func (c *Client) GetDomain() string {
	return c.Domain
}

func (c *Client) GetUserID() string {
	return c.UserID
}
