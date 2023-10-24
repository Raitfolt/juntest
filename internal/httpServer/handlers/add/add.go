package add

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-resty/resty/v2"
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

		var age int
		resty.New().R().SetResult(&age).Get("https://api.agify.io/?name=" + person.Name)
		log.Info("get age", zap.String("name", person.Name), zap.Int("age", person.Age))

		var gender string
		resty.New().R().SetResult(&gender).Get("https://api.genderize.io/?name=" + person.Name)
		log.Info("get gender", zap.String("name", person.Name), zap.String("gender", person.Gender))

		var nationality string
		resty.New().R().SetResult(&nationality).Get("https://api.nationalize.io/?name=" + person.Name)
		log.Info("get nationality", zap.String("name", person.Name), zap.String("age", person.Nationality))
	}
}
