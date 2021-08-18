package oauth

import (
	"gopkg.in/oauth2.v3"
)

type CustomClientStore struct {
	oauth2.ClientStore
}

type CustomClientInfo struct {
	oauth2.ClientInfo
	ID     string
	Secret string
	Domain string
	UserID string
}

func NewCustomStore() oauth2.ClientStore {
	s := new(CustomClientStore)
	return s
}

func (a *CustomClientStore) GetByID(id string) (oauth2.ClientInfo, error) {
	info, err := userQuery(id)
	return &info, err
}

func (a *CustomClientInfo) GetID() string {
	return a.ID
}

func (a *CustomClientInfo) GetSecret() string {
	return a.Secret
}

func (a *CustomClientInfo) GetDomain() string {
	return a.Domain
}

func (a *CustomClientInfo) GetUserID() string {
	return a.UserID
}
