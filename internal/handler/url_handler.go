package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"shortener/internal/model"

	"shortener/internal/service"
)

type URLHandler struct {
	service URLService
}

type URLService interface {
	Create(
		ctx context.Context,
		originalURL string,
	) (*model.URL, error)

	FindByCode(
		ctx context.Context,
		code string,
	) (*model.URL, error)
}

func NewURLHandler(service URLService) *URLHandler {
	return &URLHandler{
		service: service,
	}
}

type createURLRequest struct {
	URL string `json:"url"`
}

type createURLResponse struct {
	Code     string `json:"code"`
	ShortURL string `json:"short_url"`
}

func (h *URLHandler) Create(w http.ResponseWriter, r *http.Request) {

	var request createURLRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"invalid_json",
			"the request body contains invalid JSON",
		)

		return
	}

	url, err := h.service.Create(
		r.Context(),
		request.URL,
	)

	if err != nil {

		if errors.Is(err, service.ErrInvalidURL) {
			writeError(
				w,
				http.StatusBadRequest,
				"invalid_url",
				"the provided URL is invalid",
			)

			return
		}

		writeError(
			w,
			http.StatusInternalServerError,
			"internal_error",
			"an internal error occurred",
		)

		return
	}

	response := createURLResponse{
		Code: url.Code,
		ShortURL: "http://localhost:8080/" +
			url.Code,
	}

	writeJSON(
		w,
		http.StatusCreated,
		response,
	)
}

func (h *URLHandler) Redirect(
	w http.ResponseWriter,
	r *http.Request,
) {

	code := r.PathValue("code")

	url, err := h.service.FindByCode(
		r.Context(),
		code,
	)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	http.Redirect(
		w,
		r,
		url.OriginalURL,
		http.StatusFound,
	)
}
