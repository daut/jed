package tokens

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"
)

type Token struct {
	PlainText string    `json:"token"`
	Hash      []byte    `json:"-"`
	AdminID   int32     `json:"-"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func GenerateToken(adminID int32, duration time.Duration) (*Token, error) {
	token := &Token{
		AdminID:   adminID,
		ExpiresAt: time.Now().Add(duration),
	}

	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.PlainText = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(token.PlainText))
	token.Hash = hash[:]

	return token, nil
}
