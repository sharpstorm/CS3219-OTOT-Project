package database

import (
	"backend.cs3219.comp.nus.edu.sg/model"
)

//go:generate mockgen -destination=../mocks/mock_database_api_token_adapter.go -build_flags=-mod=mod -package=mocks backend.cs3219.comp.nus.edu.sg/database DatabaseApiTokenAdapter
type DatabaseApiTokenAdapter interface {
	CreateApiToken(token string) error
	DeleteApiToken(id int) error
	SetApiTokenState(id int, isActive bool) error
	IsValidToken(token string) (bool, error)
	GetAllApiTokens() ([]*model.ApiToken, error)
}

type databaseApiTokenAdapter struct {
	dbAdapter DatabaseAdapter[model.ApiToken]
}

func NewDatabaseApiTokenAdapter(connector *DatabaseConnection) DatabaseApiTokenAdapter {
	return &databaseApiTokenAdapter{
		dbAdapter: newDatabaseAdapter[model.ApiToken](connector),
	}
}

func (adapter *databaseApiTokenAdapter) CreateApiToken(token string) error {
	return adapter.dbAdapter.Execute(
		"INSERT INTO api_tokens (token, is_enabled, created_at) VALUES(?, true, NOW())",
		token,
	)
}

func (adapter *databaseApiTokenAdapter) SetApiTokenState(id int, isActive bool) error {
	return adapter.dbAdapter.Execute(
		"UPDATE api_tokens SET is_enabled=? WHERE token_id=?",
		isActive,
		id,
	)
}

func (adapter *databaseApiTokenAdapter) DeleteApiToken(id int) error {
	return adapter.dbAdapter.Execute(
		"DELETE FROM api_tokens WHERE token_id = ?;",
		id,
	)
}

func (adapter *databaseApiTokenAdapter) IsValidToken(token string) (bool, error) {
	row, err := adapter.dbAdapter.QuerySingle(
		"SELECT * FROM api_tokens WHERE token=? AND is_enabled = TRUE",
		token,
	)
	if err != nil {
		return false, err
	}

	return row != nil, nil
}

func (adapter *databaseApiTokenAdapter) GetAllApiTokens() ([]*model.ApiToken, error) {
	results, err := adapter.dbAdapter.QueryMany("SELECT * FROM api_tokens")
	if err != nil {
		return nil, err
	}
	return results, nil
}
