package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"backend.cs3219.comp.nus.edu.sg/auth"
	"github.com/julienschmidt/httprouter"
)

type baseController struct {
	authenticator auth.TokenAuthenticator
}

func (controller *baseController) readIntParam(paramName string, params httprouter.Params) *int {
	cardIdString := params.ByName(paramName)
	if cardIdString == "" {
		return nil
	}

	cardId, err := strconv.Atoi(cardIdString)
	if err != nil {
		return nil
	}
	return &cardId
}

func (controller *baseController) readJson(req *http.Request, container interface{}) error {
	rawData, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(rawData, container)
}

func (controller *baseController) writeJson(resp http.ResponseWriter, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return controller.writeJsonType(resp, 200, jsonData)
}

func (controller *baseController) writeNotFound(resp http.ResponseWriter) {
	controller.writeJsonType(resp, 404, []byte("{}"))
}

func (controller *baseController) writeInternalError(resp http.ResponseWriter) {
	controller.writeJsonType(resp, 500, []byte("{}"))
}

func (controller *baseController) writeBadRequest(resp http.ResponseWriter) {
	controller.writeJsonType(resp, 400, []byte("{}"))
}

func (controller *baseController) writeError(resp http.ResponseWriter, code int, errorMsg string) {
	var response struct {
		ErrorMsg string `json:"errorMsg"`
	}
	response.ErrorMsg = errorMsg
	data, err := json.Marshal(response)

	if err != nil {
		controller.writeJsonType(resp, code, []byte("{}"))
	} else {
		controller.writeJsonType(resp, code, data)
	}
}

func (controller *baseController) writeJsonType(resp http.ResponseWriter, code int, data []byte) error {
	resp.Header().Add("Content-Type", "application/json")
	resp.WriteHeader(code)
	_, err := resp.Write(data)
	return err
}

func (controller *baseController) authenticateRequest(resp http.ResponseWriter, req *http.Request) bool {
	if !controller.isBearerValid(req) {
		resp.WriteHeader(403)
		resp.Write([]byte("Unauthorized"))
		return false
	}
	return true
}

func (controller *baseController) isBearerValid(req *http.Request) bool {
	authHeader := req.Header.Get("authorization")
	if authHeader == "" {
		return false
	}
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return false
	}

	bearerToken := strings.Split(authHeader, " ")[1]
	return controller.authenticator.IsValidToken(bearerToken)
}
