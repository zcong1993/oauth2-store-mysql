package token

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/models"
)

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	db.AutoMigrate(new(Token))
	return &Store{db: db}
}

type Token struct {
	gorm.Model
	ExpiredAt int64  `json:"expired_at" gorm:"index"`
	Code      string `json:"code" gorm:"index"`
	Access    string `json:"access" gorm:"index"`
	Refresh   string `json:"refresh" gorm:"index"`
	Data      string `json:"data"`
}

func (t *Token) TableName() string {
	return "oauth2_token"
}

func (s *Store) Create(info oauth2.TokenInfo) error {
	bf, _ := json.Marshal(info)
	item := &Token{
		Data: string(bf),
	}

	if code := info.GetCode(); code != "" {
		item.Code = code
		item.ExpiredAt = info.GetCodeCreateAt().Add(info.GetCodeExpiresIn()).Unix()
	} else {
		item.Access = info.GetAccess()
		item.ExpiredAt = info.GetAccessCreateAt().Add(info.GetAccessExpiresIn()).Unix()

		if refresh := info.GetRefresh(); refresh != "" {
			item.Refresh = info.GetRefresh()
			item.ExpiredAt = info.GetRefreshCreateAt().Add(info.GetRefreshExpiresIn()).Unix()
		}
	}

	return s.db.Create(item).Error
}

func toTokenInfo(data string) oauth2.TokenInfo {
	var tf models.Token
	json.Unmarshal([]byte(data), &tf)
	return &tf
}

func (s *Store) RemoveByCode(code string) error {
	return s.db.Delete(&Token{}, "code = ?", code).Error
}

func (s *Store) RemoveByAccess(access string) error {
	return s.db.Delete(&Token{}, "access = ?", access).Error
}

func (s *Store) RemoveByRefresh(refresh string) error {
	return s.db.Delete(&Token{}, "refresh = ?", refresh).Error
}

func (s *Store) GetByCode(code string) (oauth2.TokenInfo, error) {
	var info Token
	err := s.db.Where("code = ?", code).First(&info).Error
	if err != nil {
		return nil, err
	}
	return toTokenInfo(info.Data), nil
}

func (s *Store) GetByAccess(access string) (oauth2.TokenInfo, error) {
	var info Token
	err := s.db.Where("access = ?", access).First(&info).Error
	if err != nil {
		return nil, err
	}
	return toTokenInfo(info.Data), nil
}

func (s *Store) GetByRefresh(refresh string) (oauth2.TokenInfo, error) {
	var info Token
	err := s.db.Where("refresh = ?", refresh).First(&info).Error
	if err != nil {
		return nil, err
	}
	return toTokenInfo(info.Data), nil
}
