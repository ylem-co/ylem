package command

import (
	"ylem_api/model/persister"

	"github.com/google/uuid"
)

type DeleteOauthClientCommand struct {
	OrganizationUuid uuid.UUID `json:"organization_uuid"`
	UserUuid         uuid.UUID `json:"user_uuid"`
}

type DeleteOauthClientHandler struct {
	persister persister.EntityPersister
}

func (h *DeleteOauthClientHandler) Handle(uid string) (bool, error) {
	err := h.persister.DeleteOauthClientByUuid(uid)
	if err != nil {
		return false, err
	}

	err = h.persister.DeleteOauthTokensByClientUuid(uid)
	if err != nil {
		return false, err
	}

	return true, nil
}

func NewDeleteOauthClientHandler() *DeleteOauthClientHandler {
	h := &DeleteOauthClientHandler{
		persister: persister.Instance(),
	}

	return h
}
