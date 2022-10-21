package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	"backend.cs3219.comp.nus.edu.sg/database"
	"backend.cs3219.comp.nus.edu.sg/mocks"
	"backend.cs3219.comp.nus.edu.sg/model"
	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CardControllerTestSuite struct {
	suite.Suite
	dbAdapter    database.DatabaseCardAdapter
	ctx          context.Context
	seedModels   []*model.Card
	unauthHeader map[string][]string
	authHeader   map[string][]string
}

type StubResponseWriter struct {
	headers map[string][]string
	body    []byte
	status  int
}

var EMPTY_PARAMS = httprouter.Params{}

const (
	UNAUTH_TOKEN = "AAA"
	AUTH_TOKEN   = "BBB"
	VALID_URL    = "http://example.com"
	INVALID_URL  = "notaurl"
)

func newResponseWriter() *StubResponseWriter {
	return &StubResponseWriter{
		headers: make(map[string][]string),
	}
}

func (writer *StubResponseWriter) Header() http.Header {
	return writer.headers
}

func (writer *StubResponseWriter) Write(body []byte) (int, error) {
	writer.body = body
	return len(body), nil
}

func (writer *StubResponseWriter) WriteHeader(statusCode int) {
	writer.status = statusCode
}

func (suite *CardControllerTestSuite) SetupSuite() {
	suite.seedModels = []*model.Card{
		{
			Id:       101,
			UniqueId: "CARD-101",
			Pokemon:  "AAA",
			ImageUrl: "http://url1.something.com",
		},
		{
			Id:       102,
			UniqueId: "CARD-102",
			Pokemon:  "BBB",
			ImageUrl: "http://url2.something.com",
		},
		{
			Id:       103,
			UniqueId: "CARD-103",
			Pokemon:  "CCC",
			ImageUrl: "http://url3.something.com",
		},
	}
	suite.unauthHeader = map[string][]string{
		"Authorization": {
			fmt.Sprintf("Bearer %s", UNAUTH_TOKEN),
		},
	}
	suite.authHeader = map[string][]string{
		"Authorization": {
			fmt.Sprintf("Bearer %s", AUTH_TOKEN),
		},
	}
}

func (suite *CardControllerTestSuite) TestGetAllCards() {
	mockCtrl := gomock.NewController(suite.T())
	defer mockCtrl.Finish()

	authenticator := mocks.NewMockTokenAuthenticator(mockCtrl)
	cardAdapter := mocks.NewMockDatabaseCardAdapter(mockCtrl)
	gomock.InOrder(
		cardAdapter.EXPECT().GetAllCards().Return(nil, errors.New("Test error")),
		cardAdapter.EXPECT().GetAllCards().Return(suite.seedModels, nil),
	)
	gomock.InOrder(
		authenticator.EXPECT().IsValidToken(UNAUTH_TOKEN).Return(false),
		authenticator.EXPECT().IsValidToken(AUTH_TOKEN).Return(true).Times(2),
	)
	controller := &cardController{
		db: cardAdapter,
		baseController: baseController{
			authenticator: authenticator,
		},
	}

	// Case: Unauthorized GET
	request := buildHTTPRequest(suite.unauthHeader, nil)
	responseStub := newResponseWriter()
	controller.getAllCards(responseStub, request, EMPTY_PARAMS)
	assert.Equal(suite.T(), 401, responseStub.status)

	// Case: Authorized GET, DB Error
	request = buildHTTPRequest(suite.authHeader, nil)
	responseStub = newResponseWriter()
	controller.getAllCards(responseStub, request, EMPTY_PARAMS)
	assert.Equal(suite.T(), 500, responseStub.status)

	// Case: Authorized GET, No Error
	responseStub = newResponseWriter()
	controller.getAllCards(responseStub, request, EMPTY_PARAMS)
	assert.Equal(suite.T(), 200, responseStub.status)
	result := make([]model.Card, 0)
	err := json.Unmarshal(responseStub.body, &result)
	assert.Nil(suite.T(), err)
	assert.Contains(suite.T(), result, *(suite.seedModels[0]))
	assert.Contains(suite.T(), result, *(suite.seedModels[1]))
	assert.Contains(suite.T(), result, *(suite.seedModels[2]))
}

func (suite *CardControllerTestSuite) TestGetCard() {
	mockCtrl := gomock.NewController(suite.T())
	defer mockCtrl.Finish()

	authenticator := mocks.NewMockTokenAuthenticator(mockCtrl)
	cardAdapter := mocks.NewMockDatabaseCardAdapter(mockCtrl)
	gomock.InOrder(
		cardAdapter.EXPECT().GetCard(gomock.Any()).Return(nil, errors.New("Test error")),
		cardAdapter.EXPECT().GetCard(gomock.Any()).Return(suite.seedModels[0], nil),
	)
	gomock.InOrder(
		authenticator.EXPECT().IsValidToken(UNAUTH_TOKEN).Return(false),
		authenticator.EXPECT().IsValidToken(AUTH_TOKEN).Return(true).Times(4),
	)
	controller := &cardController{
		db: cardAdapter,
		baseController: baseController{
			authenticator: authenticator,
		},
	}

	// Case: Unauthorized GET
	request := buildHTTPRequest(suite.unauthHeader, nil)
	responseStub := newResponseWriter()
	controller.getCard(responseStub, request, buildRouteParams("101"))
	assert.Equal(suite.T(), 401, responseStub.status)

	// Case: No Route Params
	request = buildHTTPRequest(suite.authHeader, nil)
	responseStub = newResponseWriter()
	controller.getCard(responseStub, request, EMPTY_PARAMS)
	assert.Equal(suite.T(), 400, responseStub.status)

	// Case: Bad Route Param - Not Number
	responseStub = newResponseWriter()
	controller.getCard(responseStub, request, buildRouteParams("asdf"))
	assert.Equal(suite.T(), 400, responseStub.status)

	// Case: Authorized GET, DB Error
	responseStub = newResponseWriter()
	controller.getCard(responseStub, request, buildRouteParams("101"))
	assert.Equal(suite.T(), 500, responseStub.status)

	// Case: Authorized GET, No Error
	responseStub = newResponseWriter()
	controller.getCard(responseStub, request, buildRouteParams("101"))
	assert.Equal(suite.T(), 200, responseStub.status)
	var result model.Card
	err := json.Unmarshal(responseStub.body, &result)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), result, *(suite.seedModels[0]))
}

func (suite *CardControllerTestSuite) TestCreateCard() {
	mockCtrl := gomock.NewController(suite.T())
	defer mockCtrl.Finish()

	authenticator := mocks.NewMockTokenAuthenticator(mockCtrl)
	cardAdapter := mocks.NewMockDatabaseCardAdapter(mockCtrl)
	gomock.InOrder(
		// Card already exists
		cardAdapter.EXPECT().GetCardByUniqueId(gomock.Eq("CARD-101")).Return(suite.seedModels[0], nil),

		// DB Error 1
		cardAdapter.EXPECT().GetCardByUniqueId(gomock.Eq("CARD-101")).Return(nil, errors.New("Test Error")),

		// DB Error 2
		cardAdapter.EXPECT().GetCardByUniqueId(gomock.Eq("CARD-101")).Return(nil, nil),
		cardAdapter.EXPECT().CreateCard(gomock.Any()).Return(nil, errors.New("Test Error")),

		// Successful Create
		cardAdapter.EXPECT().GetCardByUniqueId(gomock.Any()).Return(nil, nil),
		cardAdapter.EXPECT().CreateCard(gomock.Eq(
			&model.Card{
				UniqueId: "CARD-200",
				Pokemon:  "AAA",
				ImageUrl: VALID_URL,
			},
		)).Return(&model.Card{
			Id:       200,
			UniqueId: "CARD-200",
			Pokemon:  "AAA",
			ImageUrl: VALID_URL,
		}, nil),
	)
	gomock.InOrder(
		authenticator.EXPECT().IsValidToken(UNAUTH_TOKEN).Return(false),
		authenticator.EXPECT().IsValidToken(AUTH_TOKEN).Return(true).AnyTimes(),
	)
	controller := &cardController{
		db: cardAdapter,
		baseController: baseController{
			authenticator: authenticator,
		},
	}

	// Unauthorized POST
	request := buildHTTPRequest(suite.unauthHeader, nil)
	responseStub := newResponseWriter()
	controller.createCard(responseStub, request, EMPTY_PARAMS)
	assert.Equal(suite.T(), 401, responseStub.status)

	// Authorized, no body
	request = buildHTTPRequest(suite.authHeader, "")
	responseStub = newResponseWriter()
	controller.createCard(responseStub, request, EMPTY_PARAMS)
	assert.Equal(suite.T(), 400, responseStub.status)

	// Authorized, no Unique ID
	request = buildHTTPRequest(suite.authHeader, &model.Card{
		UniqueId: "",
		Pokemon:  "AAA",
		ImageUrl: VALID_URL,
	})
	responseStub = newResponseWriter()
	controller.createCard(responseStub, request, EMPTY_PARAMS)
	assert.Equal(suite.T(), 400, responseStub.status)

	// Authorized, no Pokemon
	request = buildHTTPRequest(suite.authHeader, &model.Card{
		UniqueId: "CARD-200",
		Pokemon:  "",
		ImageUrl: VALID_URL,
	})
	responseStub = newResponseWriter()
	controller.createCard(responseStub, request, EMPTY_PARAMS)
	assert.Equal(suite.T(), 400, responseStub.status)

	// Authorized, no image URL
	request = buildHTTPRequest(suite.authHeader, &model.Card{
		UniqueId: "CARD-200",
		Pokemon:  "AAA",
		ImageUrl: "",
	})
	responseStub = newResponseWriter()
	controller.createCard(responseStub, request, EMPTY_PARAMS)
	assert.Equal(suite.T(), 400, responseStub.status)

	// Authorized, bad image URL
	request = buildHTTPRequest(suite.authHeader, &model.Card{
		UniqueId: "CARD-200",
		Pokemon:  "AAA",
		ImageUrl: INVALID_URL,
	})
	responseStub = newResponseWriter()
	controller.createCard(responseStub, request, EMPTY_PARAMS)
	assert.Equal(suite.T(), 400, responseStub.status)

	// Authorized, card already exists
	request = buildHTTPRequest(suite.authHeader, &model.Card{
		UniqueId: "CARD-101",
		Pokemon:  "AAA",
		ImageUrl: VALID_URL,
	})
	responseStub = newResponseWriter()
	controller.createCard(responseStub, request, EMPTY_PARAMS)
	assert.Equal(suite.T(), 403, responseStub.status)

	// DB Error 1
	request = buildHTTPRequest(suite.authHeader, &model.Card{
		UniqueId: "CARD-101",
		Pokemon:  "AAA",
		ImageUrl: VALID_URL,
	})
	responseStub = newResponseWriter()
	controller.createCard(responseStub, request, EMPTY_PARAMS)
	assert.Equal(suite.T(), 500, responseStub.status)

	// DB Error 2
	request = buildHTTPRequest(suite.authHeader, &model.Card{
		UniqueId: "CARD-101",
		Pokemon:  "AAA",
		ImageUrl: VALID_URL,
	})
	responseStub = newResponseWriter()
	controller.createCard(responseStub, request, EMPTY_PARAMS)
	assert.Equal(suite.T(), 500, responseStub.status)

	// Successful Create
	request = buildHTTPRequest(suite.authHeader, model.Card{
		UniqueId: "CARD-200",
		Pokemon:  "AAA",
		ImageUrl: VALID_URL,
	})
	responseStub = newResponseWriter()
	controller.createCard(responseStub, request, EMPTY_PARAMS)
	assert.Equal(suite.T(), 200, responseStub.status)
	var result model.Card
	err := json.Unmarshal(responseStub.body, &result)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), model.Card{
		Id:       200,
		UniqueId: "CARD-200",
		Pokemon:  "AAA",
		ImageUrl: VALID_URL,
	}, result)
}

func (suite *CardControllerTestSuite) TestEditCard() {
	mockCtrl := gomock.NewController(suite.T())
	defer mockCtrl.Finish()

	authenticator := mocks.NewMockTokenAuthenticator(mockCtrl)
	cardAdapter := mocks.NewMockDatabaseCardAdapter(mockCtrl)
	gomock.InOrder(
		cardAdapter.EXPECT().GetCardByUniqueId(gomock.Eq("CARD-102")).Return(suite.seedModels[1], nil),

		// Card not found
		cardAdapter.EXPECT().GetCardByUniqueId(gomock.Eq("CARD-300")).Return(nil, nil),
		cardAdapter.EXPECT().GetCard(gomock.Eq(200)).Return(nil, nil),

		// DB Error 1
		cardAdapter.EXPECT().GetCardByUniqueId(gomock.Eq("CARD-101")).Return(nil, errors.New("Test Error")),

		// DB Error 2
		cardAdapter.EXPECT().GetCardByUniqueId(gomock.Eq("CARD-101")).Return(suite.seedModels[0], nil),
		cardAdapter.EXPECT().GetCard(gomock.Eq(101)).Return(nil, errors.New("Test Error")),

		// DB Error 3
		cardAdapter.EXPECT().GetCardByUniqueId(gomock.Eq("CARD-101")).Return(suite.seedModels[0], nil),
		cardAdapter.EXPECT().GetCard(gomock.Eq(101)).Return(suite.seedModels[0], nil),
		cardAdapter.EXPECT().EditCard(gomock.Any()).Return(errors.New("Test Error")),

		// Success Call 1
		cardAdapter.EXPECT().GetCardByUniqueId(gomock.Eq("CARD-101")).Return(suite.seedModels[0], nil),
		cardAdapter.EXPECT().GetCard(gomock.Eq(101)).Return(suite.seedModels[0], nil),
		cardAdapter.EXPECT().EditCard(gomock.Eq(
			&model.Card{
				Id:       101,
				UniqueId: "CARD-101",
				Pokemon:  "XXXX",
				ImageUrl: VALID_URL,
			},
		)).Return(nil),

		// Sucess Call 2
		cardAdapter.EXPECT().GetCardByUniqueId(gomock.Eq("CARD-300")).Return(nil, nil),
		cardAdapter.EXPECT().GetCard(gomock.Eq(101)).Return(suite.seedModels[0], nil),
		cardAdapter.EXPECT().EditCard(gomock.Eq(
			&model.Card{
				Id:       101,
				UniqueId: "CARD-300",
				Pokemon:  "BBB",
				ImageUrl: VALID_URL,
			},
		)).Return(nil),
	)
	gomock.InOrder(
		authenticator.EXPECT().IsValidToken("AAA").Return(false),
		authenticator.EXPECT().IsValidToken("BBB").Return(true).AnyTimes(),
	)
	controller := &cardController{
		db: cardAdapter,
		baseController: baseController{
			authenticator: authenticator,
		},
	}

	// Unauthorized PUT
	request := buildHTTPRequest(suite.unauthHeader, nil)
	responseStub := newResponseWriter()
	controller.editCard(responseStub, request, buildRouteParams("101"))
	assert.Equal(suite.T(), 401, responseStub.status)

	// Authorized, no route param
	request = buildHTTPRequest(suite.authHeader, &model.Card{
		Id:       101,
		UniqueId: "CARD-300",
		Pokemon:  "AAA",
		ImageUrl: VALID_URL,
	})
	responseStub = newResponseWriter()
	controller.editCard(responseStub, request, EMPTY_PARAMS)
	assert.Equal(suite.T(), 400, responseStub.status)

	// Authorized, no body
	request = buildHTTPRequest(suite.authHeader, "")
	responseStub = newResponseWriter()
	controller.editCard(responseStub, request, buildRouteParams("101"))
	assert.Equal(suite.T(), 400, responseStub.status)

	// Authorized, no ID
	request = buildHTTPRequest(suite.authHeader, &model.Card{
		Id:       0,
		UniqueId: "ABCD",
		Pokemon:  "AAA",
		ImageUrl: VALID_URL,
	})
	responseStub = newResponseWriter()
	controller.editCard(responseStub, request, buildRouteParams("101"))
	assert.Equal(suite.T(), 400, responseStub.status)

	// Authorized, no Unique ID
	request = buildHTTPRequest(suite.authHeader, &model.Card{
		Id:       101,
		UniqueId: "",
		Pokemon:  "AAA",
		ImageUrl: VALID_URL,
	})
	responseStub = newResponseWriter()
	controller.editCard(responseStub, request, buildRouteParams("101"))
	assert.Equal(suite.T(), 400, responseStub.status)

	// Authorized, no Pokemon
	request = buildHTTPRequest(suite.authHeader, &model.Card{
		Id:       101,
		UniqueId: "CARD-300",
		Pokemon:  "",
		ImageUrl: VALID_URL,
	})
	responseStub = newResponseWriter()
	controller.editCard(responseStub, request, buildRouteParams("101"))
	assert.Equal(suite.T(), 400, responseStub.status)

	// Authorized, no URL
	request = buildHTTPRequest(suite.authHeader, &model.Card{
		Id:       101,
		UniqueId: "CARD-300",
		Pokemon:  "AAA",
		ImageUrl: "",
	})
	responseStub = newResponseWriter()
	controller.editCard(responseStub, request, buildRouteParams("101"))
	assert.Equal(suite.T(), 400, responseStub.status)

	// Authorized, Bad URL
	request = buildHTTPRequest(suite.authHeader, &model.Card{
		Id:       101,
		UniqueId: "CARD-300",
		Pokemon:  "AAA",
		ImageUrl: INVALID_URL,
	})
	responseStub = newResponseWriter()
	controller.editCard(responseStub, request, buildRouteParams("101"))
	assert.Equal(suite.T(), 400, responseStub.status)

	// Authorized, Route and Body ID Mismatch
	request = buildHTTPRequest(suite.authHeader, &model.Card{
		Id:       102,
		UniqueId: "ABCD",
		Pokemon:  "AAA",
		ImageUrl: VALID_URL,
	})
	responseStub = newResponseWriter()
	controller.editCard(responseStub, request, buildRouteParams("101"))
	assert.Equal(suite.T(), 400, responseStub.status)

	// Authorized, Card with same ID already exists
	request = buildHTTPRequest(suite.authHeader, &model.Card{
		Id:       101,
		UniqueId: "CARD-102",
		Pokemon:  "AAA",
		ImageUrl: VALID_URL,
	})
	responseStub = newResponseWriter()
	controller.editCard(responseStub, request, buildRouteParams("101"))
	assert.Equal(suite.T(), 403, responseStub.status)

	// Authorized, Target Card not found
	request = buildHTTPRequest(suite.authHeader, &model.Card{
		Id:       200,
		UniqueId: "CARD-300",
		Pokemon:  "AAA",
		ImageUrl: VALID_URL,
	})
	responseStub = newResponseWriter()
	controller.editCard(responseStub, request, buildRouteParams("200"))
	assert.Equal(suite.T(), 404, responseStub.status)

	// Authorized, Database Error
	request = buildHTTPRequest(suite.authHeader, &model.Card{
		Id:       101,
		UniqueId: "CARD-101",
		Pokemon:  "XXXX",
		ImageUrl: VALID_URL,
	})
	responseStub = newResponseWriter()
	controller.editCard(responseStub, request, buildRouteParams("101"))
	assert.Equal(suite.T(), 500, responseStub.status)

	// DB Error 2
	request = buildHTTPRequest(suite.authHeader, &model.Card{
		Id:       101,
		UniqueId: "CARD-101",
		Pokemon:  "XXXX",
		ImageUrl: VALID_URL,
	})
	responseStub = newResponseWriter()
	controller.editCard(responseStub, request, buildRouteParams("101"))
	assert.Equal(suite.T(), 500, responseStub.status)

	// DB Error 3
	request = buildHTTPRequest(suite.authHeader, &model.Card{
		Id:       101,
		UniqueId: "CARD-101",
		Pokemon:  "XXXX",
		ImageUrl: VALID_URL,
	})
	responseStub = newResponseWriter()
	controller.editCard(responseStub, request, buildRouteParams("101"))
	assert.Equal(suite.T(), 500, responseStub.status)

	// Authorized, Change not Unique ID field
	request = buildHTTPRequest(suite.authHeader, &model.Card{
		Id:       101,
		UniqueId: "CARD-101",
		Pokemon:  "XXXX",
		ImageUrl: VALID_URL,
	})
	responseStub = newResponseWriter()
	controller.editCard(responseStub, request, buildRouteParams("101"))
	assert.Equal(suite.T(), 200, responseStub.status)
	var result model.Card
	err := json.Unmarshal(responseStub.body, &result)
	assert.Nil(suite.T(), err)

	assert.Equal(suite.T(), model.Card{
		Id:       101,
		UniqueId: "CARD-101",
		Pokemon:  "XXXX",
		ImageUrl: VALID_URL,
	}, result)

	// Authorized, Change Unique ID
	request = buildHTTPRequest(suite.authHeader, &model.Card{
		Id:       101,
		UniqueId: "CARD-300",
		Pokemon:  "BBB",
		ImageUrl: VALID_URL,
	})
	responseStub = newResponseWriter()
	controller.editCard(responseStub, request, buildRouteParams("101"))
	assert.Equal(suite.T(), 200, responseStub.status)
	err = json.Unmarshal(responseStub.body, &result)
	assert.Nil(suite.T(), err)

	assert.Equal(suite.T(), model.Card{
		Id:       101,
		UniqueId: "CARD-300",
		Pokemon:  "BBB",
		ImageUrl: VALID_URL,
	}, result)
}

func (suite *CardControllerTestSuite) TestDeleteCard() {
	mockCtrl := gomock.NewController(suite.T())
	defer mockCtrl.Finish()

	authenticator := mocks.NewMockTokenAuthenticator(mockCtrl)
	cardAdapter := mocks.NewMockDatabaseCardAdapter(mockCtrl)
	gomock.InOrder(
		// DB Error 1
		cardAdapter.EXPECT().GetCard(gomock.Any()).Return(nil, errors.New("Test Error")),

		// DB Error 2
		cardAdapter.EXPECT().GetCard(gomock.Eq(101)).Return(suite.seedModels[0], nil),
		cardAdapter.EXPECT().DeleteCard(gomock.Any()).Return(errors.New("Test error")),

		cardAdapter.EXPECT().GetCard(gomock.Eq(101)).Return(nil, nil),

		// Successful Delete
		cardAdapter.EXPECT().GetCard(gomock.Eq(101)).Return(suite.seedModels[0], nil),
		cardAdapter.EXPECT().DeleteCard(gomock.Any()).Return(nil),
	)
	gomock.InOrder(
		authenticator.EXPECT().IsValidToken(UNAUTH_TOKEN).Return(false),
		authenticator.EXPECT().IsValidToken(AUTH_TOKEN).Return(true).Times(6),
	)
	controller := &cardController{
		db: cardAdapter,
		baseController: baseController{
			authenticator: authenticator,
		},
	}

	// Unauthorized DELETE
	request := buildHTTPRequest(suite.unauthHeader, nil)
	responseStub := newResponseWriter()
	controller.deleteCard(responseStub, request, buildRouteParams("101"))
	assert.Equal(suite.T(), 401, responseStub.status)

	// Authorized, Empty route params
	request = buildHTTPRequest(suite.authHeader, nil)
	responseStub = newResponseWriter()
	controller.deleteCard(responseStub, request, EMPTY_PARAMS)
	assert.Equal(suite.T(), 400, responseStub.status)

	// Authorized, Bad Card ID
	request = buildHTTPRequest(suite.authHeader, nil)
	responseStub = newResponseWriter()
	controller.deleteCard(responseStub, request, buildRouteParams("asdf"))
	assert.Equal(suite.T(), 400, responseStub.status)

	// Authorized, Database error
	request = buildHTTPRequest(suite.authHeader, nil)
	responseStub = newResponseWriter()
	controller.deleteCard(responseStub, request, buildRouteParams("101"))
	assert.Equal(suite.T(), 500, responseStub.status)

	// Authorized, Database error
	request = buildHTTPRequest(suite.authHeader, nil)
	responseStub = newResponseWriter()
	controller.deleteCard(responseStub, request, buildRouteParams("101"))
	assert.Equal(suite.T(), 500, responseStub.status)

	// Authorized, Card not found
	request = buildHTTPRequest(suite.authHeader, nil)
	responseStub = newResponseWriter()
	controller.deleteCard(responseStub, request, buildRouteParams("101"))
	assert.Equal(suite.T(), 404, responseStub.status)

	// Successful delete
	request = buildHTTPRequest(suite.authHeader, nil)
	responseStub = newResponseWriter()
	controller.deleteCard(responseStub, request, buildRouteParams("101"))
	assert.Equal(suite.T(), 200, responseStub.status)
	var result map[string]bool
	err := json.Unmarshal(responseStub.body, &result)
	assert.Nil(suite.T(), err)
	assert.True(suite.T(), result["success"])
}

func TestCardControllerTestSuite(t *testing.T) {
	suite.Run(t, new(CardControllerTestSuite))
}

func buildHTTPRequest(headers map[string][]string, bodyData interface{}) *http.Request {
	req := &http.Request{
		Header: headers,
	}

	if bodyData != nil {
		jsonContent, _ := json.Marshal(bodyData)
		req.Body = io.NopCloser(bytes.NewReader(jsonContent))
	}

	return req
}

func buildRouteParams(cardId string) httprouter.Params {
	return httprouter.Params{
		{
			Key:   "cardId",
			Value: cardId,
		},
	}
}
