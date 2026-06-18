package security

import "time"

type Builder interface {
	CreateToken(idRol int, nombre string, email string, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}
