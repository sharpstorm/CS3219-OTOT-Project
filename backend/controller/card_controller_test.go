package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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
	dbAdapter  database.DatabaseCardAdapter
	ctx        context.Context
	seedModels []*model.Card
}

type StubResponseWriter struct {
	headers map[string][]string
	body    []byte
	status  int
}

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
			ImageUrl: "imageUrl1",
		},
		{
			Id:       102,
			UniqueId: "CARD-102",
			Pokemon:  "BBB",
			ImageUrl: "imageUrl2",
		},
		{
			Id:       103,
			UniqueId: "CARD-103",
			Pokemon:  "CCC",
			ImageUrl: "imageUrl3",
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
		authenticator.EXPECT().IsValidToken("AAA").Return(false),
		authenticator.EXPECT().IsValidToken("BBB").Return(true).Times(2),
	)
	controller := &cardController{
		db: cardAdapter,
		baseController: baseController{
			authenticator: authenticator,
		},
	}

	request := &http.Request{
		Header: map[string][]string{
			"Authorization": {
				"Bearer AAA",
			},
		},
	}
	responseStub := newResponseWriter()
	controller.getAllCards(responseStub, request, httprouter.Params{})
	assert.NotEqual(suite.T(), responseStub.status, 200)

	request = &http.Request{
		Header: map[string][]string{
			"Authorization": {
				"Bearer BBB",
			},
		},
	}
	responseStub = newResponseWriter()
	controller.getAllCards(responseStub, request, httprouter.Params{})
	assert.NotEqual(suite.T(), 200, responseStub.status)

	responseStub = newResponseWriter()
	controller.getAllCards(responseStub, request, httprouter.Params{})
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
		authenticator.EXPECT().IsValidToken("AAA").Return(false),
		authenticator.EXPECT().IsValidToken("BBB").Return(true).Times(4),
	)
	controller := &cardController{
		db: cardAdapter,
		baseController: baseController{
			authenticator: authenticator,
		},
	}

	request := &http.Request{
		Header: map[string][]string{
			"Authorization": {
				"Bearer AAA",
			},
		},
	}
	responseStub := newResponseWriter()
	controller.getCard(responseStub, request, httprouter.Params{
		{
			Key:   "cardId",
			Value: "101",
		},
	})
	assert.NotEqual(suite.T(), responseStub.status, 200)

	request = &http.Request{
		Header: map[string][]string{
			"Authorization": {
				"Bearer BBB",
			},
		},
	}

	responseStub = newResponseWriter()
	controller.getCard(responseStub, request, httprouter.Params{})
	assert.NotEqual(suite.T(), 200, responseStub.status)

	responseStub = newResponseWriter()
	controller.getCard(responseStub, request, httprouter.Params{
		{
			Key:   "cardId",
			Value: "asdf",
		},
	})
	assert.NotEqual(suite.T(), 200, responseStub.status)

	responseStub = newResponseWriter()
	controller.getCard(responseStub, request, httprouter.Params{
		{
			Key:   "cardId",
			Value: "101",
		},
	})
	assert.NotEqual(suite.T(), 200, responseStub.status)

	responseStub = newResponseWriter()
	controller.getCard(responseStub, request, httprouter.Params{
		{
			Key:   "cardId",
			Value: "101",
		},
	})
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
		cardAdapter.EXPECT().GetCardByUniqueId(gomock.Eq("CARD-101")).Return(suite.seedModels[0], nil),
		cardAdapter.EXPECT().GetCardByUniqueId(gomock.Any()).Return(nil, nil),
		cardAdapter.EXPECT().CreateCard(gomock.Eq(
			&model.Card{
				UniqueId: "CARD-200",
				Pokemon:  "AAA",
				ImageUrl: "Image",
			},
		)).Return(&model.Card{
			Id:       200,
			UniqueId: "CARD-200",
			Pokemon:  "AAA",
			ImageUrl: "Image",
		}, nil),
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

	request := &http.Request{
		Header: map[string][]string{
			"Authorization": {
				"Bearer AAA",
			},
		},
	}
	responseStub := newResponseWriter()
	controller.createCard(responseStub, request, httprouter.Params{})
	assert.NotEqual(suite.T(), responseStub.status, 200)

	request = &http.Request{
		Header: map[string][]string{
			"Authorization": {
				"Bearer BBB",
			},
		},
		Body: io.NopCloser(bytes.NewReader([]byte(""))),
	}

	responseStub = newResponseWriter()
	controller.createCard(responseStub, request, httprouter.Params{})
	assert.NotEqual(suite.T(), 200, responseStub.status)

	request = &http.Request{
		Header: map[string][]string{
			"Authorization": {
				"Bearer BBB",
			},
		},
		Body: io.NopCloser(bytes.NewReader([]byte(""))),
	}
	responseStub = newResponseWriter()
	controller.createCard(responseStub, request, httprouter.Params{})
	assert.NotEqual(suite.T(), 200, responseStub.status)

	bodyData, _ := json.Marshal(model.Card{
		UniqueId: "",
		Pokemon:  "AAA",
		ImageUrl: "Image",
	})
	request.Body = io.NopCloser(bytes.NewReader(bodyData))
	responseStub = newResponseWriter()
	controller.createCard(responseStub, request, httprouter.Params{})
	assert.NotEqual(suite.T(), 200, responseStub.status)

	bodyData, _ = json.Marshal(model.Card{
		UniqueId: "CARD-200",
		Pokemon:  "",
		ImageUrl: "Image",
	})
	request.Body = io.NopCloser(bytes.NewReader(bodyData))
	responseStub = newResponseWriter()
	controller.createCard(responseStub, request, httprouter.Params{})
	assert.NotEqual(suite.T(), 200, responseStub.status)

	bodyData, _ = json.Marshal(model.Card{
		UniqueId: "CARD-200",
		Pokemon:  "AAA",
		ImageUrl: "",
	})
	request.Body = io.NopCloser(bytes.NewReader(bodyData))
	responseStub = newResponseWriter()
	controller.createCard(responseStub, request, httprouter.Params{})
	assert.NotEqual(suite.T(), 200, responseStub.status)

	bodyData, _ = json.Marshal(model.Card{
		UniqueId: "CARD-101",
		Pokemon:  "AAA",
		ImageUrl: "Image",
	})
	request.Body = io.NopCloser(bytes.NewReader(bodyData))
	responseStub = newResponseWriter()
	controller.createCard(responseStub, request, httprouter.Params{})
	assert.NotEqual(suite.T(), 200, responseStub.status)

	bodyData, _ = json.Marshal(model.Card{
		UniqueId: "CARD-200",
		Pokemon:  "AAA",
		ImageUrl: "Image",
	})
	request.Body = io.NopCloser(bytes.NewReader(bodyData))
	responseStub = newResponseWriter()
	controller.createCard(responseStub, request, httprouter.Params{})
	assert.Equal(suite.T(), 200, responseStub.status)
	var result model.Card
	err := json.Unmarshal(responseStub.body, &result)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), model.Card{
		Id:       200,
		UniqueId: "CARD-200",
		Pokemon:  "AAA",
		ImageUrl: "Image",
	}, result)
}

func (suite *CardControllerTestSuite) TestEditCard() {
	mockCtrl := gomock.NewController(suite.T())
	defer mockCtrl.Finish()

	authenticator := mocks.NewMockTokenAuthenticator(mockCtrl)
	cardAdapter := mocks.NewMockDatabaseCardAdapter(mockCtrl)
	gomock.InOrder(
		cardAdapter.EXPECT().GetCardByUniqueId(gomock.Eq("CARD-102")).Return(suite.seedModels[0], nil),
		cardAdapter.EXPECT().GetCardByUniqueId(gomock.Eq("CARD-300")).Return(nil, nil),
		cardAdapter.EXPECT().GetCard(gomock.Eq(101)).Return(suite.seedModels[0], nil),
		cardAdapter.EXPECT().EditCard(gomock.Eq(
			&model.Card{
				Id:       101,
				UniqueId: "CARD-300",
				Pokemon:  "BBB",
				ImageUrl: "Image2",
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

	request := &http.Request{
		Header: map[string][]string{
			"Authorization": {
				"Bearer AAA",
			},
		},
	}
	responseStub := newResponseWriter()
	controller.editCard(responseStub, request, httprouter.Params{})
	assert.NotEqual(suite.T(), responseStub.status, 200)

	bodyData, _ := json.Marshal(model.Card{
		Id:       101,
		UniqueId: "CARD-300",
		Pokemon:  "AAA",
		ImageUrl: "Image",
	})
	request = &http.Request{
		Header: map[string][]string{
			"Authorization": {
				"Bearer BBB",
			},
		},
		Body: io.NopCloser(bytes.NewReader(bodyData)),
	}

	responseStub = newResponseWriter()
	controller.editCard(responseStub, request, httprouter.Params{})
	assert.NotEqual(suite.T(), 200, responseStub.status)

	request = &http.Request{
		Header: map[string][]string{
			"Authorization": {
				"Bearer BBB",
			},
		},
		Body: io.NopCloser(bytes.NewReader([]byte(""))),
	}
	httpParams := httprouter.Params{
		{
			Key:   "cardId",
			Value: "101",
		},
	}
	responseStub = newResponseWriter()
	controller.editCard(responseStub, request, httpParams)
	assert.NotEqual(suite.T(), 200, responseStub.status)

	bodyData, _ = json.Marshal(model.Card{
		Id:       101,
		UniqueId: "",
		Pokemon:  "AAA",
		ImageUrl: "Image",
	})
	request.Body = io.NopCloser(bytes.NewReader(bodyData))
	responseStub = newResponseWriter()
	controller.editCard(responseStub, request, httpParams)
	assert.NotEqual(suite.T(), 200, responseStub.status)

	bodyData, _ = json.Marshal(model.Card{
		Id:       101,
		UniqueId: "CARD-201",
		Pokemon:  "",
		ImageUrl: "Image",
	})
	request.Body = io.NopCloser(bytes.NewReader(bodyData))
	responseStub = newResponseWriter()
	controller.editCard(responseStub, request, httpParams)
	assert.NotEqual(suite.T(), 200, responseStub.status)

	bodyData, _ = json.Marshal(model.Card{
		Id:       101,
		UniqueId: "CARD-201",
		Pokemon:  "AAA",
		ImageUrl: "",
	})
	request.Body = io.NopCloser(bytes.NewReader(bodyData))
	responseStub = newResponseWriter()
	controller.editCard(responseStub, request, httpParams)
	assert.NotEqual(suite.T(), 200, responseStub.status)

	bodyData, _ = json.Marshal(model.Card{
		Id:       102,
		UniqueId: "CARD-200",
		Pokemon:  "AAA",
		ImageUrl: "Image",
	})
	request.Body = io.NopCloser(bytes.NewReader(bodyData))
	responseStub = newResponseWriter()
	controller.editCard(responseStub, request, httpParams)
	assert.NotEqual(suite.T(), 200, responseStub.status)

	bodyData, _ = json.Marshal(model.Card{
		Id:       101,
		UniqueId: "CARD-102",
		Pokemon:  "AAA",
		ImageUrl: "Image",
	})
	request.Body = io.NopCloser(bytes.NewReader(bodyData))
	responseStub = newResponseWriter()
	controller.editCard(responseStub, request, httpParams)
	assert.NotEqual(suite.T(), 200, responseStub.status)

	bodyData, _ = json.Marshal(model.Card{
		Id:       101,
		UniqueId: "CARD-300",
		Pokemon:  "BBB",
		ImageUrl: "Image2",
	})
	request.Body = io.NopCloser(bytes.NewReader(bodyData))
	responseStub = newResponseWriter()
	controller.editCard(responseStub, request, httpParams)
	assert.Equal(suite.T(), 200, responseStub.status)
	var result model.Card
	err := json.Unmarshal(responseStub.body, &result)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), model.Card{
		Id:       101,
		UniqueId: "CARD-300",
		Pokemon:  "BBB",
		ImageUrl: "Image2",
	}, result)
}

func (suite *CardControllerTestSuite) TestDeleteCard() {
	mockCtrl := gomock.NewController(suite.T())
	defer mockCtrl.Finish()

	authenticator := mocks.NewMockTokenAuthenticator(mockCtrl)
	cardAdapter := mocks.NewMockDatabaseCardAdapter(mockCtrl)
	gomock.InOrder(
		cardAdapter.EXPECT().DeleteCard(gomock.Any()).Return(errors.New("Test error")),
		cardAdapter.EXPECT().DeleteCard(gomock.Any()).Return(nil),
	)
	gomock.InOrder(
		authenticator.EXPECT().IsValidToken("AAA").Return(false),
		authenticator.EXPECT().IsValidToken("BBB").Return(true).Times(4),
	)
	controller := &cardController{
		db: cardAdapter,
		baseController: baseController{
			authenticator: authenticator,
		},
	}

	request := &http.Request{
		Header: map[string][]string{
			"Authorization": {
				"Bearer AAA",
			},
		},
	}
	responseStub := newResponseWriter()
	controller.deleteCard(responseStub, request, httprouter.Params{
		{
			Key:   "cardId",
			Value: "101",
		},
	})
	assert.NotEqual(suite.T(), responseStub.status, 200)

	request = &http.Request{
		Header: map[string][]string{
			"Authorization": {
				"Bearer BBB",
			},
		},
	}

	responseStub = newResponseWriter()
	controller.deleteCard(responseStub, request, httprouter.Params{})
	assert.NotEqual(suite.T(), 200, responseStub.status)

	responseStub = newResponseWriter()
	controller.deleteCard(responseStub, request, httprouter.Params{
		{
			Key:   "cardId",
			Value: "asdf",
		},
	})
	assert.NotEqual(suite.T(), 200, responseStub.status)

	responseStub = newResponseWriter()
	controller.deleteCard(responseStub, request, httprouter.Params{
		{
			Key:   "cardId",
			Value: "101",
		},
	})
	assert.NotEqual(suite.T(), 200, responseStub.status)

	responseStub = newResponseWriter()
	controller.deleteCard(responseStub, request, httprouter.Params{
		{
			Key:   "cardId",
			Value: "101",
		},
	})
	assert.Equal(suite.T(), 200, responseStub.status)
	var result map[string]bool
	err := json.Unmarshal(responseStub.body, &result)
	assert.Nil(suite.T(), err)
	assert.True(suite.T(), result["success"])
}

func TestCardControllerTestSuite(t *testing.T) {
	suite.Run(t, new(CardControllerTestSuite))
}
