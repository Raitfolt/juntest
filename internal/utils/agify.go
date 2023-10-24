package utils

import (
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"
)

type Age struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
}

func Agify(name string, log *zap.Logger) (int, error) {
	var age Age
	resp, err := http.Get("https://api.agify.io/?name=" + name)
	if err != nil {
		log.Error("get age", zap.String("error", err.Error()))
		return 0, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("read agify response body", zap.String("error", err.Error()))
		return 0, err
	}
	err = json.Unmarshal(body, &age)
	if err != nil {
		log.Error("unmarshall agify response body", zap.String("error", err.Error()))
	}
	resp.Body.Close()
	log.Info("get age", zap.String("name", name), zap.Int("age", age.Age))
	return age.Age, nil
}
