package command

import (
	"crypto/rand"
	"encoding/hex"
	"ylem_api/model/entity"
	"ylem_api/model/persister"

	"github.com/google/uuid"
)

type CreateOauthClientCommand struct {
	OrganizationUuid uuid.UUID `json:"organization_uuid"`
	UserUuid         uuid.UUID `json:"user_uuid"`
	Name             string    `json:"name"`
}

type CreateOauthClientHandler struct {
	persister persister.EntityPersister
}

func (h *CreateOauthClientHandler) Handle(c CreateOauthClientCommand) (uuid.UUID, error) {
	s, err := generateSecureToken(32)
	if err != nil {
		return uuid.Nil, err
	}
	cl := &entity.OauthClient{
		Uuid:             uuid.New(),
		UserUuid:         c.UserUuid,
		OrganizationUuid: c.OrganizationUuid,
		Name:             c.Name,
		Secret:           s,
	}

	return cl.Uuid, h.persister.SaveOauthClient(cl)
}

func generateSecureToken(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func NewCreateOauthClientHandler() *CreateOauthClientHandler {
	h := &CreateOauthClientHandler{
		persister: persister.Instance(),
	}

	return h
}
