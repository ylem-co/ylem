package oauth

import (
	"context"
	"ylem_api/model/entity"
	"ylem_api/model/persister"
	"ylem_api/model/repository"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/google/uuid"
)

type GormOauthClientStore struct {
	repository repository.OauthClientRepository
	persister  persister.EntityPersister
}

func (s *GormOauthClientStore) GetByID(ctx context.Context, id string) (oauth2.ClientInfo, error) {
	return s.repository.FindByUuid(uuid.MustParse(id))
}

func (s *GormOauthClientStore) Create(ctx context.Context, client *entity.OauthClient) (oauth2.ClientInfo, error) {
	return client, s.persister.SaveOauthClient(client)
}

func NewOauthClientStore() oauth2.ClientStore {
	return &GormOauthClientStore{
		repository: repository.NewOauthClientRepository(),
		persister:  persister.Instance(),
	}
}
