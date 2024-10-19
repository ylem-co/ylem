package entities

import (
	regexp2 "regexp"
	"time"
)

type Sms struct {
	Id          int64       `json:"-"`
	Integration Integration `json:"integration"`
	Code        string      `json:"-"`
	IsConfirmed bool        `json:"is_confirmed"`
	RequestedAt time.Time   `json:"requested_at"`
}

func (s Sms) CanResendSms() bool {
	return !s.IsConfirmed
}

const IntegrationTypeSms = "sms"

func IsMobilePhoneValid(Number string) bool {
	return regexp2.
		MustCompile(`^((\+[0-9]{1,3})|(\+?\([0-9]{1,3}\)))[\s-]?(?:\(0?[0-9]{1,5}\)|[0-9]{1,5})[-\s]?[0-9][\d\s-]{5,7}\s?(?:x[\d-]{0,4})?$`).
		MatchString(Number)
}
