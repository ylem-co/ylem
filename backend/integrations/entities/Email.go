package entities

import "time"

type Email struct {
	Id          int64       `json:"-"`
	Integration Integration `json:"integration"`
	Code        string      `json:"-"`
	IsConfirmed bool        `json:"is_confirmed"`
	RequestedAt time.Time   `json:"requested_at"`
}

const IntegrationTypeEmail = "email"

func (e Email) CanResendEmail() bool {
	// @todo what criteria?
	return !e.IsConfirmed
}
