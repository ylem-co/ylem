package repository

import (
	"errors"
	"ylem_api/helpers"
	"ylem_api/model/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OauthClientRepository interface {
	FindById(int64) (*entity.OauthClient, error)
	FindByUuid(uuid.UUID) (*entity.OauthClient, error)
	FindAllByUserUuid(uuid.UUID) ([]*entity.OauthClient, error)
}

type GormClientRepository struct {
	db *gorm.DB
}

func (r *GormClientRepository) FindById(id int64) (*entity.OauthClient, error) {
	client := &entity.OauthClient{}
	result := r.db.First(client, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}

	return client, nil
}

func (r *GormClientRepository) FindByUuid(uid uuid.UUID) (*entity.OauthClient, error) {
	client := &entity.OauthClient{}
	result := r.db.Where("uuid = ? AND deleted_at IS NULL", uid).Take(client)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}

	return client, nil
}

func (r *GormClientRepository) FindAllByUserUuid(uid uuid.UUID) ([]*entity.OauthClient, error) {
	var clients []*entity.OauthClient
	result := r.db.Where("user_uuid = ? AND deleted_at IS NULL", uid).Find(&clients)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}

	return clients, nil
}

func NewOauthClientRepository() OauthClientRepository {
	return &GormClientRepository{
		db: helpers.GormInstance(),
	}
}
