package database

import (
	"context"
	"testing"
	"time"

	"backend.cs3219.comp.nus.edu.sg/model"
	"backend.cs3219.comp.nus.edu.sg/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ApiTokenAdapterTestSuite struct {
	suite.Suite
	conn       *DatabaseConnection
	ctx        context.Context
	seedModels []*model.ApiToken
}

func (suite *ApiTokenAdapterTestSuite) SetupSuite() {
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

	suite.seedModels = []*model.ApiToken{
		{
			Id:        1,
			Token:     "AAA",
			IsEnabled: true,
			CreatedAt: time.Now(),
		},
		{
			Id:        2,
			Token:     "BBB",
			IsEnabled: true,
			CreatedAt: time.Now(),
		},
		{
			Id:        3,
			Token:     "CCC",
			IsEnabled: false,
			CreatedAt: time.Now(),
		},
		{
			Id:        4,
			Token:     "DDD",
			IsEnabled: true,
			CreatedAt: time.Now(),
		},
	}

	_, _ = conn.Conn.NewTruncateTable().Model(&model.ApiToken{}).Cascade().Exec(suite.ctx)
	for _, item := range suite.seedModels {
		_, err = conn.Conn.NewInsert().Model(item).ExcludeColumn("token_id").Exec(suite.ctx)
		assert.Nil(suite.T(), err)
	}
}

func (suite *ApiTokenAdapterTestSuite) TestCreateApiToken() {
	adapter := NewDatabaseApiTokenAdapter(suite.conn)
	token := "EEE"
	err := adapter.CreateApiToken(token)
	assert.Nil(suite.T(), err)

	results := make([]*model.ApiToken, 0)
	suite.conn.Conn.NewSelect().Model(&model.ApiToken{}).Where("token = ?", token).Scan(suite.ctx, &results)

	assert.Equal(suite.T(), 1, len(results))
}

func (suite *ApiTokenAdapterTestSuite) TestSetApiTokenState() {
	adapter := NewDatabaseApiTokenAdapter(suite.conn)
	err := adapter.SetApiTokenState(1, false)
	assert.Nil(suite.T(), err)

	results := make([]*model.ApiToken, 0)
	suite.conn.Conn.NewSelect().Model(&model.ApiToken{}).Where("token_id = ?", 1).Scan(suite.ctx, &results)

	assert.Equal(suite.T(), 1, len(results))
	assert.False(suite.T(), results[0].IsEnabled)
}

func (suite *ApiTokenAdapterTestSuite) TestDeleteApiToken() {
	adapter := NewDatabaseApiTokenAdapter(suite.conn)
	err := adapter.DeleteApiToken(2)
	assert.Nil(suite.T(), err)

	results := make([]*model.ApiToken, 0)
	suite.conn.Conn.NewSelect().Model(&model.ApiToken{}).Where("token_id = ?", 2).Scan(suite.ctx, &results)

	assert.Equal(suite.T(), 0, len(results))
}

func (suite *ApiTokenAdapterTestSuite) TestIsValidToken() {
	adapter := NewDatabaseApiTokenAdapter(suite.conn)
	isValid, err := adapter.IsValidToken("CCC")
	assert.Nil(suite.T(), err)
	assert.False(suite.T(), isValid)

	isValid, err = adapter.IsValidToken("DDD")
	assert.Nil(suite.T(), err)
	assert.True(suite.T(), isValid)
}

func (suite *ApiTokenAdapterTestSuite) TestGetAllApiTokens() {
	adapter := NewDatabaseApiTokenAdapter(suite.conn)
	tokens, err := adapter.GetAllApiTokens()
	assert.Nil(suite.T(), err)
	assert.Greater(suite.T(), len(tokens), 2)
}

func TestApiTokenAdapterTestSuite(t *testing.T) {
	suite.Run(t, new(ApiTokenAdapterTestSuite))
}
