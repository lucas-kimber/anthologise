package api

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"

	"github.com/google/uuid"
)

func createUserID() string {
	return uuid.NewString()
}

func creatEditToken() (string, [32]byte) {

	b := make([]byte, 32)
	rand.Read(b)

	editToken := base64.RawURLEncoding.EncodeToString(b)
	tokenHash := sha256.Sum256([]byte(editToken))

	return editToken, tokenHash
}

func (s *server) verifyUser(userID, editToken string) (bool, error) {

	targetHash, err := s.store.GetTokenHash(userID)

	if err != nil {
		return false, err
	}

	givenHash := sha256.Sum256([]byte(editToken))

	return subtle.ConstantTimeCompare(targetHash[:], givenHash[:]) == 1, nil
}

func authMiddleware() {}
