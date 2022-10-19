package controller

import (
	"log"
	"net/http"

	"backend.cs3219.comp.nus.edu.sg/auth"
	"backend.cs3219.comp.nus.edu.sg/database"
	"backend.cs3219.comp.nus.edu.sg/model"
	"backend.cs3219.comp.nus.edu.sg/server"
	"github.com/julienschmidt/httprouter"
)

type CardController interface {
	Attach(server server.HTTPServer)
}

type cardController struct {
	baseController
	db database.DatabaseCardAdapter
}

func NewCardController(db *database.DatabaseConnection, authenticator auth.TokenAuthenticator) CardController {
	return &cardController{
		db: database.NewDatabaseCardAdapter(db),
		baseController: baseController{
			authenticator: authenticator,
		},
	}
}

func (controller *cardController) Attach(server server.HTTPServer) {
	server.Get("/api/card", controller.getAllCards)
	server.Post("/api/card", controller.createCard)
	server.Get("/api/card/:cardId", controller.getCard)
	server.Put("/api/card/:cardId", controller.editCard)
	server.Delete("/api/card/:cardId", controller.deleteCard)
}

func (controller *cardController) getAllCards(
	resp http.ResponseWriter,
	req *http.Request,
	params httprouter.Params,
) {
	if !controller.authenticateRequest(resp, req) {
		return
	}

	cards, err := controller.db.GetAllCards()
	if err != nil {
		controller.writeInternalError(resp)
		return
	}

	err = controller.writeJson(resp, cards)
	if err != nil {
		log.Println("Failed to write response for getAllCards")
	}
}

func (controller *cardController) getCard(
	resp http.ResponseWriter,
	req *http.Request,
	params httprouter.Params,
) {
	if !controller.authenticateRequest(resp, req) {
		return
	}

	cardIdParam := controller.readIntParam("cardId", params)
	if cardIdParam == nil {
		controller.writeBadRequest(resp)
		return
	}
	cardId := *cardIdParam

	card, err := controller.db.GetCard(cardId)
	if err != nil {
		controller.writeInternalError(resp)
		return
	}

	if card == nil {
		controller.writeNotFound(resp)
		return
	}

	err = controller.writeJson(resp, card)
	if err != nil {
		log.Println("Failed to write response for getCard")
	}
}

func (controller *cardController) createCard(
	resp http.ResponseWriter,
	req *http.Request,
	params httprouter.Params,
) {
	if !controller.authenticateRequest(resp, req) {
		return
	}

	var cardData model.Card
	err := controller.readJson(req, &cardData)
	if err != nil {
		controller.writeBadRequest(resp)
		return
	}

	if cardData.ImageUrl == "" || cardData.Pokemon == "" || cardData.UniqueId == "" {
		controller.writeBadRequest(resp)
		return
	}

	existingCard, err := controller.db.GetCardByUniqueId(cardData.UniqueId)
	if err != nil {
		controller.writeInternalError(resp)
		return
	}
	if existingCard != nil {
		controller.writeError(resp, 403, "A card with the same unique ID already exists")
		return
	}

	card, err := controller.db.CreateCard(&cardData)
	if err != nil {
		controller.writeInternalError(resp)
		return
	}

	err = controller.writeJson(resp, card)
	if err != nil {
		log.Println("Failed to write response for createCard")
	}
}

func (controller *cardController) editCard(
	resp http.ResponseWriter,
	req *http.Request,
	params httprouter.Params,
) {
	if !controller.authenticateRequest(resp, req) {
		return
	}

	cardIdParam := controller.readIntParam("cardId", params)
	if cardIdParam == nil {
		controller.writeBadRequest(resp)
		return
	}
	cardId := *cardIdParam

	var cardData model.Card
	err := controller.readJson(req, &cardData)
	if err != nil {
		controller.writeBadRequest(resp)
		return
	}

	if cardData.ImageUrl == "" || cardData.Pokemon == "" || cardData.UniqueId == "" {
		controller.writeBadRequest(resp)
		return
	}

	if cardData.Id != cardId {
		controller.writeBadRequest(resp)
		return
	}

	existingCard, err := controller.db.GetCardByUniqueId(cardData.UniqueId)
	if err != nil {
		controller.writeInternalError(resp)
		return
	}
	if existingCard != nil && existingCard.Id != cardId {
		controller.writeError(resp, 403, "A card with the same unique ID already exists")
		return
	}

	targetCard, err := controller.db.GetCard(cardId)
	if err != nil {
		controller.writeInternalError(resp)
		return
	}
	if targetCard == nil {
		controller.writeError(resp, 403, "No such card exists")
		return
	}

	err = controller.db.EditCard(&cardData)
	if err != nil {
		controller.writeInternalError(resp)
		return
	}

	err = controller.writeJson(resp, &cardData)
	if err != nil {
		log.Println("Failed to write response for editCard")
	}
}

func (controller *cardController) deleteCard(
	resp http.ResponseWriter,
	req *http.Request,
	params httprouter.Params,
) {
	if !controller.authenticateRequest(resp, req) {
		return
	}

	cardIdParam := controller.readIntParam("cardId", params)
	if cardIdParam == nil {
		controller.writeBadRequest(resp)
		return
	}
	cardId := *cardIdParam

	err := controller.db.DeleteCard(cardId)
	if err != nil {
		controller.writeInternalError(resp)
		return
	}

	var response struct {
		Success bool `json:"success"`
	}
	response.Success = true
	err = controller.writeJson(resp, &response)
	if err != nil {
		log.Println("Failed to write response for deleteCard")
	}
}
