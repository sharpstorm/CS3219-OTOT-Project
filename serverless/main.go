/* Built on Google's Cloud Functions V1 */

package functions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type QueryResponse struct {
	ErrorMessage string     `json:"errorMessage"`
	Data         *PriceData `json:"data"`
}

type QueryRequest struct {
	CardUniqueId string `json:"cardUniqueId"`
}

type ApiResponse struct {
	Data struct {
		TcgPlayer struct {
			UpdatedAt string                 `json:"updatedAt"`
			Prices    map[string]*PriceRange `json:"prices"`
		} `json:"tcgplayer"`
	} `json:"data"`
}

type PriceData struct {
	UpdatedAt string                 `json:"updatedAt"`
	Prices    map[string]*PriceRange `json:"prices"`
}

type PriceRange struct {
	LowPrice    float32 `json:"low"`
	MidPrice    float32 `json:"mid"`
	HighPrice   float32 `json:"high"`
	MarketPrice float32 `json:"market"`
}

// helloGet is an HTTP Cloud Function.
func GetPrice(w http.ResponseWriter, r *http.Request) {
	apiToken := os.Getenv("API_KEY")
	if apiToken == "" {
		writeError(w, 500, "Server is not configured correctly")
		return
	}

	if r.Method != http.MethodPost {
		writeError(w, 405, "Only POST is allowed")
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeError(w, 400, "Bad Request Body")
		return
	}

	var reqData QueryRequest
	err = json.Unmarshal(reqBody, &reqData)
	if err != nil || reqData.CardUniqueId == "" {
		writeError(w, 400, "Bad Request Body")
		return
	}

	apiResponse, err := getTcgApiResponse(apiToken, reqData.CardUniqueId)
	if err != nil {
		writeError(w, 503, "Downstream server could not process the request")
		return
	}

	priceData, err := parseTcgApiResponse(apiResponse)
	if err != nil {
		writeError(w, 500, "Unable to parse downstream response")
		return
	}

	jsonResponse, err := json.Marshal(QueryResponse{
		ErrorMessage: "",
		Data:         priceData,
	})
	if err != nil {
		writeError(w, 500, "Internal error encountered")
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func getTcgApiResponse(apiToken string, cardId string) ([]byte, error) {
	url := fmt.Sprintf("https://api.pokemontcg.io/v2/cards/%s", cardId)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Api-Key", apiToken)
	httpClient := &http.Client{
		Timeout: time.Second * 5,
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func parseTcgApiResponse(response []byte) (*PriceData, error) {
	var dataContainer ApiResponse
	err := json.Unmarshal(response, &dataContainer)
	if err != nil {
		return nil, err
	}

	return &PriceData{
		UpdatedAt: dataContainer.Data.TcgPlayer.UpdatedAt,
		Prices:    dataContainer.Data.TcgPlayer.Prices,
	}, nil
}

func writeError(w http.ResponseWriter, code int, message string) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)

	content, err := json.Marshal(QueryResponse{
		ErrorMessage: message,
		Data:         nil,
	})

	if err != nil {
		w.Write([]byte("{}"))
		return
	}

	w.Write(content)
}
