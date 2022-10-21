package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"backend.cs3219.comp.nus.edu.sg/auth"
	"backend.cs3219.comp.nus.edu.sg/controller"
	"backend.cs3219.comp.nus.edu.sg/database"
	"backend.cs3219.comp.nus.edu.sg/model"
	"backend.cs3219.comp.nus.edu.sg/server"
	"backend.cs3219.comp.nus.edu.sg/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type E2ESuite struct {
	suite.Suite
	ctx        context.Context
	seedTokens []*model.ApiToken
	seedCards  []*model.Card
	server     *http.Server
	baseUrl    string

	unauthHeader map[string][]string
	authHeader   map[string][]string
}

func (suite *E2ESuite) SetupSuite() {
	suite.ctx = context.Background()
	suite.unauthHeader = map[string][]string{
		"Authorization": {
			fmt.Sprintf("Bearer ABCD"),
		},
	}
	suite.authHeader = map[string][]string{
		"Authorization": {
			fmt.Sprintf("Bearer BBB"),
		},
	}
	appConfig := util.LoadEnvVariables()
	dbConn, err := database.ConnectDatabase(
		appConfig.DbUrl,
		appConfig.DbUsername,
		appConfig.DbPassword,
		appConfig.DbName,
	)
	if err != nil {
		suite.T().Fatal("Failed to connect to database")
		return
	}

	tokenAuthenticator := auth.NewTokenAuthenticator(dbConn)
	server := server.CreateHTTPServer(uint16(appConfig.Port))
	controller := controller.NewCardController(dbConn, tokenAuthenticator)
	controller.Attach(server)
	suite.server = &http.Server{
		Handler: server.GetRouter(),
		Addr:    fmt.Sprintf(":%d", appConfig.Port),
	}
	suite.baseUrl = fmt.Sprintf("http://localhost:%d", appConfig.Port)

	go func() {
		suite.server.ListenAndServe()
	}()

	suite.seedTokens = []*model.ApiToken{
		{
			Id:        1,
			Token:     "BBB",
			IsEnabled: true,
			CreatedAt: time.Now(),
		},
		{
			Id:        2,
			Token:     "CCC",
			IsEnabled: true,
			CreatedAt: time.Now(),
		},
	}

	suite.seedCards = []*model.Card{
		{
			Id:       1,
			UniqueId: "C1",
			Pokemon:  "AAA",
			ImageUrl: "http://example.com/a",
		},
		{
			Id:       2,
			UniqueId: "C2",
			Pokemon:  "BBB",
			ImageUrl: "http://example.com/b",
		},
		{
			Id:       3,
			UniqueId: "C3",
			Pokemon:  "CCC",
			ImageUrl: "http://example.com/c",
		},
		{
			Id:       4,
			UniqueId: "C4",
			Pokemon:  "DDD",
			ImageUrl: "http://example.com/d",
		},
	}

	_, _ = dbConn.Conn.NewTruncateTable().Model(&model.ApiToken{}).Cascade().Exec(suite.ctx)
	_, _ = dbConn.Conn.NewTruncateTable().Model(&model.Card{}).Cascade().Exec(suite.ctx)
	for _, item := range suite.seedTokens {
		_, err = dbConn.Conn.NewInsert().Model(item).ExcludeColumn("token_id").Exec(suite.ctx)
		assert.Nil(suite.T(), err)
	}
	for _, item := range suite.seedCards {
		_, err = dbConn.Conn.NewInsert().Model(item).ExcludeColumn("card_id").Exec(suite.ctx)
		assert.Nil(suite.T(), err)
	}
}

func (suite *E2ESuite) TearDownTestSuite() {
	suite.server.Shutdown(suite.ctx)
}

func (suite *E2ESuite) Test_A_Get() {
	suite.launchRequest(
		suite.newRequest(http.MethodGet, "/api/card", nil, suite.unauthHeader),
		401,
	)

	resp := suite.launchRequest(
		suite.newRequest(http.MethodGet, "/api/card", nil, suite.authHeader),
		200,
	)
	var cards []*model.Card
	assert.Nil(suite.T(), suite.parseJSONResponse(resp, &cards))
	assert.Equal(suite.T(), 4, len(cards))
}

func (suite *E2ESuite) Test_B_GetSpecific() {
	suite.launchRequest(
		suite.newRequest(http.MethodGet, "/api/card/1", nil, suite.unauthHeader),
		401,
	)

	suite.launchRequest(
		suite.newRequest(http.MethodGet, "/api/card/asdf", nil, suite.authHeader),
		400,
	)

	resp := suite.launchRequest(
		suite.newRequest(http.MethodGet, "/api/card/1", nil, suite.authHeader),
		200,
	)
	var card *model.Card
	assert.Nil(suite.T(), suite.parseJSONResponse(resp, &card))
	assert.Equal(suite.T(), suite.seedCards[0], card)
}

func (suite *E2ESuite) Test_C_Create() {
	refCard := &model.Card{
		UniqueId: "C1",
		Pokemon:  "AAA",
		ImageUrl: "http://example.com/a",
	}

	copy := *refCard
	suite.launchRequest(
		suite.newRequest(http.MethodPost, "/api/card", &copy, suite.unauthHeader),
		401,
	)

	copy = *refCard
	copy.UniqueId = ""
	suite.launchRequest(
		suite.newRequest(http.MethodPost, "/api/card", &copy, suite.authHeader),
		400,
	)

	copy = *refCard
	copy.Pokemon = ""
	suite.launchRequest(
		suite.newRequest(http.MethodPost, "/api/card", &copy, suite.authHeader),
		400,
	)

	copy = *refCard
	copy.ImageUrl = ""
	suite.launchRequest(
		suite.newRequest(http.MethodPost, "/api/card", &copy, suite.authHeader),
		400,
	)

	copy = *refCard
	copy.ImageUrl = "not_a_url"
	suite.launchRequest(
		suite.newRequest(http.MethodPost, "/api/card", &copy, suite.authHeader),
		400,
	)

	copy = *refCard
	suite.launchRequest(
		suite.newRequest(http.MethodPost, "/api/card", &copy, suite.authHeader),
		403,
	)

	copy = *refCard
	copy.UniqueId = "A1"
	resp := suite.launchRequest(
		suite.newRequest(http.MethodPost, "/api/card", &copy, suite.authHeader),
		200,
	)
	var card *model.Card
	assert.Nil(suite.T(), suite.parseJSONResponse(resp, &card))
	assert.Equal(suite.T(), copy.UniqueId, card.UniqueId)
	assert.Equal(suite.T(), copy.Pokemon, card.Pokemon)
	assert.Equal(suite.T(), copy.ImageUrl, card.ImageUrl)

	// Check Records
	resp = suite.launchRequest(
		suite.newRequest(http.MethodGet, "/api/card", nil, suite.authHeader),
		200,
	)
	var cards []*model.Card
	assert.Nil(suite.T(), suite.parseJSONResponse(resp, &cards))
	assert.Equal(suite.T(), 5, len(cards))
}

func (suite *E2ESuite) Test_D_Edit() {
	refCard := &model.Card{
		Id:       1,
		UniqueId: "C1",
		Pokemon:  "BBB",
		ImageUrl: "http://example.com/a",
	}

	copy := *refCard
	suite.launchRequest(
		suite.newRequest(http.MethodPut, "/api/card/1", &copy, suite.unauthHeader),
		401,
	)

	copy = *refCard
	copy.UniqueId = ""
	suite.launchRequest(
		suite.newRequest(http.MethodPut, "/api/card/1", &copy, suite.authHeader),
		400,
	)

	copy = *refCard
	copy.Pokemon = ""
	suite.launchRequest(
		suite.newRequest(http.MethodPut, "/api/card/1", &copy, suite.authHeader),
		400,
	)

	copy = *refCard
	copy.ImageUrl = ""
	suite.launchRequest(
		suite.newRequest(http.MethodPut, "/api/card/1", &copy, suite.authHeader),
		400,
	)

	copy = *refCard
	copy.ImageUrl = "not_a_url"
	suite.launchRequest(
		suite.newRequest(http.MethodPut, "/api/card/1", &copy, suite.authHeader),
		400,
	)

	copy = *refCard
	suite.launchRequest(
		suite.newRequest(http.MethodPut, "/api/card/2", &copy, suite.authHeader),
		400,
	)

	copy = *refCard
	copy.UniqueId = "C2"
	suite.launchRequest(
		suite.newRequest(http.MethodPut, "/api/card/1", &copy, suite.authHeader),
		403,
	)

	copy = *refCard
	copy.Id = 100
	copy.UniqueId = "D2"
	suite.launchRequest(
		suite.newRequest(http.MethodPut, "/api/card/100", &copy, suite.authHeader),
		404,
	)

	copy = *refCard
	resp := suite.launchRequest(
		suite.newRequest(http.MethodPut, "/api/card/1", &copy, suite.authHeader),
		200,
	)
	var card *model.Card
	assert.Nil(suite.T(), suite.parseJSONResponse(resp, &card))
	assert.Equal(suite.T(), copy.Id, card.Id)
	assert.Equal(suite.T(), copy.UniqueId, card.UniqueId)
	assert.Equal(suite.T(), copy.Pokemon, card.Pokemon)
	assert.Equal(suite.T(), copy.ImageUrl, card.ImageUrl)

	resp = suite.launchRequest(
		suite.newRequest(http.MethodGet, "/api/card/1", nil, suite.authHeader),
		200,
	)
	assert.Nil(suite.T(), suite.parseJSONResponse(resp, &card))
	assert.Equal(suite.T(), copy.Id, card.Id)
	assert.Equal(suite.T(), copy.UniqueId, card.UniqueId)
	assert.Equal(suite.T(), copy.Pokemon, card.Pokemon)
	assert.Equal(suite.T(), copy.ImageUrl, card.ImageUrl)
}

func (suite *E2ESuite) Test_E_Delete() {
	suite.launchRequest(
		suite.newRequest(http.MethodDelete, "/api/card/1", nil, suite.unauthHeader),
		401,
	)

	suite.launchRequest(
		suite.newRequest(http.MethodDelete, "/api/card/100", nil, suite.authHeader),
		404,
	)

	suite.launchRequest(
		suite.newRequest(http.MethodDelete, "/api/card/3", nil, suite.authHeader),
		200,
	)

	// Check Records
	resp := suite.launchRequest(
		suite.newRequest(http.MethodGet, "/api/card", nil, suite.authHeader),
		200,
	)
	var cards []*model.Card
	assert.Nil(suite.T(), suite.parseJSONResponse(resp, &cards))
	assert.Equal(suite.T(), 4, len(cards))
}

func (suite *E2ESuite) launchRequest(req *http.Request, expectedStatus int) *http.Response {
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := httpClient.Do(req)
	assert.Nil(suite.T(), err)
	suite.assertResponseStatus(resp, expectedStatus)
	return resp
}

func (suite *E2ESuite) newRequest(method string, endpoint string, body interface{}, headers map[string][]string) *http.Request {
	var bodyStream io.Reader = nil
	if body != nil {
		jsonContent, _ := json.Marshal(body)
		bodyStream = io.NopCloser(bytes.NewReader(jsonContent))
	}

	req, err := http.NewRequest(method, suite.baseUrl+endpoint, bodyStream)
	req.Header = headers
	assert.Nil(suite.T(), err)

	return req
}

func (suite *E2ESuite) assertResponseStatus(resp *http.Response, status int) {
	assert.Equal(suite.T(), status, resp.StatusCode)
}

func (suite *E2ESuite) parseJSONResponse(resp *http.Response, container interface{}) error {
	bodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(bodyData, &container)
}

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, new(E2ESuite))
}
