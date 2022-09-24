package auth

import (
	"log"

	"backend.cs3219.comp.nus.edu.sg/database"
)

//go:generate mockgen -destination=../mocks/mock_token_authenticator.go -build_flags=-mod=mod -package=mocks backend.cs3219.comp.nus.edu.sg/auth TokenAuthenticator
type TokenAuthenticator interface {
	IsValidToken(token string) bool
}

type tokenAuthenticator struct {
	tokenAdapter database.DatabaseApiTokenAdapter
}

func NewTokenAuthenticator(db *database.DatabaseConnection) TokenAuthenticator {
	return &tokenAuthenticator{
		tokenAdapter: database.NewDatabaseApiTokenAdapter(db),
	}
}

func (authenticator *tokenAuthenticator) IsValidToken(token string) bool {
	isValid, err := authenticator.tokenAdapter.IsValidToken(token)
	if err != nil {
		log.Println(err)
		return false
	}
	return isValid
}
