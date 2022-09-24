package database

import (
	"log"

	"backend.cs3219.comp.nus.edu.sg/model"
)

//go:generate mockgen -destination=../mocks/mock_database_card_adapter.go -build_flags=-mod=mod -package=mocks backend.cs3219.comp.nus.edu.sg/database DatabaseCardAdapter
type DatabaseCardAdapter interface {
	CreateCard(card *model.Card) (*model.Card, error)
	EditCard(card *model.Card) error
	DeleteCard(id int) error
	GetCard(id int) (*model.Card, error)
	GetCardByUniqueId(uniqueId string) (*model.Card, error)
	GetAllCards() ([]*model.Card, error)
}

type databaseCardAdapter struct {
	dbAdapter DatabaseAdapter[model.Card]
}

func NewDatabaseCardAdapter(connector *DatabaseConnection) DatabaseCardAdapter {
	return &databaseCardAdapter{
		dbAdapter: newDatabaseAdapter[model.Card](connector),
	}
}

func (adapter *databaseCardAdapter) CreateCard(card *model.Card) (*model.Card, error) {
	result, err := adapter.dbAdapter.QuerySingle(
		"INSERT INTO cards (card_unique_id, card_pokemon, card_image) VALUES(?, ?, ?) RETURNING card_id;",
		card.UniqueId,
		card.Pokemon,
		card.ImageUrl,
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	cardDuplicated := *card
	cardDuplicated.Id = result.Id
	return &cardDuplicated, nil
}

func (adapter *databaseCardAdapter) EditCard(card *model.Card) error {
	return adapter.dbAdapter.Execute(
		"UPDATE cards SET card_unique_id=?, card_pokemon=?, card_image=? WHERE card_id=?;",
		card.UniqueId,
		card.Pokemon,
		card.ImageUrl,
		card.Id,
	)
}

func (adapter *databaseCardAdapter) DeleteCard(id int) error {
	return adapter.dbAdapter.Execute(
		"DELETE FROM cards WHERE card_id=?",
		id,
	)
}

func (adapter *databaseCardAdapter) GetCard(id int) (*model.Card, error) {
	result, err := adapter.dbAdapter.QuerySingle(
		"SELECT * FROM cards WHERE card_id=?",
		id,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (adapter *databaseCardAdapter) GetCardByUniqueId(uniqueId string) (*model.Card, error) {
	result, err := adapter.dbAdapter.QuerySingle(
		"SELECT * FROM cards WHERE card_unique_id=?",
		uniqueId,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (adapter *databaseCardAdapter) GetAllCards() ([]*model.Card, error) {
	results, err := adapter.dbAdapter.QueryMany("SELECT * FROM cards")
	if err != nil {
		return nil, err
	}
	return results, nil
}
