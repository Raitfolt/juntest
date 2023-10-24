package utils

import (
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"
)

type CountryInfo struct {
	CountryID   string  `json:"country_id"`
	Probability float64 `json:"probability"`
}

type PersonInfo struct {
	Count   int           `json:"count"`
	Name    string        `json:"name"`
	Country []CountryInfo `json:"country"`
}

func Nationalize(name string, log *zap.Logger) (string, error) {
	var person PersonInfo
	resp, err := http.Get("https://api.nationalize.io/?name=" + name)
	if err != nil {
		log.Error("get nationalize", zap.String("error", err.Error()))
		return "", err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("read nationalize response body", zap.String("error", err.Error()))
		return "", err
	}
	err = json.Unmarshal(body, &person)
	if err != nil {
		log.Error("unmarshall nationalize response body", zap.String("error", err.Error()))
	}
	resp.Body.Close()

	log.Info("get nation", zap.String("name", name), zap.String("nation", person.Country[0].CountryID))
	return person.Country[0].CountryID, nil
}
