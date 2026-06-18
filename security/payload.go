package security

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrorInvalidToken = errors.New("Token invalido")
	ErrorExpiredToken = errors.New("Token expirado")
)

type Payload struct {
	IDUsuario uuid.UUID `json:"idUsuario"`
	IDRol     int       `json:"idRol"`
	Nombre    string    `json:"nombre"`
	Email     string    `json:"email"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(idRol int, nombre string, email string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		IDUsuario: tokenID,
		IDRol:     idRol,
		Nombre:    nombre,
		Email:     email,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrorExpiredToken
	}
	return nil
}
