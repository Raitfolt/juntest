package utils

import (
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"
)

type Gender struct {
	Count       int     `json:"count"`
	Name        string  `json:"name"`
	Gender      string  `json:"gender"`
	Probability float64 `json:"probability"`
}

func Genderize(name string, log *zap.Logger) (string, error) {
	var gender Gender
	resp, err := http.Get("https://api.genderize.io/?name=" + name)
	if err != nil {
		log.Error("get gender", zap.String("error", err.Error()))
		return "", err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("read genderize response body", zap.String("error", err.Error()))
		return "", err
	}
	err = json.Unmarshal(body, &gender)
	if err != nil {
		log.Error("unmarshall genderize response body", zap.String("error", err.Error()))
	}
	resp.Body.Close()
	log.Info("get gender", zap.String("name", name), zap.String("gender", gender.Gender))
	return gender.Gender, nil
}
