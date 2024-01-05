package get

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
)

type Request struct {
	Alias string `json:"alias"`
}

type Response struct {
	response.Response
	URL string `json:"url"`
}

type URLGetter interface {
	GetURL(alias string) (string, error)
}

func Getter(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.get"

		log = log.With(
			slog.String("op", op), slog.String("request id", middleware.GetReqID(r.Context())),
		)
		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}
		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, response.Error("Invalid request"))
			render.JSON(w, r, response.ValidationError(validateErr))

			return
		}
		alias := req.Alias
		if alias == "" {
			log.Error("alias is required", slog.String("path", r.URL.Path))
			render.JSON(w, r, response.Error("alias is required"))
			return
		}

		url, err := urlGetter.GetURL(alias)

		if err != nil {

			log.Error("alias not found", sl.Err(err))
			render.JSON(w, r, response.Error("alias not found"))
			return
		}

		render.JSON(w, r, Response{
			Response: response.OK(),
			URL:      url,
		})
	}
}