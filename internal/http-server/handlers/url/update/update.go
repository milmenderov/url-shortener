package update

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
	URL   string `json:"url"`
	Alias string `json:"alias"`
}

type Response struct {
	response.Response
}

type URLUpdater interface {
	UpdateURL(newUrl string, alias string) error
}

func Updater(log *slog.Logger, urlUpdater URLUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.updater"

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
		url := req.URL
		if url == "" {
			log.Error("url is required", slog.String("path", r.URL.Path))
			render.JSON(w, r, response.Error("url is required"))
			return
		}
		err = urlUpdater.UpdateURL(url, alias)
		if err != nil {
			log.Error("failed to update url", sl.Err(err))
			render.JSON(w, r, response.Error("failed to update url"))
			return
		}
		render.JSON(w, r, Response{
			Response: response.OK(),
		})
	}
}
