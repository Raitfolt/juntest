package add

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Raitfolt/juntest/internal/utils"
	"go.uber.org/zap"
)

type Person struct {
	Name        string `json:"name" validate:"required"`
	Surname     string `json:"surname" validate:"required"`
	Patronymic  string `json:"patronymic,omitempty"`
	Age         int    `json:"age,omitempty"`
	Gender      string `json:"gender,omitempty"`
	Nationality string `json:"nationality,omitempty"`
}

type Gender struct {
	Count       int     `json:"count"`
	Name        string  `json:"name"`
	Gender      string  `json:"gender"`
	Probability float32 `json:"probability"`
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
			return
		}

		err = json.Unmarshal(data, &person)
		if err != nil {
			log.Error("unmarshall body", zap.String("error", err.Error()))
			return
		}
		log.Info("body decoded", zap.String("name", person.Name), zap.String("surname", person.Surname))

		age, err := utils.Agify(person.Name, log)
		if err != nil {
			return
		}

		gender, err := utils.Genderize(person.Name, log)
		if err != nil {
			return
		}

		nation, err := utils.Nationalize(person.Name, log)
		if err != nil {
			return
		}

		_ = age
		_ = gender
		_ = nation

	}
}
