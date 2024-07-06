package token

type Maker interface {
	CreateToken(username int64) (string, error)
	VerifyToken(token string) (*Payload, error)
}
