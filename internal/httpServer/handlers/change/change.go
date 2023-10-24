package change

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Raitfolt/juntest/internal/response"
	"github.com/Raitfolt/juntest/internal/utils"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type Person struct {
	Name       string `json:"name" validate:"required"`
	Surname    string `json:"surname" validate:"required"`
	Patronymic string `json:"patronymic,omitempty"`
}

type PersonSaver interface {
	NewPerson(name, surname, patronymic string, age int, gender, nationality string) (int64, error)
}

func Change(log *zap.Logger, personSaver PersonSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("change person information")

		var person Person

		data, err := io.ReadAll(r.Body)
		if err != nil {
			//TODO: do a function with 2 errors: to json and to logs
			render.JSON(w, r, response.Error("read body from request"))
			log.Error("read body from request", zap.String("error", err.Error()))
			return
		}

		err = json.Unmarshal(data, &person)
		if err != nil {
			render.JSON(w, r, response.Error("read body from request"))
			log.Error("unmarshall body", zap.String("error", err.Error()))
			return
		}
		log.Info("body decoded", zap.String("name", person.Name), zap.String("surname", person.Surname))

		age, err := utils.Agify(person.Name, log)
		if err != nil {
			log.Error("get age", zap.String("name", person.Name), zap.String("error", err.Error()))
			return
		}

		gender, err := utils.Genderize(person.Name, log)
		if err != nil {
			log.Error("get gender", zap.String("name", person.Name), zap.String("error", err.Error()))
			return
		}

		nation, err := utils.Nationalize(person.Name, log)
		if err != nil {
			log.Error("get nation", zap.String("name", person.Name), zap.String("error", err.Error()))
			return
		}

		pos, err := personSaver.NewPerson(
			person.Name,
			person.Surname,
			person.Patronymic,
			age,
			gender,
			nation)
		if err != nil {
			render.JSON(w, r, response.Error("record insert error"))
			log.Error("record insert", zap.String("error", err.Error()))
			return
		}
		log.Info("record added", zap.Int64("id", pos))
		render.JSON(w, r, Response{
			Response: response.OK(),
			ID:       pos,
		})
	}
}

type Response struct {
	response.Response
	ID int64 `json:"id,omitempty"`
}
