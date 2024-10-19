package persister

import (
	"ylem_api/helpers"
	"ylem_api/model/entity"

	"gorm.io/gorm"
)

type EntityPersister interface {
	SaveOauthClient(*entity.OauthClient) error
	SaveOauthToken(*entity.OauthToken) error
	DeleteOauthToken(id uint) error
	DeleteOauthClientByUuid(uuid string) error
	DeleteOauthTokensByClientUuid(uuid string) error
}

type GormEntityPersister struct {
	db *gorm.DB
}

func (p *GormEntityPersister) SaveOauthClient(s *entity.OauthClient) error {
	if s.ID > 0 {
		result := p.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(s)
		return result.Error
	}

	result := p.db.Create(s)
	return result.Error
}

func (p *GormEntityPersister) SaveOauthToken(t *entity.OauthToken) error {
	if t.ID > 0 {
		result := p.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(t)
		return result.Error
	}

	result := p.db.Create(t)
	return result.Error
}

func (p *GormEntityPersister) DeleteOauthToken(id uint) error {
	return p.db.Delete(&entity.OauthToken{}, id).Error
}

func (p *GormEntityPersister) DeleteOauthClientByUuid(uid string) error {
	return p.db.Where("uuid = ?", uid).Delete(&entity.OauthClient{}).Error
}

func (p *GormEntityPersister) DeleteOauthTokensByClientUuid(uid string) error {
	return p.db.Where("oauth_client_uuid = ?", uid).Delete(&entity.OauthToken{}).Error
}

func Instance() EntityPersister {
	return &GormEntityPersister{
		db: helpers.GormInstance(),
	}
}
