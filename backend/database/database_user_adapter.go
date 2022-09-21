package database

import (
	"backend.cs3219.comp.nus.edu.sg/model"
)

type DatabaseUserAdapter interface {
	CreateUser(user *model.User) (*model.User, error)
	EditUser(user *model.User) error
	DeleteUser(id int) error
	GetUser(id int) (*model.User, error)
	GetAllUsers() ([]*model.User, error)
}

type databaseUserAdapter struct {
	dbAdapter DatabaseAdapter[model.User]
}

func NewDatabaseUserAdapter(connector *DatabaseConnection) DatabaseUserAdapter {
	return &databaseUserAdapter{
		dbAdapter: newDatabaseAdapter[model.User](connector),
	}
}

func (adapter *databaseUserAdapter) CreateUser(user *model.User) (*model.User, error) {
	result, err := adapter.dbAdapter.QuerySingle(
		"INSERT INTO users (user_name, user_password) VALUES(?, ?);",
		user.Username,
		user.Password,
	)
	if err != nil {
		return nil, err
	}
	userDuplicated := *user
	userDuplicated.Id = result.Id
	return &userDuplicated, nil
}

func (adapter *databaseUserAdapter) EditUser(user *model.User) error {
	return adapter.dbAdapter.Execute(
		"UPDATE users SET user_name=?, user_password=? WHERE user_id=?;",
		user.Username,
		user.Password,
		user.Id,
	)
}

func (adapter *databaseUserAdapter) DeleteUser(id int) error {
	return adapter.dbAdapter.Execute(
		"DELETE FROM users WHERE user_id=?",
		id,
	)
}

func (adapter *databaseUserAdapter) GetUser(id int) (*model.User, error) {
	result, err := adapter.dbAdapter.QuerySingle(
		"SELECT * FROM users WHERE user_id=?",
		id,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (adapter *databaseUserAdapter) GetAllUsers() ([]*model.User, error) {
	results, err := adapter.dbAdapter.QueryMany("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	return results, nil
}
