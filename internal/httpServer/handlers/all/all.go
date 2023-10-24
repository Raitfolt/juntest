package all

import (
	"net/http"

	"github.com/Raitfolt/juntest/internal/response"
	"github.com/Raitfolt/juntest/internal/storage/psql"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type PersonGetter interface {
	ListPersons() ([]psql.Person, error)
}

func List(log *zap.Logger, personGetter PersonGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("get all persons")

		persons, err := personGetter.ListPersons()

		if err != nil {
			render.JSON(w, r, response.Error("get list of persons error"))
			log.Error("get list of persons", zap.String("error", err.Error()))
			return
		}

		log.Info("list of persons loaded")

		//jsonPersons, err := json.Marshal(persons)
		if err != nil {
			render.JSON(w, r, response.Error("persons to json marshal error"))
			log.Error("persons to json marshal", zap.String("error", err.Error()))
			return
		}
		render.JSON(w, r, persons)
	}
}
