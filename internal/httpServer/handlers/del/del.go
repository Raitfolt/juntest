package del

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Raitfolt/juntest/internal/response"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type DeleteID struct {
	ID int64 `json:"id" validate:"required"`
}

type PersonDeleter interface {
	DeletePerson(id int64) error
}

func Delete(log *zap.Logger, personDeleter PersonDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("delete person by ID")

		var deleteID DeleteID

		data, err := io.ReadAll(r.Body)
		if err != nil {
			render.JSON(w, r, response.Error("read body from request"))
			log.Error("read body from request", zap.String("error", err.Error()))
			return
		}

		err = json.Unmarshal(data, &deleteID)
		if err != nil {
			render.JSON(w, r, response.Error("read body from request"))
			log.Error("unmarshall body", zap.String("error", err.Error()))
			return
		}
		log.Info("body decoded", zap.Int64("id", deleteID.ID))

		err = personDeleter.DeletePerson(deleteID.ID)
		if err != nil {
			render.JSON(w, r, response.Error("delete record error"))
			log.Error("record delete", zap.String("error", err.Error()))
			return
		}
		log.Info("record deleted", zap.Int64("id", deleteID.ID))
		render.JSON(w, r, Response{
			Response: response.OK(),
			ID:       deleteID.ID,
		})
	}
}

type Response struct {
	response.Response
	ID int64 `json:"id,omitempty"`
}
