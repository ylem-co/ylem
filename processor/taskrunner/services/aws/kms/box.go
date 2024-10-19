package kms

import (
    "errors"
)

type SecretBox struct {
    Sealed         bool
    EncryptedValue []byte
    PlainValue     []byte
}

func (s *SecretBox) UnmarshalJSON(data []byte) error {
    s.Sealed = false
    s.PlainValue = data

    return nil
}

func (s *SecretBox) MarshalJSON() ([]byte, error) {
    if s.Sealed {
        return nil, errors.New("can't marshal a value as the box is sealed")
    }

    v := s.PlainValue
    v = append([]byte(`"`), v...)
    v = append(v, []byte(`"`)...)

    return v, nil
}

func (s *SecretBox) SetPlainValue(value []byte) *SecretBox {
    s.PlainValue = value
    s.EncryptedValue = nil
    s.Sealed = false

    return s
}

func (s *SecretBox) SetEncryptedValue(value []byte) *SecretBox {
    s.EncryptedValue = value

    return s
}

func (s *SecretBox) Seal() {
    s.PlainValue = nil
    s.Sealed = true
}

func (s *SecretBox) Open(value []byte) {
    s.PlainValue = value
    s.Sealed = false
}

func NewOpenSecretBox(value []byte) SecretBox {
    return SecretBox{
        Sealed: false,
        EncryptedValue: nil,
        PlainValue: value,
    }
}

func NewSealedSecretBox(value []byte) SecretBox {
    return SecretBox{
        Sealed: true,
        EncryptedValue: value,
        PlainValue: nil,
    }
}
