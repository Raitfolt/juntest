package add

import (
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"
)

type Person struct {
	Name       string `json:"name" validate:"required"`
	Surname    string `json:"surname" validate:"required"`
	Patronymic string `json:"patronymic,omitempty"`
}

type PersonSaver interface {
	NewPerson(name, surname, patronymic string) (int64, error)
}

func New(log *zap.Logger, personSaver PersonSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("add new person")

		var person Person

		data, err := io.ReadAll(r.Body)
		if err != nil {
			log.Error("read body from request", zap.String("error", err.Error()))
		}

		err = json.Unmarshal(data, &person)
		if err != nil {
			log.Error("unmarshall body", zap.String("error", err.Error()))
		}
		//TODO: need to complete
	}
}
