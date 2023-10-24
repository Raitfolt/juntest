package change

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Raitfolt/juntest/internal/response"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type Person struct {
	ID          int    `json:"id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Surname     string `json:"surname" validate:"required"`
	Patronymic  string `json:"patronymic,omitempty"`
	Age         int    `json:"age"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
}

type PersonChanger interface {
	ChangePerson(id int, name, surname, patronymic string, age int, gender, nationality string) (int64, error)
}

func Change(log *zap.Logger, personChanger PersonChanger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("change person information")

		var person Person

		data, err := io.ReadAll(r.Body)
		if err != nil {
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

		pos, err := personChanger.ChangePerson(
			person.ID,
			person.Name,
			person.Surname,
			person.Patronymic,
			person.Age,
			person.Gender,
			person.Nationality)
		if err != nil {
			render.JSON(w, r, response.Error("record change error"))
			log.Error("record change", zap.String("error", err.Error()))
			return
		}
		log.Info("record changed", zap.Int64("id", pos))
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
