package repository

import (
	"errors"
	"ylem_api/helpers"
	"ylem_api/model/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OauthTokenRepository interface {
	FindByAccessToken(token string) (*entity.OauthToken, error)
	FindByRefreshToken(token string) (*entity.OauthToken, error)
	FindByAccessTokenAndClient(token string, clientUuid uuid.UUID) (*entity.OauthToken, error)
	FindByRefreshTokenAndClient(token string, clientUuid uuid.UUID) (*entity.OauthToken, error)
}

type GormOauthTokenRepository struct {
	db *gorm.DB
}

func (r *GormOauthTokenRepository) FindByAccessToken(accessToken string) (*entity.OauthToken, error) {
	token := &entity.OauthToken{}
	result := r.db.Preload("OauthClient").Where("access_token = ?", accessToken).Take(token)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}

	return token, nil
}

func (r *GormOauthTokenRepository) FindByRefreshToken(refreshToken string) (*entity.OauthToken, error) {
	token := &entity.OauthToken{}
	result := r.db.Preload("OauthClient").Where("refresh_token = ?", refreshToken).Take(token)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}

	return token, nil
}

func (r *GormOauthTokenRepository) FindByAccessTokenAndClient(accessToken string, clientUuid uuid.UUID) (*entity.OauthToken, error) {
	token := &entity.OauthToken{}
	result := r.db.Preload("OauthClient", "uuid = ?", clientUuid).Where("access_token = ?", accessToken).Take(token)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}

	return token, nil
}

func (r *GormOauthTokenRepository) FindByRefreshTokenAndClient(refreshToken string, clientUuid uuid.UUID) (*entity.OauthToken, error) {
	token := &entity.OauthToken{}
	result := r.db.Preload("OauthClient", "uuid = ?", clientUuid).Where("refresh_token = ?", refreshToken).Take(token)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}

	return token, nil
}

func NewOauthTokenRepository() OauthTokenRepository {
	return &GormOauthTokenRepository{
		db: helpers.GormInstance(),
	}
}
