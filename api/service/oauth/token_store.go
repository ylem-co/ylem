package oauth

import (
	"context"
	"errors"
	"ylem_api/model/entity"
	"ylem_api/model/persister"
	"ylem_api/model/repository"
	"ylem_api/service/ylem_users"

	"github.com/go-oauth2/oauth2/v4"
	log "github.com/sirupsen/logrus"
)

type GormOauthTokenStore struct {
	tokenRepository  repository.OauthTokenRepository
	clientRepository repository.OauthClientRepository
	persister        persister.EntityPersister
	ylemUsersClient  *ylem_users.Client
}

// create and store the new token information
func (s *GormOauthTokenStore) Create(ctx context.Context, info oauth2.TokenInfo) error {
	t := entity.NewToken()
	t.SetAccess(info.GetAccess())
	t.SetClientID(info.GetClientID())
	t.SetUserID(info.GetUserID())
	t.SetScope(info.GetScope())
	t.SetCodeCreateAt(info.GetCodeCreateAt())
	t.SetAccessCreateAt(info.GetAccessCreateAt())
	t.SetAccessExpiresIn(info.GetAccessExpiresIn())
	t.SetRefresh(info.GetRefresh())
	t.SetRefreshCreateAt(info.GetRefreshCreateAt())
	t.SetRefreshExpiresIn(info.GetAccessExpiresIn())

	var err error
	oauthClient, err := s.clientRepository.FindByUuid(t.OauthClientUuid)
	if err != nil {
		log.Error(err)
		return err
	}

	t.InternalToken, err = s.ylemUsersClient.IssueJWT(oauthClient.OrganizationUuid, oauthClient.UserUuid)
	if err != nil {
		log.Error(err)
		return err
	}

	return s.persister.SaveOauthToken(t)
}

// delete the authorization code
func (s *GormOauthTokenStore) RemoveByCode(ctx context.Context, code string) error {
	return errors.New("not implemented")
}

// use the access token to delete the token information
func (s *GormOauthTokenStore) RemoveByAccess(ctx context.Context, access string) error {
	t, err := s.tokenRepository.FindByAccessToken(access)
	if err != nil {
		return err
	}

	return s.persister.DeleteOauthToken(t.ID)
}

// use the refresh token to delete the token information
func (s *GormOauthTokenStore) RemoveByRefresh(ctx context.Context, refresh string) error {
	t, err := s.tokenRepository.FindByRefreshToken(refresh)
	if err != nil {
		return err
	}

	return s.persister.DeleteOauthToken(t.ID)
}

// use the authorization code for token information data
func (s *GormOauthTokenStore) GetByCode(ctx context.Context, code string) (oauth2.TokenInfo, error) {
	return nil, errors.New("not implemented")
}

// use the access token for token information data
func (s *GormOauthTokenStore) GetByAccess(ctx context.Context, access string) (oauth2.TokenInfo, error) {
	t, err := s.tokenRepository.FindByAccessToken(access)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return t, nil
}

// use the refresh token for token information data
func (s *GormOauthTokenStore) GetByRefresh(ctx context.Context, refresh string) (oauth2.TokenInfo, error) {
	t, err := s.tokenRepository.FindByRefreshToken(refresh)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return t, nil
}

func NewOauthTokenStore() oauth2.TokenStore {
	return &GormOauthTokenStore{
		tokenRepository:  repository.NewOauthTokenRepository(),
		clientRepository: repository.NewOauthClientRepository(),
		ylemUsersClient:  ylem_users.NewClient(),
		persister:        persister.Instance(),
	}
}
