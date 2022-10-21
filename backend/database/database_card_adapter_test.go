package database

import (
	"context"
	"testing"

	"backend.cs3219.comp.nus.edu.sg/model"
	"backend.cs3219.comp.nus.edu.sg/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CardAdapterTestSuite struct {
	suite.Suite
	conn       *DatabaseConnection
	ctx        context.Context
	seedModels []*model.Card
}

func (suite *CardAdapterTestSuite) SetupSuite() {
	config := util.LoadEnvVariables()
	conn, err := ConnectDatabase(
		config.DbUrl,
		config.DbUsername,
		config.DbPassword,
		config.DbName,
	)
	if err != nil || conn == nil {
		suite.T().Fatal("Failed to read db config")
	}
	suite.ctx = context.Background()
	suite.conn = conn

	suite.seedModels = []*model.Card{
		{
			Id:       1,
			UniqueId: "CARD-001",
			Pokemon:  "AAA",
			ImageUrl: "imageUrl1",
		},
		{
			Id:       2,
			UniqueId: "CARD-002",
			Pokemon:  "BBB",
			ImageUrl: "imageUrl2",
		},
		{
			Id:       3,
			UniqueId: "CARD-003",
			Pokemon:  "CCC",
			ImageUrl: "imageUrl3",
		},
	}

	_, _ = conn.Conn.NewTruncateTable().Model(&model.Card{}).Cascade().Exec(suite.ctx)
	_, err = conn.Conn.NewInsert().Model(suite.seedModels[0]).ExcludeColumn("card_id").Exec(suite.ctx)
	assert.Nil(suite.T(), err)

	_, err = conn.Conn.NewInsert().Model(suite.seedModels[1]).ExcludeColumn("card_id").Exec(suite.ctx)
	assert.Nil(suite.T(), err)

	_, err = conn.Conn.NewInsert().Model(suite.seedModels[2]).ExcludeColumn("card_id").Exec(suite.ctx)
	assert.Nil(suite.T(), err)
}

func (suite *CardAdapterTestSuite) TestNewModel() {
	adapter := NewDatabaseCardAdapter(suite.conn)
	testModel := &model.Card{
		Id:       4,
		UniqueId: "CARD-004",
		Pokemon:  "DDD",
		ImageUrl: "Image4",
	}
	createdModel, err := adapter.CreateCard(testModel)
	assert.Nil(suite.T(), err)

	results := make([]*model.Card, 0)
	suite.conn.Conn.NewSelect().Model(&model.Card{}).Scan(suite.ctx, &results)

	assert.Contains(suite.T(), results, createdModel)
}

func (suite *CardAdapterTestSuite) TestDeleteModel() {
	adapter := NewDatabaseCardAdapter(suite.conn)
	err := adapter.DeleteCard(2)
	assert.Nil(suite.T(), err)

	results := make([]*model.Card, 0)
	suite.conn.Conn.NewSelect().Model(&model.Card{}).Scan(suite.ctx, &results)

	for _, item := range results {
		if item.Id == 2 {
			suite.T().Fatal("Failed to delete card")
		}
	}
}

func (suite *CardAdapterTestSuite) TestUpdateModel() {
	adapter := NewDatabaseCardAdapter(suite.conn)
	changedModel := &model.Card{
		Id:       1,
		UniqueId: "CARD-111",
		Pokemon:  "Another",
		ImageUrl: "anotherUrl",
	}
	err := adapter.EditCard(changedModel)
	assert.Nil(suite.T(), err)

	results := make([]*model.Card, 0)
	suite.conn.Conn.NewSelect().Model(&model.Card{}).Scan(suite.ctx, &results)

	assert.Contains(suite.T(), results, changedModel)
}

func (suite *CardAdapterTestSuite) TestGetCard() {
	adapter := NewDatabaseCardAdapter(suite.conn)
	retrievedModel, err := adapter.GetCard(3)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), suite.seedModels[2], retrievedModel)

	retrievedModel, err = adapter.GetCard(10)
	assert.Nil(suite.T(), err)
	assert.Nil(suite.T(), retrievedModel)
}

func (suite *CardAdapterTestSuite) TestGetCardByUniqueId() {
	adapter := NewDatabaseCardAdapter(suite.conn)
	retrievedModel, err := adapter.GetCardByUniqueId("CARD-003")
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), suite.seedModels[2], retrievedModel)

	retrievedModel, err = adapter.GetCardByUniqueId("ASDF")
	assert.Nil(suite.T(), err)
	assert.Nil(suite.T(), retrievedModel)
}

func (suite *CardAdapterTestSuite) TestGetAllCards() {
	adapter := NewDatabaseCardAdapter(suite.conn)
	retrievedModels, err := adapter.GetAllCards()
	assert.Nil(suite.T(), err)

	assert.Greater(suite.T(), len(retrievedModels), 1)
}

func TestCardAdapterTestSuite(t *testing.T) {
	suite.Run(t, new(CardAdapterTestSuite))
}
